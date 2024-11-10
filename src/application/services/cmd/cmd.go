package cmd

import (
	"errors"
	"fmt"
	"github.com/go-clarum/agent/application/services/cmd/interfaces"
	"github.com/go-clarum/agent/application/services/cmd/internal"
	"github.com/go-clarum/agent/infrastructure/logging"
)

type service struct {
	endpoints map[string]*internal.Endpoint
	logger    *logging.Logger
}

func NewCommandService() interfaces.CommandService {
	return &service{
		endpoints: make(map[string]*internal.Endpoint),
		logger:    logging.NewLogger("CommandService"),
	}
}

func (s *service) InitializeEndpoint(name string, cmdComponents []string, warmupMillis int32) error {
	newEndpoint, err := internal.NewEndpoint(name, cmdComponents, warmupMillis)

	if err != nil {
		s.logger.Errorf("failed to initialize endpoint - %s", err)
		return err
	}

	if oldEndpoint, exists := s.endpoints[newEndpoint.Name]; exists {
		s.logger.Infof("endpoint [%s] already exists - replacing", oldEndpoint.Name)
		go func() {
			_ = oldEndpoint.Shutdown()
		}()
	}

	if err := newEndpoint.Start(); err != nil {
		s.logger.Errorf("failed to initialize endpoint - %s", err)
		return err
	}

	s.endpoints[newEndpoint.Name] = newEndpoint
	logging.Infof("registered endpoint [%s]", newEndpoint.Name)

	return nil
}

func (s *service) ShutdownEndpoint(name string) error {
	if endpoint, exists := s.endpoints[name]; exists {
		s.logger.Infof("shutting down endpoint [%s]", endpoint.Name)
		err := endpoint.Shutdown()
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
