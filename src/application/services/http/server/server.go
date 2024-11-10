package server

import (
	"github.com/go-clarum/agent/application/services/http/server/actions"
	"github.com/go-clarum/agent/application/services/http/server/interfaces"
	"github.com/go-clarum/agent/application/services/http/server/internal"
	"github.com/go-clarum/agent/infrastructure/logging"
	"net/http"
)

type service struct {
	endpoints map[string]*internal.Endpoint
	logger    *logging.Logger
}

func NewHttpServerService() interfaces.ServerService {
	return &service{
		endpoints: make(map[string]*internal.Endpoint),
		logger:    logging.NewLogger("HttpServerService"),
	}
}

func (s *service) InitializeEndpoint(is *actions.InitEndpointAction) error {
	newEndpoint := internal.NewEndpoint(is)

	if oldEndpoint, exists := s.endpoints[newEndpoint.Name]; exists {
		s.logger.Infof("endpoint [%s] already exists - replacing", oldEndpoint.Name)
		oldEndpoint.Shutdown()
	}

	newEndpoint.Start()

	s.endpoints[newEndpoint.Name] = newEndpoint
	logging.Infof("registered HTTP server endpoint [%s]", newEndpoint.Name)

	return nil
}

func (s *service) SendAction(sendAction *actions.SendAction) error {
	endpoint, exists := s.endpoints[sendAction.EndpointName]
	if !exists {
		s.logger.Errorf("HTTP server endpoint [%s] not found - action [%s] will not be executed",
			sendAction.EndpointName, sendAction.Name)
	}

	return endpoint.Send(sendAction)
}

func (s *service) ReceiveAction(receiveAction *actions.ReceiveAction) (*http.Request, error) {
	endpoint, exists := s.endpoints[receiveAction.EndpointName]
	if !exists {
		s.logger.Errorf("HTTP server endpoint [%s] not found - action [%s] will not be executed",
			receiveAction.EndpointName, receiveAction.Name)
	}

	return endpoint.Receive(receiveAction)
}
