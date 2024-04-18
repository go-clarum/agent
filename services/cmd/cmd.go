package cmd

import (
	"context"
	. "github.com/go-clarum/agent/api/cmd"
	"github.com/go-clarum/agent/logging"
	"google.golang.org/grpc"
)

type service struct {
	UnimplementedCmdServiceServer
}

func RegisterCmdService(server *grpc.Server) {
	logging.Infof("Registering CommandService")
	RegisterCmdServiceServer(server, &service{})
}

func (s *service) CreateCommandEndpoint(ctx context.Context, ce *CommandEndpoint) (*CommandCreateResponse, error) {
	return &CommandCreateResponse{}, nil
}
