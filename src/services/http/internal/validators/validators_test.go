package validators

import (
	_ "github.com/go-clarum/agent/config"
	"github.com/go-clarum/agent/logging"
	"github.com/go-clarum/agent/services/http/internal/constants"
	"net/http"
	"testing"
)

var logger = logging.NewLogger("validators test: ")

func TestValidatePathOK(t *testing.T) {
	req := createRealRequest()

	if err := ValidatePath("myPath/some/api", req.URL, logger); err != nil {
		t.Errorf("No header validation error expected, but got %s", err)
	}
}

func TestValidatePathError(t *testing.T) {
	req := createRealRequest()

	err := ValidatePath("blup/", req.URL, logger)

	if err == nil {
		t.Errorf("Path validation error expected, but got none")
	}

	if err.Error() != "validation error - path mismatch - expected [blup] but received [myPath/some/api]" {
		t.Errorf("Path validation error message is unexpected")
	}
}

func TestValidateMethodOK(t *testing.T) {
	req := createRealRequest()

	if err := ValidateHttpMethod(http.MethodPost, req.Method, logger); err != nil {
		t.Errorf("No header validation error expected, but got %s", err)
	}
}

func TestValidateMethodError(t *testing.T) {
	req := createRealRequest()

	err := ValidateHttpMethod(http.MethodOptions, req.Method, logger)

	if err == nil {
		t.Errorf("Method validation error expected, but got none")
	}

	if err.Error() != "validation error - method mismatch - expected [OPTIONS] but received [POST]" {
		t.Errorf("Path validation error message is unexpected")
	}
}

func TestValidateHeadersOK(t *testing.T) {
	expectedHeaders := make(map[string]string)
	expectedHeaders["Connection"] = "keep-alive"
	expectedHeaders[constants.ContentTypeHeaderName] = "application/json"
	expectedHeaders[constants.AuthorizationHeaderName] = "Bearer 0b79bab50daca910b000d4f1a2b675d604257e42"

	req := createRealRequest()

	if err := ValidateHttpHeaders(&expectedHeaders, req.Header, logger); err != nil {
		t.Errorf("No header validation error expected, but got %s", err)
	}
}

func TestValidateHeaderValueError(t *testing.T) {
	expectedHeaders := make(map[string]string)
	expectedHeaders["Connection"] = "keep-alive"
	expectedHeaders[constants.ContentTypeHeaderName] = "application/json"
	expectedHeaders[constants.AuthorizationHeaderName] = "something else"

	req := createRealRequest()

	err := ValidateHttpHeaders(&expectedHeaders, req.Header, logger)

	if err == nil {
		t.Errorf("Header validation error expected, but got none")
	}

	if err.Error() != "validation error - header <authorization> mismatch - expected [something else] but received [[Bearer 0b79bab50daca910b000d4f1a2b675d604257e42]]" {
		t.Errorf("Header validation error message is unexpected")
	}
}

func TestValidateMissingHeaderError(t *testing.T) {
	expectedHeaders := make(map[string]string)
	expectedHeaders["Connection"] = "keep-alive"
	expectedHeaders[constants.ContentTypeHeaderName] = "application/json"
	expectedHeaders[constants.AuthorizationHeaderName] = "something else"
	expectedHeaders["traceid"] = "124245132"

	req := createRealRequest()

	err := ValidateHttpHeaders(&expectedHeaders, req.Header, logger)

	if err == nil {
		t.Errorf("Header validation error expected, but got none")
	}

	if err.Error() != "validation error - header <traceid> missing" {
		t.Errorf("Header validation error message is unexpected")
	}
}

func TestValidateQueryParamsOK(t *testing.T) {
	expectedParams := make(map[string][]string)
	expectedParams["param1"] = []string{"value1"}
	expectedParams["param2"] = []string{"value2"}

	req := createRealRequest()
	qParams := req.URL.Query()
	qParams.Set("param1", "value1")
	qParams.Set("param2", "value2")
	req.URL.RawQuery = qParams.Encode()

	if err := ValidateHttpQueryParams(&expectedParams, req.URL, logger); err != nil {
		t.Errorf("No query param validation error expected, but got %s", err)
	}
}

func TestValidateQueryParamsParamMismatch(t *testing.T) {
	expectedParams := make(map[string][]string)
	expectedParams["param1"] = []string{"value1"}
	expectedParams["param2"] = []string{"value2"}

	req := createRealRequest()
	qParams := req.URL.Query()
	qParams.Set("param1", "value1")
	qParams.Set("param3", "value2")
	req.URL.RawQuery = qParams.Encode()

	err := ValidateHttpQueryParams(&expectedParams, req.URL, logger)
	if err == nil {
		t.Errorf("Query param validation error expected, but got none")
	}

	if err.Error() != "validation error - query param <param2> missing" {
		t.Errorf("Query param validation error message is unexpected")
	}
}

func TestValidateQueryParamsValueMismatch(t *testing.T) {
	expectedParams := make(map[string][]string)
	expectedParams["param1"] = []string{"value1"}
	expectedParams["param2"] = []string{"value2"}

	req := createRealRequest()
	qParams := req.URL.Query()
	qParams.Set("param1", "value1")
	qParams.Set("param2", "value22")
	req.URL.RawQuery = qParams.Encode()

	err := ValidateHttpQueryParams(&expectedParams, req.URL, logger)
	if err == nil {
		t.Errorf("Query param validation error expected, but got none")
	}

	if err.Error() != "validation error - query param <param2> values mismatch - expected [[value2]] but received [[value22]]" {
		t.Errorf("Query param validation error message is unexpected")
	}
}

func TestValidateQueryParamsMultiValueOK(t *testing.T) {
	expectedParams := make(map[string][]string)
	expectedParams["param1"] = []string{"value1", "value3"}

	req := createRealRequest()
	qParams := req.URL.Query()
	qParams.Set("param1", "value1")
	qParams.Add("param1", "value2")
	qParams.Add("param1", "value3")
	req.URL.RawQuery = qParams.Encode()

	if err := ValidateHttpQueryParams(&expectedParams, req.URL, logger); err != nil {
		t.Errorf("No query param validation error expected, but got %s", err)
	}
}

func TestValidateQueryParamsMultiValueMismatch(t *testing.T) {
	expectedParams := make(map[string][]string)
	expectedParams["param1"] = []string{"value1", "value2", "value4"}

	req := createRealRequest()
	qParams := req.URL.Query()
	qParams.Set("param1", "value1")
	qParams.Add("param1", "value2")
	qParams.Add("param1", "value3")
	req.URL.RawQuery = qParams.Encode()

	err := ValidateHttpQueryParams(&expectedParams, req.URL, logger)
	if err == nil {
		t.Errorf("Query param validation error expected, but got none")
	}

	if err.Error() != "validation error - query param <param1> values mismatch - expected [[value1 value2 value4]] but received [[value1 value2 value3]]" {
		t.Errorf("Query param validation error message is unexpected")
	}
}

func createRealRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "myPath/some/api", nil)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set(constants.ContentTypeHeaderName, "application/json")
	req.Header.Set(constants.AuthorizationHeaderName, "Bearer 0b79bab50daca910b000d4f1a2b675d604257e42")

	return req
}
