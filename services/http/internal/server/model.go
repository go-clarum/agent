package server

import (
	"fmt"
	api "github.com/go-clarum/agent/api/http"
	"github.com/go-clarum/agent/services/http/internal"
	"time"
)

// the purpose of this layer is to separate the internal model from the grpc one
// only data type mapping should happen here, no business logic (like setting defaults)

type initializeRequest struct {
	name           string
	port           uint
	contentType    string
	timeoutSeconds time.Duration
}

type sendAction struct {
	name         string
	payloadType  internal.PayloadType
	statusCode   int
	headers      map[string]string
	payload      string
	endpointName string
}

type receiveAction struct {
	name         string
	url          string
	path         string
	method       string
	queryParams  map[string][]string
	headers      map[string]string
	payload      string
	payloadType  internal.PayloadType
	endpointName string
}

func NewInitializeRequestFrom(is *api.InitializeServerRequest) *initializeRequest {
	return &initializeRequest{
		name:           is.Name,
		port:           uint(is.Port),
		timeoutSeconds: time.Duration(is.TimeoutSeconds) * time.Second,
	}
}

func NewSendActionFrom(sa *api.ServerSendActionRequest) *sendAction {
	return &sendAction{
		name:         sa.Name,
		statusCode:   int(sa.StatusCode),
		headers:      sa.Headers,
		payload:      sa.Payload,
		endpointName: sa.EndpointName,
	}
}

func NewReceiveActionFrom(ra *api.ServerReceiveActionRequest) *receiveAction {
	return &receiveAction{
		name:         ra.Name,
		url:          ra.Url,
		path:         ra.Path,
		method:       ra.Method,
		queryParams:  parseQueryParams(ra.QueryParams),
		headers:      ra.Headers,
		payload:      ra.Payload,
		payloadType:  internal.PayloadType(ra.PayloadType),
		endpointName: ra.EndpointName,
	}
}

func (action *receiveAction) ToString() string {
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

func parseQueryParams(apiQueryParams map[string]*api.StringsList) map[string][]string {
	result := make(map[string][]string)

	if apiQueryParams != nil {
		for key, value := range apiQueryParams {
			result[key] = value.Values
		}
	}

	return result
}
