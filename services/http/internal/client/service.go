package client

import (
	"github.com/go-clarum/agent/logging"
	"net/http"
)

type ClientEndpoint interface {
	InitializeEndpoint(req *initRequest) error
	SendAction(sendAction *sendAction) error
	ReceiveAction(receiveAction *receiveAction) (*http.Response, error)
}

type service struct {
	endpoints map[string]*endpoint
	logger    *logging.Logger
}

func NewHttpClientService() *service {
	return &service{
		endpoints: make(map[string]*endpoint),
		logger:    logging.NewLogger("HttpClientService"),
	}
}

func (s *service) InitializeEndpoint(req *initRequest) error {
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

func (s *service) SendAction(sendAction *sendAction) error {
	endpoint, exists := s.endpoints[sendAction.endpointName]
	if !exists {
		s.logger.Errorf("HTTP client endpoint [%s] not found - action [%s] will not be executed",
			sendAction.endpointName, sendAction.name)
	}

	return endpoint.send(sendAction)
}

func (s *service) ReceiveAction(receiveAction *receiveAction) (*http.Response, error) {
	endpoint, exists := s.endpoints[receiveAction.endpointName]
	if !exists {
		s.logger.Errorf("HTTP client endpoint [%s] not found - action [%s] will not be executed",
			receiveAction.endpointName, receiveAction.name)
	}

	return endpoint.receive(receiveAction)
}
