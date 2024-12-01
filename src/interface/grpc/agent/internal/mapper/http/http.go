package http

import (
	clientCommands "github.com/go-clarum/agent/application/command/http/client/commands"
	"github.com/go-clarum/agent/application/command/http/common/model"
	serverCommands "github.com/go-clarum/agent/application/command/http/server/commands"
	api "github.com/go-clarum/agent/interface/grpc/agent/internal/api/commands/http"
	"time"
)

func NewClientInitCommandFrom(is *api.InitClientCommand) *clientCommands.InitEndpointCommand {
	return &clientCommands.InitEndpointCommand{
		Name:           is.Name,
		BaseUrl:        is.BaseUrl,
		ContentType:    is.ContentType,
		TimeoutSeconds: time.Duration(is.TimeoutSeconds) * time.Second,
	}
}

func NewClientSendActionFrom(sa *api.ClientSendActionCommand) *clientCommands.SendCommand {
	return &clientCommands.SendCommand{
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

func NewClientReceiveActionFrom(sa *api.ClientReceiveActionCommand) *clientCommands.ReceiveCommand {
	return &clientCommands.ReceiveCommand{
		Name:         sa.Name,
		PayloadType:  model.PayloadType(sa.PayloadType),
		StatusCode:   int(sa.StatusCode),
		Headers:      sa.Headers,
		Payload:      sa.Payload,
		EndpointName: sa.EndpointName,
	}
}

func NewServerInitRequestFrom(is *api.InitServerCommand) *serverCommands.InitEndpointCommand {
	return &serverCommands.InitEndpointCommand{
		Name:           is.Name,
		Port:           uint(is.Port),
		TimeoutSeconds: time.Duration(is.TimeoutSeconds) * time.Second,
	}
}

func NewServerSendActionFrom(sa *api.ServerSendActionCommand) *serverCommands.SendCommand {
	return &serverCommands.SendCommand{
		Name:         sa.Name,
		StatusCode:   int(sa.StatusCode),
		Headers:      sa.Headers,
		Payload:      sa.Payload,
		EndpointName: sa.EndpointName,
	}
}

func NewServerReceiveActionFrom(ra *api.ServerReceiveActionCommand) *serverCommands.ReceiveCommand {
	return &serverCommands.ReceiveCommand{
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
