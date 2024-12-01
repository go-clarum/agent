package agent

import (
	"context"
	"github.com/go-clarum/agent/application/command"
	"github.com/go-clarum/agent/application/services/agent"
	"github.com/go-clarum/agent/infrastructure/config"
	"github.com/go-clarum/agent/infrastructure/logging"
	"github.com/go-clarum/agent/interface/grpc/agent/internal/api"
	"github.com/go-clarum/agent/interface/grpc/agent/internal/mapper"
	"google.golang.org/grpc"
	"io"
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

func (s *grpcService) Session(stream grpc.BidiStreamingServer[api.ActionCommand, api.CommandResponse]) error {
	logging.Debug("session started")

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		var inCommand = mapper.TranslateCommand(in)
		result := command.GetMediator().DelegateCommand(inCommand)
		var outResult = mapper.TranslateResult(result)

		if err := stream.Send(outResult); err != nil {
			return err
		}
	}

	logging.Debug("session closed")
	return nil
}
