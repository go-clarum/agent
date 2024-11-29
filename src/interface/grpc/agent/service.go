package agent

import (
	"context"
	"github.com/go-clarum/agent/application/services/agent"
	"github.com/go-clarum/agent/infrastructure/config"
	"github.com/go-clarum/agent/infrastructure/logging"
	"github.com/go-clarum/agent/interface/grpc/agent/internal/api"
	"github.com/go-clarum/agent/interface/grpc/agent/internal/api/commands/logs"
	"google.golang.org/grpc"
)

var agentService = agent.NewAgentService()

type grpcService struct {
	api.UnimplementedAgentApiServer
}

func RegisterAgentApi(server *grpc.Server) {
	logging.Infof("registering AgentService")
	api.RegisterAgentApiServer(server, &grpcService{})
}

func (s *grpcService) Status(ctx context.Context, request *api.StatusRequest) (*api.StatusResponse, error) {
	logging.Infof("signaling agent status")

	return &api.StatusResponse{
		Version: config.Version(),
	}, nil
}

func (s *grpcService) Shutdown(ctx context.Context, request *api.ShutdownRequest) (*api.ShutdownResponse, error) {
	logging.Infof("received shutdown command")

	defer agentService.Shutdown()
	return &api.ShutdownResponse{}, nil
}

func (s *grpcService) Logs(req *logs.LogsRequest, stream api.AgentApi_LogsServer) error {
	logging.Debugf("log listener %s connected", req.ListenerName)
	logChannel := logging.LogEmitter.Subscribe()
	defer logging.LogEmitter.Unsubscribe(logChannel)

	for {
		select {
		case <-stream.Context().Done():
			logging.Debugf("closing log stream for %s ", req.ListenerName)
			return nil
		case message := <-logChannel:
			if err := stream.Send(&logs.LogEntry{Message: message}); err != nil {
				logging.Errorf("received error while streaming logs: %s", err.Error())
			}
		}
	}
}
