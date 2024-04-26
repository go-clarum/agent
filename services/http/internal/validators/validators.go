package validators

import (
	"errors"
	"fmt"
	api "github.com/go-clarum/agent/api/http"
	"github.com/go-clarum/agent/arrays"
	"github.com/go-clarum/agent/logging"
	clarumstrings "github.com/go-clarum/agent/validators/strings"
	"github.com/go-clarum/clarum-json/comparator"
	"github.com/go-clarum/clarum-json/recorder"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

func ValidatePath(expectedPath string, actualUrl *url.URL, logger *logging.Logger) error {
	cleanedExpected := cleanPath(expectedPath)
	cleanedActual := cleanPath(actualUrl.Path)

	if cleanedExpected != cleanedActual {
		return handleError(logger, "validation error - path mismatch - expected [%s] but received [%s]",
			cleanedExpected, cleanedActual)
	} else {
		logger.Info("path validation successful")
	}

	return nil
}

func ValidateHttpMethod(expectedMethod string, actualMethod string, logger *logging.Logger) error {
	if expectedMethod != actualMethod {
		return handleError(logger, "validation error - method mismatch - expected [%s] but received [%s]",
			expectedMethod, actualMethod)
	} else {
		logger.Info("method validation successful")
	}

	return nil
}

func ValidateHttpHeaders(expectedHeaders map[string]string, actualHeaders http.Header, logger *logging.Logger) error {
	if err := validateHeaders(expectedHeaders, actualHeaders); err != nil {
		return handleError(logger, "%s", err)
	} else {
		logger.Info("header validation successful")
	}

	return nil
}

// According to the official specification, HTTP headers must be compared in a case-insensitive way
func validateHeaders(expectedHeaders map[string]string, actualHeaders http.Header) error {
	lowerCaseReceivedHeaders := make(map[string][]string)
	for header, values := range actualHeaders {
		lowerCaseReceivedHeaders[strings.ToLower(header)] = values
	}

	for header, expectedValue := range expectedHeaders {
		lowerCaseExpectedHeader := strings.ToLower(header)
		if receivedValues, exists := lowerCaseReceivedHeaders[lowerCaseExpectedHeader]; exists {
			if !arrays.Contains(receivedValues, expectedValue) {
				return errors.New(fmt.Sprintf("validation error - header <%s> mismatch - expected [%s] but received [%s]",
					lowerCaseExpectedHeader, expectedValue, receivedValues))
			}
		} else {
			return errors.New(fmt.Sprintf("validation error - header <%s> missing", lowerCaseExpectedHeader))
		}
	}

	return nil
}

func ValidateHttpQueryParams(expectedQueryParams *map[string][]string, actualUrl *url.URL, logger *logging.Logger) error {
	if err := validateQueryParams(expectedQueryParams, actualUrl.Query()); err != nil {
		return handleError(logger, "%s", err)
	} else {
		logger.Info("query params validation successful")
	}

	return nil
}

// validate query parameters based on these rules
//
//	-> validate that the param exists
//	-> that the values match
func validateQueryParams(expectedQueryParams *map[string][]string, params url.Values) error {
	for param, expectedValues := range *expectedQueryParams {
		if receivedValues, exists := params[param]; exists {
			for _, expectedValue := range expectedValues {
				if !arrays.Contains(receivedValues, expectedValue) {
					return errors.New(fmt.Sprintf("validation error - query param <%s> values mismatch - expected [%v] but received [%s]",
						param, expectedValues, receivedValues))
				}
			}
		} else {
			return errors.New(fmt.Sprintf("validation error - query param <%s> missing", param))
		}
	}

	return nil
}

func ValidateHttpStatusCode(expectedStatusCode int, actualStatusCode int, logger *logging.Logger) error {
	if actualStatusCode != expectedStatusCode {
		return handleError(logger, "validation error - status mismatch - expected [%d] but received [%d]",
			expectedStatusCode, actualStatusCode)
	} else {
		logger.Info("status validation successful")
	}

	return nil
}

func ValidateHttpPayload(expectedPayload *string, actualPayload io.ReadCloser,
	payloadType api.PayloadType, logger *logging.Logger) error {
	defer closeBody(logger, actualPayload)

	if clarumstrings.IsBlank(*expectedPayload) {
		logger.Info("message payload is empty - no body validation will be done")
		return nil
	}

	bodyBytes, err := io.ReadAll(actualPayload)
	if err != nil {
		return handleError(logger, "could not read response body - %s", err)
	}

	if err := validatePayload(expectedPayload, bodyBytes, payloadType, logger); err != nil {
		return handleError(logger, "%s", err)
	} else {
		logger.Info("payload validation successful")
	}

	return nil
}

func closeBody(logger *logging.Logger, body io.ReadCloser) {
	if err := body.Close(); err != nil {
		logger.Errorf("unable to close body - %s", err)
	}
}

func validatePayload(expected *string, actual []byte, payloadType api.PayloadType, logger *logging.Logger) error {

	if len(actual) == 0 {
		return errors.New(fmt.Sprintf("validation error - payload missing - expected [%s] but received no payload",
			expected))
	} else if payloadType == api.PayloadType_Plaintext {
		receivedPayload := string(actual)

		if *expected != receivedPayload {
			return errors.New(fmt.Sprintf("validation error - payload mismatch - expected [%s] but received [%s]",
				expected, receivedPayload))
		}
	} else if payloadType == api.PayloadType_Json {
		jsonComparator := comparator.NewComparator().
			Recorder(recorder.NewDefaultRecorder()).
			Build()

		reporterLog, errs := jsonComparator.Compare([]byte(*expected), actual)

		if errs != nil {
			logger.Infof("json validation log: %s", reporterLog)
			return errors.New(fmt.Sprintf("json validation errors: [%s]", errs))
		}
		logger.Debugf("json payload validation log: %s", reporterLog)
	}

	return nil
}

func handleError(logger *logging.Logger, format string, a ...any) error {
	errorMessage := fmt.Sprintf(format, a...)
	logger.Errorf(errorMessage)
	return errors.New(errorMessage)
}

// path.Clean() does not remove leading "/", so we do that ourselves
func cleanPath(pathToClean string) string {
	return strings.TrimPrefix(path.Clean(pathToClean), "/")
}
