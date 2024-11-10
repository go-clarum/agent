package agent

import (
	"github.com/go-clarum/agent/application/services/agent/interfaces"
	"github.com/go-clarum/agent/infrastructure/config"
	"github.com/go-clarum/agent/infrastructure/logging"
	"os"
)

type service struct {
	logger *logging.Logger
}

func NewAgentService() interfaces.AgentService {
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
