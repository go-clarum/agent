package cmd

import (
	"context"
	"fmt"
	. "github.com/go-clarum/agent/api/cmd"
	"github.com/go-clarum/agent/logging"
	"github.com/go-clarum/agent/services/cmd/internal/service"
	"google.golang.org/grpc"
)

var commandService = service.NewCommandService()

type grpcService struct {
	UnimplementedCmdServiceServer
}

func RegisterCmdService(server *grpc.Server) {
	logging.Infof("Registering CommandService")
	RegisterCmdServiceServer(server, &grpcService{})
}

func (s *grpcService) InitEndpoint(ctx context.Context, req *InitEndpointRequest) (*InitEndpointResponse, error) {
	err := commandService.InitializeEndpoint(req.Name, req.CmdComponents, req.WarmupMillis)

	return &InitEndpointResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}

func (s *grpcService) ShutdownEndpoint(ctx context.Context, req *ShutdownEndpointRequest) (*ShutdownEndpointResponse, error) {
	err := commandService.ShutdownEndpoint(req.Name)

	return &ShutdownEndpointResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}
