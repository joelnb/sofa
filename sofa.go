package sofa

import (
	"net/http"
)

// ServerDetails represents the details returned from a CouchDB server when
// requesting the root page.
type ServerDetails struct {
	CouchDB string                 `json:"couchdb"`
	UUID    string                 `json:"uuid"`
	Version string                 `json:"version"`
	Vendor  map[string]interface{} `json:"vendor"`
}

type ServerResponse struct {
	RawResponse *http.Response
	ResultBody  *struct {
		OK  bool   `json:"ok"`
		ID  string `json:"id"`
		Rev string `json:"rev"`
	}
}

func (resp *ServerResponse) HasBody() bool {
	return resp.ResultBody != nil
}
