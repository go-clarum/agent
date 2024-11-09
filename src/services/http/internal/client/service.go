package client

import (
	"github.com/go-clarum/agent/logging"
	"net/http"
)

type ClientService interface {
	InitializeEndpoint(req *InitRequest) error
	SendAction(sendAction *SendAction) error
	ReceiveAction(receiveAction *ReceiveAction) (*http.Response, error)
}

type service struct {
	endpoints map[string]*endpoint
	logger    *logging.Logger
}

func NewHttpClientService() ClientService {
	return &service{
		endpoints: make(map[string]*endpoint),
		logger:    logging.NewLogger("HttpClientService"),
	}
}

func (s *service) InitializeEndpoint(req *InitRequest) error {
	newEndpoint, err := newEndpoint(req)

	if err != nil {
		s.logger.Errorf("failed to initialize HTTP client endpoint - %s", err)
		return err
	}

	if oldEndpoint, exists := s.endpoints[newEndpoint.name]; exists {
		s.logger.Infof("HTTP client endpoint [%s] already exists - replacing", oldEndpoint.name)
	}

	s.endpoints[newEndpoint.name] = newEndpoint
	logging.Infof("registered HTTP client endpoint [%s]", newEndpoint.name)

	return nil
}

func (s *service) SendAction(sendAction *SendAction) error {
	endpoint, exists := s.endpoints[sendAction.EndpointName]
	if !exists {
		s.logger.Errorf("HTTP client endpoint [%s] not found - action [%s] will not be executed",
			sendAction.EndpointName, sendAction.Name)
	}

	return endpoint.send(sendAction)
}

func (s *service) ReceiveAction(receiveAction *ReceiveAction) (*http.Response, error) {
	endpoint, exists := s.endpoints[receiveAction.EndpointName]
	if !exists {
		s.logger.Errorf("HTTP client endpoint [%s] not found - action [%s] will not be executed",
			receiveAction.EndpointName, receiveAction.Name)
	}

	return endpoint.receive(receiveAction)
}
