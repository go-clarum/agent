package server

import (
	"github.com/go-clarum/agent/logging"
	"net/http"
)

type ServerService interface {
	InitializeEndpoint(is *InitRequest) error
	SendAction(sendAction *SendAction) error
	ReceiveAction(receiveAction *ReceiveAction) (*http.Request, error)
}

type service struct {
	endpoints map[string]*endpoint
	logger    *logging.Logger
}

func NewHttpServerService() ServerService {
	return &service{
		endpoints: make(map[string]*endpoint),
		logger:    logging.NewLogger("HttpServerService"),
	}
}

func (s *service) InitializeEndpoint(is *InitRequest) error {
	newEndpoint := newEndpoint(is)

	if oldEndpoint, exists := s.endpoints[newEndpoint.name]; exists {
		s.logger.Infof("endpoint [%s] already exists - replacing", oldEndpoint.name)
		oldEndpoint.shutdown()
	}

	newEndpoint.start()

	s.endpoints[newEndpoint.name] = newEndpoint
	logging.Infof("registered HTTP server endpoint [%s]", newEndpoint.name)

	return nil
}

func (s *service) SendAction(sendAction *SendAction) error {
	endpoint, exists := s.endpoints[sendAction.EndpointName]
	if !exists {
		s.logger.Errorf("HTTP server endpoint [%s] not found - action [%s] will not be executed",
			sendAction.EndpointName, sendAction.Name)
	}

	return endpoint.send(sendAction)
}

func (s *service) ReceiveAction(receiveAction *ReceiveAction) (*http.Request, error) {
	endpoint, exists := s.endpoints[receiveAction.EndpointName]
	if !exists {
		s.logger.Errorf("HTTP server endpoint [%s] not found - action [%s] will not be executed",
			receiveAction.EndpointName, receiveAction.Name)
	}

	return endpoint.receive(receiveAction)
}
