package server

import (
	"github.com/go-clarum/agent/logging"
	"net/http"
)

type ServerEndpoint interface {
	InitializeEndpoint(is *initRequest) error
	SendAction(sendAction *sendAction) error
	ReceiveAction(receiveAction *receiveAction) (*http.Request, error)
}

type service struct {
	endpoints map[string]*endpoint
	logger    *logging.Logger
}

func NewHttpServerService() *service {
	return &service{
		endpoints: make(map[string]*endpoint),
		logger:    logging.NewLogger("HttpServerService"),
	}
}

func (s *service) InitializeEndpoint(is *initRequest) error {
	newEndpoint := newEndpoint(is)

	if oldEndpoint, exists := s.endpoints[newEndpoint.name]; exists {
		s.logger.Infof("endpoint [%s] already exists - replacing", oldEndpoint.name)
		oldEndpoint.shutdown()
	}

	newEndpoint.start()

	s.endpoints[newEndpoint.name] = newEndpoint
	logging.Infof("registered endpoint [%s]", newEndpoint.name)

	return nil
}

func (s *service) SendAction(sendAction *sendAction) error {
	endpoint, exists := s.endpoints[sendAction.endpointName]
	if !exists {
		s.logger.Errorf("HTTP server endpoint [%s] not found - action [%s] will not be executed",
			sendAction.endpointName, sendAction.name)
	}

	return endpoint.send(sendAction)
}

func (s *service) ReceiveAction(receiveAction *receiveAction) (*http.Request, error) {
	endpoint, exists := s.endpoints[receiveAction.endpointName]
	if !exists {
		s.logger.Errorf("HTTP server endpoint [%s] not found - action [%s] will not be executed",
			receiveAction.endpointName, receiveAction.name)
	}

	return endpoint.receive(receiveAction)
}
