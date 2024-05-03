package client

import (
	"fmt"
	"github.com/go-clarum/agent/services/http/internal"
	"strconv"
	"time"
)

type InitRequest struct {
	Name           string
	BaseUrl        string
	ContentType    string
	TimeoutSeconds time.Duration
}

type SendAction struct {
	Name         string
	Url          string
	Path         string
	Method       string
	QueryParams  map[string][]string
	Headers      map[string]string
	Payload      string
	EndpointName string
}

type ReceiveAction struct {
	Name         string
	PayloadType  internal.PayloadType
	StatusCode   int
	Headers      map[string]string
	Payload      string
	EndpointName string
}

func (action *SendAction) ToString() string {
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

func (action *ReceiveAction) ToString() string {
	statusCodeText := "none"
	if action.StatusCode > 0 {
		statusCodeText = strconv.Itoa(action.StatusCode)
	}

	return fmt.Sprintf(
		"["+
			"StatusCode: %s, "+
			"Headers: %s, "+
			"Payload: %s"+
			"]",
		statusCodeText, action.Headers, action.Payload)
}
