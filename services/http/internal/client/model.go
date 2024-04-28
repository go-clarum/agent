package client

import (
	"fmt"
	api "github.com/go-clarum/agent/api/http"
	"github.com/go-clarum/agent/services/http/internal"
	"strconv"
	"time"
)

// the purpose of this layer is to separate the internal model from the grpc one
// only data type mapping should happen here, no business logic (like setting defaults)

type initializeRequest struct {
	name           string
	baseUrl        string
	contentType    string
	timeoutSeconds time.Duration
}

type sendAction struct {
	name         string
	url          string
	path         string
	method       string
	queryParams  map[string][]string
	headers      map[string]string
	payload      string
	endpointName string
}

type receiveAction struct {
	name         string
	payloadType  internal.PayloadType
	statusCode   int
	headers      map[string]string
	payload      string
	endpointName string
}

func NewInitializeRequestFrom(is *api.InitializeClientRequest) *initializeRequest {
	return &initializeRequest{
		name:           is.Name,
		baseUrl:        is.BaseUrl,
		contentType:    is.ContentType,
		timeoutSeconds: time.Duration(is.TimeoutSeconds) * time.Second,
	}
}

func NewSendActionFrom(sa *api.ClientSendActionRequest) *sendAction {
	return &sendAction{
		name:         sa.Name,
		url:          sa.Url,
		path:         sa.Path,
		method:       sa.Method,
		queryParams:  parseQueryParams(sa.QueryParams),
		headers:      sa.Headers,
		payload:      sa.Payload,
		endpointName: sa.EndpointName,
	}
}

func NewReceiveActionFrom(sa *api.ClientReceiveActionRequest) *receiveAction {
	return &receiveAction{
		name:         sa.Name,
		payloadType:  internal.PayloadType(sa.PayloadType),
		statusCode:   int(sa.StatusCode),
		headers:      sa.Headers,
		payload:      sa.Payload,
		endpointName: sa.EndpointName,
	}
}

func (action *sendAction) ToString() string {
	return fmt.Sprintf(
		"["+
			"Method: %s, "+
			"BaseUrl: %s, "+
			"Path: '%s', "+
			"Headers: %s, "+
			"QueryParams: %s, "+
			"Payload: %s"+
			"]",
		action.method, action.url, action.path,
		action.headers, action.queryParams, action.payload)
}

func (action *receiveAction) ToString() string {
	statusCodeText := "none"
	if action.statusCode > 0 {
		statusCodeText = strconv.Itoa(action.statusCode)
	}

	return fmt.Sprintf(
		"["+
			"StatusCode: %s, "+
			"Headers: %s, "+
			"MessagePayload: %s"+
			"]",
		statusCodeText, action.headers, action.payload)
}

func parseQueryParams(apiQueryParams map[string]*api.StringsList) map[string][]string {
	result := make(map[string][]string)

	if apiQueryParams != nil {
		for key, value := range apiQueryParams {
			result[key] = value.Values
		}
	}

	return result
}
