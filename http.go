package sofa

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Options is a slightly modified version of url.Options which provides query options
// with JSON encoding of values.
type Options interface {
	Encode() string
}

// NewURLOptions creates a new URLOptions struct to hold HTTP query options passed to the
// CouchDB server.
func NewURLOptions() URLOptions {
	return URLOptions{
		url.Values{},
	}
}

// URLOptions is a slightly modified version of url.Values, the only difference being
// that values are automatically encoded as JSON when they are added to the URLOptions
// map as this is the format CouchDB will expect them in. This means that the
// URLOptions.Add and URLOptions.Set methods can return an error in the case that the passed
// value cannot be encoded as JSON.
type URLOptions struct {
	url.Values
}

// Set overwrites the currently-stored value for name with the JSON-encoded version of val.
func (opts *URLOptions) Set(name string, val interface{}) error {
	encoded, err := encodeValue(val)
	if err != nil {
		return err
	}

	opts.Values.Set(name, encoded)
	return nil
}

// Add adds the JSON-encoded version of val to the list of values stored for name.
func (opts *URLOptions) Add(name string, val interface{}) error {
	encoded, err := encodeValue(val)
	if err != nil {
		return err
	}

	opts.Values.Add(name, encoded)
	return nil
}

// urlConcat concatenates two strings together with a single '/' character, removing
// an instance of the character from one of the strings if required.
func urlConcat(x, y string) string {
	xs := strings.HasSuffix(x, "/")
	ys := strings.HasPrefix(y, "/")

	switch {
	case xs && ys:
		return x + y[1:]
	case !xs && !ys:
		return x + "/" + y
	}

	return x + y
}

// responseEtag returns the Etag from a response, with quotes removed.
func responseEtag(resp *http.Response) (string, error) {
	etag := resp.Header.Get("Etag")
	if etag == "" {
		return "", fmt.Errorf("couchdb: missing Etag header in response")
	}

	return etag[1 : len(etag)-1], nil
}

// newRequest returns a new http.Request with default headers set for accessing
// couchdb.
func newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// encodeValue encodes a value as JSON unless it is already a string
func encodeValue(val interface{}) (string, error) {
	switch t := val.(type) {
	case BooleanParameter:
		return t.String(), nil
	}

	// rv := reflect.ValueOf(val)
	// switch rv.Kind() {
	// case reflect.String:
	// 	return rv.String(), nil
	// }

	bytes, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
