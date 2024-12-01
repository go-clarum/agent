package server

import (
	"github.com/go-clarum/agent/application/command/common"
	"github.com/go-clarum/agent/application/command/http/server/commands"
	"github.com/go-clarum/agent/application/command/http/server/internal"
	"github.com/go-clarum/agent/infrastructure/logging"
	"net/http"
)

type handler struct {
	endpoints map[string]*internal.Endpoint
	logger    *logging.Logger
}

func NewHttpServerHandler() common.CommandHandler {
	return &handler{
		endpoints: make(map[string]*internal.Endpoint),
		logger:    logging.NewLogger("HttpServerService"),
	}
}

func (h *handler) CanHandle(command any) bool {
	return true
}

func (h *handler) Handle(command any) any {
	return nil
}

func (h *handler) InitializeEndpoint(is *commands.InitEndpointCommand) error {
	newEndpoint := internal.NewEndpoint(is)

	if oldEndpoint, exists := h.endpoints[newEndpoint.Name]; exists {
		h.logger.Infof("endpoint [%s] already exists - replacing", oldEndpoint.Name)
		oldEndpoint.Shutdown()
	}

	newEndpoint.Start()

	h.endpoints[newEndpoint.Name] = newEndpoint
	logging.Infof("registered HTTP server endpoint [%s]", newEndpoint.Name)

	return nil
}

func (h *handler) SendAction(sendAction *commands.SendCommand) error {
	endpoint, exists := h.endpoints[sendAction.EndpointName]
	if !exists {
		h.logger.Errorf("HTTP server endpoint [%s] not found - action [%s] will not be executed",
			sendAction.EndpointName, sendAction.Name)
	}

	return endpoint.Send(sendAction)
}

func (h *handler) ReceiveAction(receiveAction *commands.ReceiveCommand) (*http.Request, error) {
	endpoint, exists := h.endpoints[receiveAction.EndpointName]
	if !exists {
		h.logger.Errorf("HTTP server endpoint [%s] not found - action [%s] will not be executed",
			receiveAction.EndpointName, receiveAction.Name)
	}

	return endpoint.Receive(receiveAction)
}
