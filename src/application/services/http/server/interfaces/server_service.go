package interfaces

import (
	"github.com/go-clarum/agent/application/services/http/server/actions"
	"net/http"
)

type ServerService interface {
	InitializeEndpoint(is *actions.InitEndpointAction) error
	SendAction(sendAction *actions.SendAction) error
	ReceiveAction(receiveAction *actions.ReceiveAction) (*http.Request, error)
}
