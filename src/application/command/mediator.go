package command

import (
	"github.com/go-clarum/agent/application/command/common"
	httpClient "github.com/go-clarum/agent/application/command/http/client"
	httpServer "github.com/go-clarum/agent/application/command/http/server"
	"github.com/go-clarum/agent/infrastructure/logging"
)

var med *mediator

type mediator struct {
	logger   *logging.Logger
	handlers []common.CommandHandler
}

func GetMediator() *mediator {
	if med == nil {
		med = &mediator{
			logger: logging.NewLogger("commandMediator"),
		}
	}

	med.handlers = append(med.handlers, httpClient.NewHttpClientHandler())
	med.handlers = append(med.handlers, httpServer.NewHttpServerHandler())

	return med
}

func (m *mediator) DelegateCommand(command any) any {
	for _, h := range m.handlers {
		if h.CanHandle(command) {
			return h.Handle(command)
		}
	}

	m.logger.Errorf("unable to handle command: %v", command)
	return nil
}
