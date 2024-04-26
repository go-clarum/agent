package request

import (
	"fmt"
	"github.com/go-clarum/agent/services/http/internal/constants"
	"github.com/go-clarum/agent/services/http/internal/utils"
	"maps"
	"net/http"
	"reflect"
)

type Request struct {
	Url            string
	Path           string
	Method         string
	QueryParams    map[string][]string
	Headers        map[string]string
	RequestPayload string
}

func Get(pathElements ...string) *Request {
	return &Request{
		Method: http.MethodGet,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Head(pathElements ...string) *Request {
	return &Request{
		Method: http.MethodHead,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Post(pathElements ...string) *Request {
	return &Request{
		Method: http.MethodPost,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Put(pathElements ...string) *Request {
	return &Request{
		Method: http.MethodPut,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Delete(pathElements ...string) *Request {
	return &Request{
		Method: http.MethodDelete,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Options(pathElements ...string) *Request {
	return &Request{
		Method: http.MethodOptions,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Trace(pathElements ...string) *Request {
	return &Request{
		Method: http.MethodTrace,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Patch(pathElements ...string) *Request {
	return &Request{
		Method: http.MethodPatch,
		Path:   utils.BuildPath("", pathElements...),
	}
}

// BaseUrl - While this should normally be configured only on the HTTP client,
// this is also allowed on the message so that a client can send a request to different targets.
// When used on a message passed to an HTTP server, it will do nothing.
func (request *Request) BaseUrl(baseUrl string) *Request {
	request.Url = baseUrl
	return request
}

func (request *Request) Header(key string, value string) *Request {
	request.header(key, value)
	return request
}

func (request *Request) ContentType(value string) *Request {
	request.contentType(value)
	return request
}

func (request *Request) Authorization(value string) *Request {
	request.authorization(value)
	return request
}

func (request *Request) QueryParam(key string, values ...string) *Request {
	if request.QueryParams == nil {
		request.QueryParams = make(map[string][]string)
	}

	if _, exists := request.QueryParams[key]; exists {
		for _, value := range values {
			request.QueryParams[key] = append(request.QueryParams[key], value)
		}
	} else {
		request.QueryParams[key] = values
	}

	return request
}

func (request *Request) Payload(payload string) *Request {
	request.RequestPayload = payload
	return request
}

func (request *Request) Clone() *Request {
	return &Request{
		Url:            request.Url,
		Path:           request.Path,
		Method:         request.Method,
		QueryParams:    maps.Clone(request.QueryParams),
		Headers:        maps.Clone(request.Headers),
		RequestPayload: request.RequestPayload,
	}
}

func (request *Request) Equals(other *Request) bool {
	if request.Method != other.Method {
		return false
	} else if request.Url != other.Url {
		return false
	} else if request.Path != other.Path {
		return false
	} else if !maps.Equal(request.Headers, other.Headers) {
		return false
	} else if !reflect.DeepEqual(request.QueryParams, other.QueryParams) {
		return false
	} else if request.RequestPayload != other.RequestPayload {
		return false
	}
	return true
}

func (request *Request) ToString() string {
	return fmt.Sprintf(
		"["+
			"Method: %s, "+
			"BaseUrl: %s, "+
			"Path: '%s', "+
			"Headers: %s, "+
			"QueryParams: %s, "+
			"RequestPayload: %s"+
			"]",
		request.Method, request.Url, request.Path,
		request.Headers, request.QueryParams, request.RequestPayload)
}

func (request *Request) header(key string, value string) *Request {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}

	request.Headers[key] = value
	return request
}

func (request *Request) contentType(value string) *Request {
	return request.header(constants.ContentTypeHeaderName, value)
}

func (request *Request) authorization(value string) *Request {
	return request.header(constants.AuthorizationHeaderName, value)
}

func (request *Request) eTag(value string) *Request {
	return request.header(constants.ETagHeaderName, value)
}

func (request *Request) payload(payload string) *Request {
	request.RequestPayload = payload
	return request
}

func (request *Request) clone() Request {
	return Request{
		Headers:        maps.Clone(request.Headers),
		RequestPayload: request.RequestPayload,
	}
}
