package agent

import (
	"context"
	. "github.com/go-clarum/agent/api/agent"
	"github.com/go-clarum/agent/config"
	"github.com/go-clarum/agent/logging"
	"github.com/go-clarum/agent/services/agent/internal/service"
	"google.golang.org/grpc"
)

var agentService = service.NewAgentService()

type grpcService struct {
	UnimplementedAgentServiceServer
}

func RegisterAgentService(server *grpc.Server) {
	logging.Infof("Registering AgentService")
	RegisterAgentServiceServer(server, &grpcService{})
}

func (s *grpcService) Status(ctx context.Context, request *StatusRequest) (*StatusResponse, error) {
	return &StatusResponse{
		Version: config.Version(),
	}, nil
}

func (s *grpcService) Shutdown(ctx context.Context, request *ShutdownRequest) (*ShutdownResponse, error) {
	defer agentService.Shutdown()
	return &ShutdownResponse{}, nil
}
