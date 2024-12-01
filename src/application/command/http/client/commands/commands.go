package commands

import (
	"fmt"
	"github.com/go-clarum/agent/application/command/http/common/model"
	"strconv"
	"time"
)

type InitEndpointCommand struct {
	Name           string
	BaseUrl        string
	ContentType    string
	TimeoutSeconds time.Duration
}

type SendCommand struct {
	Name         string
	Url          string
	Path         []string
	Method       string
	QueryParams  map[string][]string
	Headers      map[string]string
	Payload      string
	EndpointName string
}

type ReceiveCommand struct {
	Name         string
	PayloadType  model.PayloadType
	StatusCode   int
	Headers      map[string]string
	Payload      string
	EndpointName string
}

func (action *SendCommand) ToString() string {
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

func (action *ReceiveCommand) ToString() string {
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
