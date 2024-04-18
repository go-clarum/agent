package http

import (
	"context"
	. "github.com/go-clarum/agent/api/http"
	"github.com/go-clarum/agent/logging"
	"google.golang.org/grpc"
)

type service struct {
	UnimplementedHttpServiceServer
}

func RegisterHttpService(server *grpc.Server) {
	logging.Infof("Registering HttpService")
	RegisterHttpServiceServer(server, &service{})
}

func (s *service) CreateClientEndpoint(ctx context.Context, ce *ClientEndpoint) (*CreateResponse, error) {
	return &CreateResponse{}, nil
}
