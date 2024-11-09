package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-clarum/agent/config"
	"github.com/go-clarum/agent/control"
	"github.com/go-clarum/agent/durations"
	"github.com/go-clarum/agent/logging"
	"github.com/go-clarum/agent/services/http/internal/constants"
	"github.com/go-clarum/agent/services/http/internal/utils"
	"github.com/go-clarum/agent/services/http/internal/validators"
	clarumstrings "github.com/go-clarum/agent/validators/strings"
	"io"
	"net/http"
	"time"
)

type endpoint struct {
	name            string
	baseUrl         string
	contentType     string
	client          *http.Client
	responseChannel chan *responsePair
	logger          *logging.Logger
}

type responsePair struct {
	response *http.Response
	error    error
}

func newEndpoint(ic *InitRequest) (*endpoint, error) {
	if clarumstrings.IsBlank(ic.Name) {
		return nil, errors.New("cannot create HTTP client endpoint - name is empty")
	}

	client := http.Client{
		Timeout: durations.GetDurationWithDefault(ic.TimeoutSeconds, 10*time.Second),
	}

	return &endpoint{
		name:            ic.Name,
		baseUrl:         ic.BaseUrl,
		contentType:     ic.ContentType,
		client:          &client,
		responseChannel: make(chan *responsePair),
		logger:          logging.NewLogger(loggerName(ic.Name)),
	}, nil
}

func (endpoint *endpoint) send(action *SendAction) error {
	if action == nil {
		return endpoint.handleError("send action is nil", nil)
	}

	endpoint.logger.Debugf("action to send [%s]", action.ToString())
	endpoint.enrichSendAction(action)
	endpoint.logger.Debugf("will send action [%s]", action.ToString())

	if err := endpoint.validateMessageToSend(action); err != nil {
		return err
	}

	req, err := endpoint.buildRequest(action)
	// we return error here directly and not in the goroutine below
	// this way we can signal to the test synchronously that there was an error
	if err != nil {
		return endpoint.handleError("canceled message", err)
	}

	go func() {
		control.RunningActions.Add(1)
		defer control.RunningActions.Done()

		endpoint.logOutgoingRequest(action.Payload, req)
		res, err := endpoint.client.Do(req)

		// we log the error here directly, but will do error handling downstream
		if err != nil {
			endpoint.logger.Errorf("error on response - %s", err)
			defer closeBody(res)
		} else {
			endpoint.logIncomingResponse(res)
		}

		responsePair := &responsePair{
			response: res,
			error:    err,
		}

		select {
		// we send the error downstream for it to be returned when an action is called
		case endpoint.responseChannel <- responsePair:
		case <-time.After(config.ActionTimeout()):
			endpoint.handleError("action timed out - no client receive action called in test", nil)
		}
	}()

	return nil
}

func closeBody(res *http.Response) {
	if res.Body != nil {
		res.Body.Close()
	}
}

// validationOptions pass by value is intentional
func (endpoint *endpoint) receive(action *ReceiveAction) (*http.Response, error) {
	if action == nil {
		return nil, endpoint.handleError("receive action is nil", nil)
	}
	endpoint.logger.Debugf("action to receive [%s]", action.ToString())

	select {
	case responsePair := <-endpoint.responseChannel:
		if responsePair.error != nil {
			return responsePair.response, endpoint.handleError("error while receiving response", responsePair.error)
		}

		endpoint.enrichReceiveAction(action)
		endpoint.logger.Debugf("validating receive action [%s]", action.ToString())

		return responsePair.response, errors.Join(
			validators.ValidateHttpStatusCode(action.StatusCode, responsePair.response.StatusCode, endpoint.logger),
			validators.ValidateHttpHeaders(action.Headers, responsePair.response.Header, endpoint.logger),
			validators.ValidateHttpPayload(&action.Payload, responsePair.response.Body,
				action.PayloadType, endpoint.logger))
	case <-time.After(config.ActionTimeout()):
		return nil, endpoint.handleError("receive action timed out - no response received for validation", nil)
	}
}

// Put missing data into a message to send: baseUrl & ContentType Header
func (endpoint *endpoint) enrichSendAction(action *SendAction) {
	if clarumstrings.IsBlank(action.Url) {
		action.Url = endpoint.baseUrl
	}

	// if no Headers have been sent by the bindings, this will be nil
	if action.Headers == nil {
		action.Headers = make(map[string]string)
	}

	if clarumstrings.IsBlank(action.Headers[constants.ContentTypeHeaderName]) {
		action.Headers[constants.ContentTypeHeaderName] = endpoint.contentType
	}
}

// Put missing data into message to receive: ContentType Header
func (endpoint *endpoint) enrichReceiveAction(action *ReceiveAction) {
	if clarumstrings.IsNotBlank(endpoint.contentType) {
		if _, exists := action.Headers[constants.ContentTypeHeaderName]; !exists {
			action.Headers[constants.ContentTypeHeaderName] = endpoint.contentType
		}
	}
}

func (endpoint *endpoint) validateMessageToSend(action *SendAction) error {
	if clarumstrings.IsBlank(action.Method) {
		return endpoint.handleError("send action is invalid - missing HTTP method", nil)
	}
	if clarumstrings.IsBlank(action.Url) {
		return endpoint.handleError("send action is invalid - missing url", nil)
	}
	if !utils.IsValidUrl(action.Url) {
		return endpoint.handleError("send action is invalid - invalid url", nil)
	}

	return nil
}

func (endpoint *endpoint) buildRequest(action *SendAction) (*http.Request, error) {
	url := utils.BuildPath(action.Url, action.Path...)

	req, err := http.NewRequest(action.Method, url, bytes.NewBufferString(action.Payload))
	if err != nil {
		endpoint.logger.Errorf("error - %s", err)
		return nil, err
	}

	for header, value := range action.Headers {
		req.Header.Set(header, value)
	}

	qParams := req.URL.Query()
	for key, values := range action.QueryParams {
		for _, value := range values {
			qParams.Add(key, value)
		}
	}
	req.URL.RawQuery = qParams.Encode()

	return req, nil
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

func (endpoint *endpoint) logOutgoingRequest(payload string, req *http.Request) {
	endpoint.logger.Infof("sending HTTP request ["+
		"method: %s, "+
		"url: %s, "+
		"headers: %s, "+
		"payload: %s"+
		"]",
		req.Method, req.URL, req.Header, payload)
}

// we read the body 'as is' for logging, after which we put it back into the response
// with an open reader so that it can be read downstream again
func (endpoint *endpoint) logIncomingResponse(res *http.Response) {
	bodyBytes, _ := io.ReadAll(res.Body)
	bodyString := ""

	err := res.Body.Close()
	if err != nil {
		endpoint.logger.Errorf("could not read response body - %s", err)
	} else {
		res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		bodyString = string(bodyBytes)
	}

	endpoint.logger.Infof("received HTTP response ["+
		"status: %s, "+
		"headers: %s, "+
		"payload: %s"+
		"]",
		res.Status, res.Header, bodyString)
}

func loggerName(endpointName string) string {
	return fmt.Sprintf("%s:", endpointName)
}
