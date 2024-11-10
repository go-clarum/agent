package client

import (
	"github.com/go-clarum/agent/application/services/http/client/actions"
	"github.com/go-clarum/agent/application/services/http/client/interfaces"
	"github.com/go-clarum/agent/application/services/http/client/internal"
	"github.com/go-clarum/agent/infrastructure/logging"
	"net/http"
)

type service struct {
	endpoints map[string]*internal.Endpoint
	logger    *logging.Logger
}

func NewHttpClientService() interfaces.ClientService {
	return &service{
		endpoints: make(map[string]*internal.Endpoint),
		logger:    logging.NewLogger("HttpClientService"),
	}
}

func (s *service) InitializeEndpoint(req *actions.InitEndpointAction) error {
	newEndpoint, err := internal.NewEndpoint(req)

	if err != nil {
		s.logger.Errorf("failed to initialize HTTP client endpoint - %s", err)
		return err
	}

	if oldEndpoint, exists := s.endpoints[newEndpoint.Name]; exists {
		s.logger.Infof("HTTP client endpoint [%s] already exists - replacing", oldEndpoint.Name)
	}

	s.endpoints[newEndpoint.Name] = newEndpoint
	logging.Infof("registered HTTP client endpoint [%s]", newEndpoint.Name)

	return nil
}

func (s *service) SendAction(sendAction *actions.SendAction) error {
	endpoint, exists := s.endpoints[sendAction.EndpointName]
	if !exists {
		s.logger.Errorf("HTTP client endpoint [%s] not found - action [%s] will not be executed",
			sendAction.EndpointName, sendAction.Name)
	}

	return endpoint.Send(sendAction)
}

func (s *service) ReceiveAction(receiveAction *actions.ReceiveAction) (*http.Response, error) {
	endpoint, exists := s.endpoints[receiveAction.EndpointName]
	if !exists {
		s.logger.Errorf("HTTP client endpoint [%s] not found - action [%s] will not be executed",
			receiveAction.EndpointName, receiveAction.Name)
	}

	return endpoint.Receive(receiveAction)
}
