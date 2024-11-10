package interfaces

import (
	client "github.com/go-clarum/agent/application/services/http/client/actions"
	"net/http"
)

type ClientService interface {
	InitializeEndpoint(req *client.InitEndpointAction) error
	SendAction(sendAction *client.SendAction) error
	ReceiveAction(receiveAction *client.ReceiveAction) (*http.Response, error)
}
