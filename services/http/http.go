package http

import (
	"context"
	. "github.com/go-clarum/agent/api/http"
	"github.com/go-clarum/agent/logging"
	"github.com/go-clarum/agent/services/http/internal/client"
	"github.com/go-clarum/agent/services/http/internal/server"
	"google.golang.org/grpc"
)

var httpClientService = client.NewHttpClientService()
var httpServerService = server.NewHttpServerService()

type grpcService struct {
	UnimplementedHttpServiceServer
}

func RegisterHttpService(server *grpc.Server) {
	logging.Infof("registering HttpService")
	RegisterHttpServiceServer(server, &grpcService{})
}

func (s *grpcService) InitClientEndpoint(ctx context.Context, ic *InitClientRequest) (*InitClientResponse, error) {
	var errorMessage string

	ir := newClientInitRequestFrom(ic)
	if err := httpClientService.InitializeEndpoint(ir); err != nil {
		errorMessage = err.Error()
	}

	return &InitClientResponse{
		Error: errorMessage,
	}, nil
}

func (s *grpcService) InitServerEndpoint(ctx context.Context, is *InitServerRequest) (*InitServerResponse, error) {
	var errorMessage string

	req := newServerInitRequestFrom(is)
	if err := httpServerService.InitializeEndpoint(req); err != nil {
		errorMessage = err.Error()
	}

	return &InitServerResponse{
		Error: errorMessage,
	}, nil
}

func (s *grpcService) ClientSendAction(ctx context.Context, sendAction *ClientSendActionRequest) (*ClientSendActionResponse, error) {
	var errorMessage string

	sa := newClientSendActionFrom(sendAction)
	if err := httpClientService.SendAction(sa); err != nil {
		errorMessage = err.Error()
	}

	return &ClientSendActionResponse{
		Error: errorMessage,
	}, nil
}

func (s *grpcService) ClientReceiveAction(ctx context.Context, receiveAction *ClientReceiveActionRequest) (*ClientReceiveActionResponse, error) {
	var errorMessage string

	ra := newClientReceiveActionFrom(receiveAction)
	if _, err := httpClientService.ReceiveAction(ra); err != nil {
		errorMessage = err.Error()
	}

	return &ClientReceiveActionResponse{
		Error: errorMessage,
	}, nil
}

func (s *grpcService) ServerSendAction(ctx context.Context, sendAction *ServerSendActionRequest) (*ServerSendActionResponse, error) {
	var errorMessage string

	sa := newServerSendActionFrom(sendAction)
	if err := httpServerService.SendAction(sa); err != nil {
		errorMessage = err.Error()
	}

	return &ServerSendActionResponse{
		Error: errorMessage,
	}, nil
}

func (s *grpcService) ServerReceiveAction(ctx context.Context, receiveAction *ServerReceiveActionRequest) (*ServerReceiveActionResponse, error) {
	var errorMessage string

	ra := newServerReceiveActionFrom(receiveAction)
	if _, err := httpServerService.ReceiveAction(ra); err != nil {
		errorMessage = err.Error()
	}

	return &ServerReceiveActionResponse{
		Error: errorMessage,
	}, nil
}
