package client

import (
	"github.com/go-clarum/agent/application/command/common"
	"github.com/go-clarum/agent/application/command/http/client/commands"
	"github.com/go-clarum/agent/application/command/http/client/internal"
	"github.com/go-clarum/agent/infrastructure/logging"
	"net/http"
)

type handler struct {
	endpoints map[string]*internal.Endpoint
	logger    *logging.Logger
}

func (h *handler) CanHandle(command any) bool {
	return true
}

func (h *handler) Handle(command any) any {
	return nil
}

func NewHttpClientHandler() common.CommandHandler {
	return &handler{
		endpoints: make(map[string]*internal.Endpoint),
		logger:    logging.NewLogger("HttpClientHandler"),
	}
}

func (h *handler) InitializeEndpoint(req *commands.InitEndpointCommand) error {
	newEndpoint, err := internal.NewEndpoint(req)

	if err != nil {
		h.logger.Errorf("failed to initialize HTTP client endpoint - %s", err)
		return err
	}

	if oldEndpoint, exists := h.endpoints[newEndpoint.Name]; exists {
		h.logger.Infof("HTTP client endpoint [%s] already exists - replacing", oldEndpoint.Name)
	}

	h.endpoints[newEndpoint.Name] = newEndpoint
	logging.Infof("registered HTTP client endpoint [%s]", newEndpoint.Name)

	return nil
}

func (h *handler) SendAction(sendAction *commands.SendCommand) error {
	endpoint, exists := h.endpoints[sendAction.EndpointName]
	if !exists {
		h.logger.Errorf("HTTP client endpoint [%s] not found - action [%s] will not be executed",
			sendAction.EndpointName, sendAction.Name)
	}

	return endpoint.Send(sendAction)
}

func (h *handler) ReceiveAction(receiveAction *commands.ReceiveCommand) (*http.Response, error) {
	endpoint, exists := h.endpoints[receiveAction.EndpointName]
	if !exists {
		h.logger.Errorf("HTTP client endpoint [%s] not found - action [%s] will not be executed",
			receiveAction.EndpointName, receiveAction.Name)
	}

	return endpoint.Receive(receiveAction)
}
