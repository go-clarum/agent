package agent

import (
	"context"
	. "github.com/go-clarum/agent/api/agent"
	"github.com/go-clarum/agent/logging"
	"google.golang.org/grpc"
)

type service struct {
	UnimplementedAgentServiceServer
}

func RegisterAgentService(server *grpc.Server) {
	logging.Infof("Registering AgentService")
	RegisterAgentServiceServer(server, &service{})
}

func (s *service) GetStatus(ctx context.Context, request *StatusRequest) (*StatusResponse, error) {
	return &StatusResponse{}, nil
}
func (s *service) Shutdown(ctx context.Context, request *ShutdownRequest) (*ShutdownResponse, error) {
	return &ShutdownResponse{}, nil
}
