package server

import (
	"fmt"
	"github.com/go-clarum/agent/services/http/internal"
	"time"
)

type InitRequest struct {
	Name           string
	Port           uint
	ContentType    string
	TimeoutSeconds time.Duration
}

type SendAction struct {
	Name         string
	PayloadType  internal.PayloadType
	StatusCode   int
	Headers      map[string]string
	Payload      string
	EndpointName string
}

type ReceiveAction struct {
	Name         string
	Url          string
	Path         []string
	Method       string
	QueryParams  map[string][]string
	Headers      map[string]string
	Payload      string
	PayloadType  internal.PayloadType
	EndpointName string
}

func (action *ReceiveAction) ToString() string {
	return fmt.Sprintf(
		"["+
			"Method: %s, "+
			"BaseUrl: %s, "+
			"Path: '%s', "+
			"Headers: %s, "+
			"QueryParams: %s, "+
			"Payload: %s"+
			"]",
		action.Method, action.Url, action.Path,
		action.Headers, action.QueryParams, action.Payload)
}
