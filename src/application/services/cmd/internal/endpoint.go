package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-clarum/agent/application/control"
	"github.com/go-clarum/agent/application/utils/durations"
	clarumstrings "github.com/go-clarum/agent/application/validators/strings"
	"github.com/go-clarum/agent/infrastructure/logging"
	"os/exec"
	"time"
)

type Endpoint struct {
	Name          string
	cmdComponents []string
	warmup        time.Duration
	cmd           *exec.Cmd
	cmdCancel     context.CancelFunc
	logger        *logging.Logger
}

func NewEndpoint(name string, cmdComponents []string, warmupMillis int32) (*Endpoint, error) {
	if clarumstrings.IsBlank(name) {
		return nil, errors.New("cannot create command Endpoint - Name is empty")
	}

	if len(cmdComponents) == 0 || clarumstrings.IsBlank(cmdComponents[0]) {
		return nil, errors.New(fmt.Sprintf("cannot create command Endpoint [%s] - cmd is empty", name))
	}

	warmupDuration := time.Duration(warmupMillis)

	return &Endpoint{
		Name:          name,
		cmdComponents: cmdComponents,
		warmup:        durations.GetDurationWithDefault(warmupDuration, 1*time.Millisecond),
		logger:        logging.NewLogger(loggerName(name)),
	}, nil
}

// Start the process from the given command & arguments.
// The process will be started into a cancelable context so that we can
// cancel it later in the post-integration test phase.
func (endpoint *Endpoint) Start() error {
	endpoint.logger.Infof("running cmd [%s]", endpoint.cmdComponents)
	ctx, cancel := context.WithCancel(context.Background())

	endpoint.cmd = exec.CommandContext(ctx, endpoint.cmdComponents[0], endpoint.cmdComponents[1:]...)
	endpoint.cmdCancel = cancel

	endpoint.logger.Debug("starting command")
	if err := endpoint.cmd.Start(); err != nil {
		return err
	} else {
		endpoint.logger.Debug("cmd start successful")
	}

	time.Sleep(endpoint.warmup)
	endpoint.logger.Debug("warmup ended")

	return nil
}

// shutdown the running process. Since the process was created with a context, we will attempt to
// call ctx.Cancel(). If it returns an error, the process will be killed just in case.
// We also wait for the action here, so that the post-integration test phase ends successfully.
// TODO: check this code again, some parts are redundant
func (endpoint *Endpoint) Shutdown() error {
	control.RunningActions.Add(1)
	defer control.RunningActions.Done()

	endpoint.logger.Infof("stopping cmd [%s]", endpoint.cmdComponents)

	if endpoint.cmdCancel != nil {
		endpoint.logger.Debug("cancelling cmd")
		endpoint.cmdCancel()

		if _, err := endpoint.cmd.Process.Wait(); err != nil {
			endpoint.logger.Errorf("cmd.Wait() returned error - [%s]", err)
			endpoint.killProcess()
			return err
		} else {
			endpoint.logger.Debug("context cancel finished successfully")
		}
	} else {
		if err := endpoint.cmd.Process.Release(); err != nil {
			endpoint.logger.Errorf("cmd.Release() returned error - [%s]", err)
			endpoint.killProcess()
			return err
		} else {
			endpoint.logger.Debug("cmd kill successful")
		}
	}

	return nil
}

func (endpoint *Endpoint) killProcess() {
	endpoint.logger.Info("killing process")

	if err := endpoint.cmd.Process.Kill(); err != nil {
		endpoint.logger.Errorf("cmd.Kill() returned error - [%s]", err)
		return
	}
}

func loggerName(cmdName string) string {
	return fmt.Sprintf("Command %s", cmdName)
}
