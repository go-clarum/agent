package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/go-clarum/agent/config"
	"github.com/go-clarum/agent/control"
	"github.com/go-clarum/agent/logging"
	"github.com/go-clarum/agent/services/http/internal/constants"
	"github.com/go-clarum/agent/services/http/internal/validators"
	clarumstrings "github.com/go-clarum/agent/validators/strings"
	"io"
	"net"
	"net/http"
	"time"
)

const contextNameKey = "endpointContext"

type endpoint struct {
	name                     string
	port                     uint
	contentType              string
	server                   *http.Server
	serverTimeout            time.Duration
	context                  context.Context
	cancelContext            context.CancelFunc
	requestValidationChannel chan *http.Request
	sendChannel              chan *sendPair
	logger                   *logging.Logger
}

type endpointContext struct {
	endpointName             string
	requestValidationChannel chan *http.Request
	sendChannel              chan *sendPair
	logger                   *logging.Logger
}

type sendPair struct {
	response *SendAction
	error    error
}

func newEndpoint(is *InitRequest) *endpoint {
	ctx, cancelCtx := context.WithCancel(context.Background())
	sendChannel := make(chan *sendPair)
	requestChannel := make(chan *http.Request)

	se := &endpoint{
		name:                     is.Name,
		port:                     is.Port,
		contentType:              is.ContentType,
		serverTimeout:            is.TimeoutSeconds,
		context:                  ctx,
		cancelContext:            cancelCtx,
		sendChannel:              sendChannel,
		requestValidationChannel: requestChannel,
		logger:                   logging.NewLogger(loggerName(is.Name)),
	}

	return se
}

// this Method is blocking, until a request is received
func (endpoint *endpoint) receive(action *ReceiveAction) (*http.Request, error) {
	endpoint.logger.Debugf("action to receive %s", action.ToString())
	endpoint.enrichReceiveAction(action)

	select {
	case receivedRequest := <-endpoint.requestValidationChannel:
		endpoint.logger.Debugf("validation action %s", action.ToString())

		return receivedRequest, errors.Join(
			validators.ValidatePath(action.Path, receivedRequest.URL, endpoint.logger),
			validators.ValidateHttpMethod(action.Method, receivedRequest.Method, endpoint.logger),
			validators.ValidateHttpHeaders(action.Headers, receivedRequest.Header, endpoint.logger),
			validators.ValidateHttpQueryParams(action.QueryParams, receivedRequest.URL, endpoint.logger),
			validators.ValidateHttpPayload(&action.Payload, receivedRequest.Body,
				action.PayloadType, endpoint.logger))
	case <-time.After(config.ActionTimeout()):
		return nil, endpoint.handleError("receive action timed out - no request received for validation", nil)
	}
}

func (endpoint *endpoint) send(action *SendAction) error {
	endpoint.enrichSendAction(action)
	err := endpoint.validateMessageToSend(action)

	// we must always send a signal downstream so that the handler is not blocked
	toSend := &sendPair{
		response: action,
		error:    err,
	}

	select {
	case endpoint.sendChannel <- toSend:
		return err
	case <-time.After(config.ActionTimeout()):
		return endpoint.handleError("send action timed out - no request received for validation", nil)
	}
}

func (endpoint *endpoint) enrichReceiveAction(action *ReceiveAction) {
	if clarumstrings.IsNotBlank(endpoint.contentType) {
		if _, exists := action.Headers[constants.ContentTypeHeaderName]; !exists {
			action.Headers[constants.ContentTypeHeaderName] = endpoint.contentType
		}
	}
}

func (endpoint *endpoint) enrichSendAction(action *SendAction) {
	// if no Headers have been sent by the bindings, this will be nil
	if action.Headers == nil {
		action.Headers = make(map[string]string)
	}

	if clarumstrings.IsBlank(action.Headers[constants.ContentTypeHeaderName]) {
		action.Headers[constants.ContentTypeHeaderName] = endpoint.contentType
	}
}

func (endpoint *endpoint) start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", requestHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", endpoint.port),
		Handler:      mux,
		WriteTimeout: endpoint.serverTimeout,
		BaseContext: func(l net.Listener) context.Context {
			endpointContext := &endpointContext{
				endpointName:             endpoint.name,
				requestValidationChannel: endpoint.requestValidationChannel,
				sendChannel:              endpoint.sendChannel,
				logger:                   endpoint.logger,
			}

			return context.WithValue(endpoint.context, contextNameKey, endpointContext)
		},
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			endpoint.logger.Errorf("error - %s", err)
		} else {
			endpoint.logger.Info("closed server")
		}

		endpoint.cancelContext()
	}()

	endpoint.server = server
}

// The requestHandler is started when the server receives a request.
// The request is sent to the requestValidationChannel to be picked up by a test action (validation).
// After sending the request to the channel, the handler is blocked until the send() test action
// provides a response message. This way we can control, inside the test, when a response will be sent.
// The handler blocks until a timeout is triggered
func requestHandler(resWriter http.ResponseWriter, request *http.Request) {
	control.RunningActions.Add(1)
	ctx := request.Context().Value(contextNameKey).(*endpointContext)
	defer finishOrRecover(ctx.logger)

	logIncomingRequest(ctx.logger, request)

	select {
	case ctx.requestValidationChannel <- request:
		ctx.logger.Debug("received request was sent to validation channel")
	case <-time.After(config.ActionTimeout()):
		ctx.logger.Warn("request handling timed out - no server receive action called in test")
	}

	select {
	case sendPair := <-ctx.sendChannel:
		// error from upstream - we send a response to close the HTTP cycle
		if sendPair.error != nil {
			sendDefaultErrorResponse(ctx.logger, "request handler received error from upstream", resWriter)
			return
		}

		// check if response is empty - we send a response to close the HTTP cycle
		if sendPair.response == nil {
			sendDefaultErrorResponse(ctx.logger, "request handler received empty ResponseMesage", resWriter)
			return
		}

		sendResponse(ctx.logger, sendPair, resWriter)
	case <-time.After(config.ActionTimeout()):
		ctx.logger.Warn("response handling timed out - no server send action called in test")
	}
}

func (endpoint *endpoint) shutdown() {
	errorCh := make(chan error, 1)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	go func() {
		if err := endpoint.server.Shutdown(ctxTimeout); err != nil {
			errorCh <- err
		}

		errorCh <- nil
	}()

	select {
	case <-ctxTimeout.Done():
		endpoint.logger.Errorf("shutdown timeout reached for endpoint [%s]", endpoint.name)
	case err := <-errorCh:
		if err != nil {
			endpoint.logger.Errorf("shutdown failed for endpoint [%s] - %s", endpoint.name, err)
		}
		endpoint.logger.Infof("successfully shutdown endpoint [%s]", endpoint.name)
	}
}

func sendResponse(logger *logging.Logger, sendPair *sendPair, resWriter http.ResponseWriter) {
	for header, value := range sendPair.response.Headers {
		resWriter.Header().Set(header, value)
	}

	resWriter.WriteHeader(sendPair.response.StatusCode)

	_, err := io.WriteString(resWriter, sendPair.response.Payload)
	if err != nil {
		logger.Errorf("could not write response body - %s", err)
	}
	logOutgoingResponse(logger, sendPair.response.StatusCode, sendPair.response.Payload, resWriter)
}

func sendDefaultErrorResponse(logger *logging.Logger, errorMessage string, resWriter http.ResponseWriter) {
	logger.Error(errorMessage)
	resWriter.WriteHeader(http.StatusInternalServerError)
	logOutgoingResponse(logger, http.StatusInternalServerError, "", resWriter)
}

func (endpoint *endpoint) validateMessageToSend(action *SendAction) error {
	if action.StatusCode < 100 || action.StatusCode > 999 {
		return endpoint.handleError(fmt.Sprintf("action to send is invalid - unsupported status code [%d]",
			action.StatusCode), nil)
	}

	return nil
}

func (endpoint *endpoint) handleError(message string, err error) error {
	var errorMessage string
	if err != nil {
		errorMessage = message + " - " + err.Error()
	} else {
		errorMessage = message
	}
	endpoint.logger.Errorf(errorMessage)
	return errors.New(endpoint.logger.Name() + " " + errorMessage)
}

func finishOrRecover(logger *logging.Logger) {
	control.RunningActions.Done()

	if r := recover(); r != nil {
		logger.Errorf("endpoint panicked: error - %s", r)
	}
}

// we read the body 'as is' for logging, after which we put it back into the request
// with an open reader so that it can be read downstream again
func logIncomingRequest(logger *logging.Logger, request *http.Request) {
	bodyBytes, _ := io.ReadAll(request.Body)
	bodyString := ""

	err := request.Body.Close()
	if err != nil {
		logger.Errorf("could not read request body - %s", err)
	} else {
		request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		bodyString = string(bodyBytes)
	}

	logger.Infof("received HTTP request ["+
		"method: %s, "+
		"url: %s, "+
		"headers: %s, "+
		"payload: %s"+
		"]",
		request.Method, request.URL.String(), request.Header, bodyString)
}

func logOutgoingResponse(logger *logging.Logger, statusCode int, payload string, res http.ResponseWriter) {
	logger.Infof("sending response ["+
		"status: %d, "+
		"headers: %s, "+
		"payload: %s"+
		"]",
		statusCode, res.Header(), payload)
}

func loggerName(endpointName string) string {
	return fmt.Sprintf("%s:", endpointName)
}
