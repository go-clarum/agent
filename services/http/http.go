package http

import (
	"context"
	"fmt"
	. "github.com/go-clarum/agent/api/http"
	"github.com/go-clarum/agent/logging"
	"github.com/go-clarum/agent/services/http/internal/client"
	"google.golang.org/grpc"
)

var httpClientService = client.NewHttpClientService()

type grpcService struct {
	UnimplementedHttpServiceServer
}

func RegisterHttpService(server *grpc.Server) {
	logging.Infof("Registering HttpService")
	RegisterHttpServiceServer(server, &grpcService{})
}

func (s *grpcService) InitClientEndpoint(ctx context.Context, ic *InitializeClientRequest) (*InitializeClientResponse, error) {
	err := httpClientService.InitializeEndpoint(ic)

	return &InitializeClientResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}

func (s *grpcService) InitServerEndpoint(ctx context.Context, is *InitializeServerRequest) (*InitializeServerResponse, error) {
	return &InitializeServerResponse{}, nil
}

func (s *grpcService) ClientSendAction(ctx context.Context, sendAction *ClientSendActionRequest) (*ClientSendActionResponse, error) {
	err := httpClientService.SendAction(sendAction)

	return &ClientSendActionResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}

func (s *grpcService) ClientReceiveAction(ctx context.Context, receiveAction *ClientReceiveActionRequest) (*ClientReceiveActionResponse, error) {
	_, err := httpClientService.ReceiveAction(receiveAction)

	return &ClientReceiveActionResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}

func (s *grpcService) ServerSendAction(ctx context.Context, sendAction *ServerSendActionRequest) (*ServerSendActionResponse, error) {
	return &ServerSendActionResponse{}, nil
}

func (s *grpcService) ServerReceiveAction(ctx context.Context, receiveAction *ServerReceiveActionRequest) (*ServerReceiveActionResponse, error) {
	return &ServerReceiveActionResponse{}, nil
}
