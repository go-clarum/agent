package service

import (
	"errors"
	"fmt"
	"github.com/go-clarum/agent/logging"
)

var endpoints map[string]*endpoint
var logger *logging.Logger

func init() {
	endpoints = make(map[string]*endpoint)
	logger = logging.NewLogger("CommandService")
}

func InitializeEndpoint(name string, cmdComponents []string, warmupSeconds int32) error {
	newEndpoint, err := newCommandEndpoint(name, cmdComponents, warmupSeconds)

	if err != nil {
		logger.Errorf("failed to initialize endpoint - %s", err)
		return err
	}

	if oldEndpoint, exists := endpoints[newEndpoint.name]; exists {
		logger.Infof("endpoint [%s] already exists - replacing", oldEndpoint.name)
		go func() {
			_ = oldEndpoint.shutdown()
		}()
	}

	if err := newEndpoint.start(); err != nil {
		logger.Errorf("failed to initialize endpoint - %s", err)
		return err
	}

	endpoints[newEndpoint.name] = newEndpoint
	logging.Debugf("registered endpoint [%s]", newEndpoint.name)

	return nil
}

func ShutdownEndpoint(name string) error {
	if endpoint, exists := endpoints[name]; exists {
		logger.Infof("shutting down endpoint [%s]", endpoint.name)
		err := endpoint.shutdown()
		if err != nil {
			logger.Errorf("error during endpoint shutdown - %s", err)
			return err
		}
		return nil
	}

	err := errors.New(fmt.Sprintf("endpoint [%s] does not exist", name))
	logger.Errorf("unable to shutdown endpoint - %s", err)
	return err
}
