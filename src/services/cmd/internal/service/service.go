package service

import (
	"errors"
	"fmt"
	"github.com/go-clarum/agent/logging"
)

type CommandService interface {
	InitializeEndpoint(name string, cmdComponents []string, warmupMillis int32) error
	ShutdownEndpoint(name string) error
}

type service struct {
	endpoints map[string]*endpoint
	logger    *logging.Logger
}

func NewCommandService() CommandService {
	return &service{
		endpoints: make(map[string]*endpoint),
		logger:    logging.NewLogger("CommandService"),
	}
}

func (s *service) InitializeEndpoint(name string, cmdComponents []string, warmupMillis int32) error {
	newEndpoint, err := newEndpoint(name, cmdComponents, warmupMillis)

	if err != nil {
		s.logger.Errorf("failed to initialize endpoint - %s", err)
		return err
	}

	if oldEndpoint, exists := s.endpoints[newEndpoint.name]; exists {
		s.logger.Infof("endpoint [%s] already exists - replacing", oldEndpoint.name)
		go func() {
			_ = oldEndpoint.shutdown()
		}()
	}

	if err := newEndpoint.start(); err != nil {
		s.logger.Errorf("failed to initialize endpoint - %s", err)
		return err
	}

	s.endpoints[newEndpoint.name] = newEndpoint
	logging.Infof("registered endpoint [%s]", newEndpoint.name)

	return nil
}

func (s *service) ShutdownEndpoint(name string) error {
	if endpoint, exists := s.endpoints[name]; exists {
		s.logger.Infof("shutting down endpoint [%s]", endpoint.name)
		err := endpoint.shutdown()
		if err != nil {
			s.logger.Errorf("error during endpoint shutdown - %s", err)
			return err
		}
		return nil
	}

	err := errors.New(fmt.Sprintf("endpoint [%s] does not exist", name))
	s.logger.Errorf("unable to shutdown endpoint - %s", err)
	return err
}
