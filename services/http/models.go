package http

import (
	api "github.com/go-clarum/agent/api/http"
	"github.com/go-clarum/agent/services/http/internal"
	"github.com/go-clarum/agent/services/http/internal/client"
	"github.com/go-clarum/agent/services/http/internal/server"
	"time"
)

// the purpose of this layer is to separate the internal model from the grpc one
// only data type mapping should happen here, no business logic (like setting defaults)

func newClientInitRequestFrom(is *api.InitClientRequest) *client.InitRequest {
	return &client.InitRequest{
		Name:           is.Name,
		BaseUrl:        is.BaseUrl,
		ContentType:    is.ContentType,
		TimeoutSeconds: time.Duration(is.TimeoutSeconds) * time.Second,
	}
}

func newClientSendActionFrom(sa *api.ClientSendActionRequest) *client.SendAction {
	return &client.SendAction{
		Name:         sa.Name,
		Url:          sa.Url,
		Path:         sa.Path,
		Method:       sa.Method,
		QueryParams:  parseQueryParams(sa.QueryParams),
		Headers:      sa.Headers,
		Payload:      sa.Payload,
		EndpointName: sa.EndpointName,
	}
}

func newClientReceiveActionFrom(sa *api.ClientReceiveActionRequest) *client.ReceiveAction {
	return &client.ReceiveAction{
		Name:         sa.Name,
		PayloadType:  internal.PayloadType(sa.PayloadType),
		StatusCode:   int(sa.StatusCode),
		Headers:      sa.Headers,
		Payload:      sa.Payload,
		EndpointName: sa.EndpointName,
	}
}

func newServerInitRequestFrom(is *api.InitServerRequest) *server.InitRequest {
	return &server.InitRequest{
		Name:           is.Name,
		Port:           uint(is.Port),
		TimeoutSeconds: time.Duration(is.TimeoutSeconds) * time.Second,
	}
}

func newServerSendActionFrom(sa *api.ServerSendActionRequest) *server.SendAction {
	return &server.SendAction{
		Name:         sa.Name,
		StatusCode:   int(sa.StatusCode),
		Headers:      sa.Headers,
		Payload:      sa.Payload,
		EndpointName: sa.EndpointName,
	}
}

func newServerReceiveActionFrom(ra *api.ServerReceiveActionRequest) *server.ReceiveAction {
	return &server.ReceiveAction{
		Name:         ra.Name,
		Url:          ra.Url,
		Path:         ra.Path,
		Method:       ra.Method,
		QueryParams:  parseQueryParams(ra.QueryParams),
		Headers:      ra.Headers,
		Payload:      ra.Payload,
		PayloadType:  internal.PayloadType(ra.PayloadType),
		EndpointName: ra.EndpointName,
	}
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
