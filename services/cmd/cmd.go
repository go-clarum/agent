package cmd

import (
	"context"
	"fmt"
	. "github.com/go-clarum/agent/api/cmd"
	"github.com/go-clarum/agent/logging"
	"github.com/go-clarum/agent/services/cmd/internal/service"
	"google.golang.org/grpc"
)

type cmdService struct {
	UnimplementedCmdServiceServer
}

func RegisterCmdService(server *grpc.Server) {
	logging.Infof("Registering CommandService")
	RegisterCmdServiceServer(server, &cmdService{})
}

func (s *cmdService) InitEndpoint(ctx context.Context, req *InitCmdEndpoint) (*InitCmdEndpointResponse, error) {
	err := service.InitializeEndpoint(req.Name, req.CmdComponents, req.WarmupSeconds)

	return &InitCmdEndpointResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}

func (s *cmdService) ShutdownEndpoint(ctx context.Context, req *ShutdownCmdEndpoint) (*ShutdownCmdEndpointResponse, error) {
	err := service.ShutdownEndpoint(req.Name)

	return &ShutdownCmdEndpointResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}
