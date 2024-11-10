package cmd

import (
	"context"
	"fmt"
	"github.com/go-clarum/agent/application/services/cmd"
	"github.com/go-clarum/agent/infrastructure/logging"
	"github.com/go-clarum/agent/interface/grpc/cmd/internal/api"
	"google.golang.org/grpc"
)

var commandService = cmd.NewCommandService()

type grpcService struct {
	api.UnimplementedCmdApiServer
}

func RegisterCmdService(server *grpc.Server) {
	logging.Infof("registering CommandService")
	api.RegisterCmdApiServer(server, &grpcService{})
}

func (s *grpcService) InitEndpoint(ctx context.Context, req *api.InitEndpointRequest) (*api.InitEndpointResponse, error) {
	err := commandService.InitializeEndpoint(req.Name, req.CmdComponents, req.WarmupMillis)

	return &api.InitEndpointResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}

func (s *grpcService) ShutdownEndpoint(ctx context.Context, req *api.ShutdownEndpointRequest) (*api.ShutdownEndpointResponse, error) {
	err := commandService.ShutdownEndpoint(req.Name)

	return &api.ShutdownEndpointResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}
