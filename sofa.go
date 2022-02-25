package sofa

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// AlwaysString allows fields which can be either an int or string (depending on the
// version) to always be unmarshaled as a string.
type AlwaysString string

// UnmarshalJSON implements the json.Unmarshaler interface, returning directly as a
// string if one is found or converting an int if that is found. Other types are not
// supported and the error from json.Unmarshal will be returned.
func (fs *AlwaysString) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		return json.Unmarshal(b, (*string)(fs))
	}

	var i int
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	*fs = AlwaysString(strconv.Itoa(i))

	return nil
}

// ServerDetails1 represents the details returned from a CouchDB version 1 server when
// requesting the root page.
type ServerDetails1 struct {
	CouchDB string                 `json:"couchdb"`
	UUID    string                 `json:"uuid"`
	Version string                 `json:"version"`
	Vendor  map[string]interface{} `json:"vendor"`
}

// ServerDetails2 represents the details returned from a CouchDB version 2 server when
// requesting the root page.
type ServerDetails2 struct {
	CouchDB  string                 `json:"couchdb"`
	Features []string               `json:"features"`
	Version  string                 `json:"version"`
	Vendor   map[string]interface{} `json:"vendor"`
}

// ServerResponse is a parsed CouchDB response which also contains a RawResponse
// field containing a pointer to the unaltered http.Response.
type ServerResponse struct {
	RawResponse *http.Response
	ResultBody  *struct {
		OK  bool   `json:"ok"`
		ID  string `json:"id"`
		Rev string `json:"rev"`
	}
}

// HasBody returns true if the response from the server has a body (will return
// false for HEAD requests).
func (resp *ServerResponse) HasBody() bool {
	return resp.ResultBody != nil
}
