package service

import (
	"github.com/go-clarum/agent/config"
	"github.com/go-clarum/agent/logging"
	"os"
)

type AgentService interface {
	Status() string
	Shutdown()
}

type service struct {
	logger *logging.Logger
}

func NewAgentService() AgentService {
	return &service{
		logger: logging.NewLogger("AgentService"),
	}
}

func (s service) Status() string {
	return config.Version()
}

func (s service) Shutdown() {
	s.logger.Info("received shutdown signal from binding")
	os.Exit(0)
}
