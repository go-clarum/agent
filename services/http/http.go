package http

import (
	"context"
	"fmt"
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
	logging.Infof("Registering HttpService")
	RegisterHttpServiceServer(server, &grpcService{})
}

func (s *grpcService) InitClientEndpoint(ctx context.Context, ic *InitClientRequest) (*InitClientResponse, error) {
	ir := newClientInitRequestFrom(ic)
	err := httpClientService.InitializeEndpoint(ir)

	return &InitClientResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}

func (s *grpcService) InitServerEndpoint(ctx context.Context, is *InitServerRequest) (*InitServerResponse, error) {
	req := newServerInitRequestFrom(is)
	err := httpServerService.InitializeEndpoint(req)

	return &InitServerResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}

func (s *grpcService) ClientSendAction(ctx context.Context, sendAction *ClientSendActionRequest) (*ClientSendActionResponse, error) {
	sa := newClientSendActionFrom(sendAction)
	err := httpClientService.SendAction(sa)

	return &ClientSendActionResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}

func (s *grpcService) ClientReceiveAction(ctx context.Context, receiveAction *ClientReceiveActionRequest) (*ClientReceiveActionResponse, error) {
	ra := newClientReceiveActionFrom(receiveAction)
	_, err := httpClientService.ReceiveAction(ra)

	return &ClientReceiveActionResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}

func (s *grpcService) ServerSendAction(ctx context.Context, sendAction *ServerSendActionRequest) (*ServerSendActionResponse, error) {
	sa := newServerSendActionFrom(sendAction)
	err := httpServerService.SendAction(sa)

	return &ServerSendActionResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}

func (s *grpcService) ServerReceiveAction(ctx context.Context, receiveAction *ServerReceiveActionRequest) (*ServerReceiveActionResponse, error) {
	ra := newServerReceiveActionFrom(receiveAction)
	_, err := httpServerService.ReceiveAction(ra)

	return &ServerReceiveActionResponse{
		Error: fmt.Sprintf("%s", err),
	}, nil
}
