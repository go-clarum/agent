package mappers

import (
	clientActions "github.com/go-clarum/agent/application/services/http/client/actions"
	"github.com/go-clarum/agent/application/services/http/common/model"
	serverActions "github.com/go-clarum/agent/application/services/http/server/actions"
	"github.com/go-clarum/agent/interface/grpc/http/internal/api"
	"time"
)

// the purpose of this layer is to separate the internal model from the grpc one
// only data type mapping should happen here, no business logic (like setting defaults)

func NewClientInitRequestFrom(is *api.InitClientRequest) *clientActions.InitEndpointAction {
	return &clientActions.InitEndpointAction{
		Name:           is.Name,
		BaseUrl:        is.BaseUrl,
		ContentType:    is.ContentType,
		TimeoutSeconds: time.Duration(is.TimeoutSeconds) * time.Second,
	}
}

func NewClientSendActionFrom(sa *api.ClientSendActionRequest) *clientActions.SendAction {
	return &clientActions.SendAction{
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

func NewClientReceiveActionFrom(sa *api.ClientReceiveActionRequest) *clientActions.ReceiveAction {
	return &clientActions.ReceiveAction{
		Name:         sa.Name,
		PayloadType:  model.PayloadType(sa.PayloadType),
		StatusCode:   int(sa.StatusCode),
		Headers:      sa.Headers,
		Payload:      sa.Payload,
		EndpointName: sa.EndpointName,
	}
}

func NewServerInitRequestFrom(is *api.InitServerRequest) *serverActions.InitEndpointAction {
	return &serverActions.InitEndpointAction{
		Name:           is.Name,
		Port:           uint(is.Port),
		TimeoutSeconds: time.Duration(is.TimeoutSeconds) * time.Second,
	}
}

func NewServerSendActionFrom(sa *api.ServerSendActionRequest) *serverActions.SendAction {
	return &serverActions.SendAction{
		Name:         sa.Name,
		StatusCode:   int(sa.StatusCode),
		Headers:      sa.Headers,
		Payload:      sa.Payload,
		EndpointName: sa.EndpointName,
	}
}

func NewServerReceiveActionFrom(ra *api.ServerReceiveActionRequest) *serverActions.ReceiveAction {
	return &serverActions.ReceiveAction{
		Name:         ra.Name,
		Url:          ra.Url,
		Path:         ra.Path,
		Method:       ra.Method,
		QueryParams:  parseQueryParams(ra.QueryParams),
		Headers:      ra.Headers,
		Payload:      ra.Payload,
		PayloadType:  model.PayloadType(ra.PayloadType),
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
