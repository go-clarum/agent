package http

import (
	"context"
	"github.com/go-clarum/agent/application/services/http/client"
	"github.com/go-clarum/agent/application/services/http/server"
	"github.com/go-clarum/agent/infrastructure/logging"
	"github.com/go-clarum/agent/interface/grpc/http/internal/api"
	"github.com/go-clarum/agent/interface/grpc/http/mappers"

	"google.golang.org/grpc"
)

var httpClientService = client.NewHttpClientService()
var httpServerService = server.NewHttpServerService()

type grpcService struct {
	api.UnimplementedHttpApiServer
}

func RegisterHttpService(server *grpc.Server) {
	logging.Infof("registering HttpService")
	api.RegisterHttpApiServer(server, &grpcService{})
}

func (s *grpcService) InitClientEndpoint(ctx context.Context, ic *api.InitClientRequest) (*api.InitClientResponse, error) {
	var errorMessage string

	ir := mappers.NewClientInitRequestFrom(ic)
	if err := httpClientService.InitializeEndpoint(ir); err != nil {
		errorMessage = err.Error()
	}

	return &api.InitClientResponse{
		Error: errorMessage,
	}, nil
}

func (s *grpcService) InitServerEndpoint(ctx context.Context, is *api.InitServerRequest) (*api.InitServerResponse, error) {
	var errorMessage string

	req := mappers.NewServerInitRequestFrom(is)
	if err := httpServerService.InitializeEndpoint(req); err != nil {
		errorMessage = err.Error()
	}

	return &api.InitServerResponse{
		Error: errorMessage,
	}, nil
}

func (s *grpcService) ClientSendAction(ctx context.Context, sendAction *api.ClientSendActionRequest) (*api.ClientSendActionResponse, error) {
	var errorMessage string

	sa := mappers.NewClientSendActionFrom(sendAction)
	if err := httpClientService.SendAction(sa); err != nil {
		errorMessage = err.Error()
	}

	return &api.ClientSendActionResponse{
		Error: errorMessage,
	}, nil
}

func (s *grpcService) ClientReceiveAction(ctx context.Context, receiveAction *api.ClientReceiveActionRequest) (*api.ClientReceiveActionResponse, error) {
	var errorMessage string

	ra := mappers.NewClientReceiveActionFrom(receiveAction)
	if _, err := httpClientService.ReceiveAction(ra); err != nil {
		errorMessage = err.Error()
	}

	return &api.ClientReceiveActionResponse{
		Error: errorMessage,
	}, nil
}

func (s *grpcService) ServerSendAction(ctx context.Context, sendAction *api.ServerSendActionRequest) (*api.ServerSendActionResponse, error) {
	var errorMessage string

	sa := mappers.NewServerSendActionFrom(sendAction)
	if err := httpServerService.SendAction(sa); err != nil {
		errorMessage = err.Error()
	}

	return &api.ServerSendActionResponse{
		Error: errorMessage,
	}, nil
}

func (s *grpcService) ServerReceiveAction(ctx context.Context, receiveAction *api.ServerReceiveActionRequest) (*api.ServerReceiveActionResponse, error) {
	var errorMessage string

	ra := mappers.NewServerReceiveActionFrom(receiveAction)
	if _, err := httpServerService.ReceiveAction(ra); err != nil {
		errorMessage = err.Error()
	}

	return &api.ServerReceiveActionResponse{
		Error: errorMessage,
	}, nil
}
