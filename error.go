package sofa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// ResponseError represents an error response from CouchDB. The method & URL of the request
// are included as well as the HTTP status code. If the request was not a HEAD request then
// the values from the JSON error returned by CouchDB will be included.
type ResponseError struct {
	Method     string
	StatusCode int
	URL        string

	Err    string
	Reason string
}

// ErrorStatus checks returns true if the provided error is a
// ResponseError with a status code which matches the one provided.
func ErrorStatus(err error, statusCode int) bool {
	return isResponseError(err) && err.(ResponseError).StatusCode == statusCode
}

// Error provodes a detailed representation of the response error
// including as much information as there is available
func (e ResponseError) Error() string {
	if e.Err == "" {
		return fmt.Sprintf("%v %v: %v", e.Method, e.URL, e.StatusCode)
	}

	return fmt.Sprintf("%v %v: (%v) %v: %v", e.Method, e.URL, e.StatusCode, e.Err, e.Reason)
}

// isResponseError checks if a given error is a ResponseError.
func isResponseError(err error) bool {
	val := reflect.TypeOf(err)
	errType := reflect.TypeOf((*ResponseError)(nil)).Elem()
	return val == errType
}

// httpResponseError parses a http.Response from CouchDB and returns a
// ResponseError with the details.
func httpResponseError(resp *http.Response) error {
	var couchErr struct{ Error, Reason string }

	if resp.Request.Method != "HEAD" {
		if err := unmarshalResponse(resp, &couchErr); err != nil {
			return fmt.Errorf("unable to parse couchdb error: %v", err)
		}
	}

	return ResponseError{
		Method:     resp.Request.Method,
		URL:        resp.Request.URL.String(),
		StatusCode: resp.StatusCode,
		Err:        couchErr.Error,
		Reason:     couchErr.Reason,
	}
}

// unmasrshallResponse attempts to unmarshal the json from the response body
// into res.
func unmarshalResponse(resp *http.Response, res interface{}) error {
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		resp.Body.Close()
		return err
	}

	return resp.Body.Close()
}
