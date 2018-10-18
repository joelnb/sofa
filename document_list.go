package sofa

import (
	"encoding/json"
)

// DocumentList is the CouchDB reprentation of a list of Documents where each document
// is contained within a Row along with some metadata. Depending on the view & the request
// paramaters used the documents may or may not be included in the row.
type DocumentList struct {
	TotalRows float64 `json:"total_rows"`
	Offset    float64 `json:"offset"`
	Rows      []Row   `json:"rows"`
}

// Size returns the number of documents currently stored in this DocumentList.
func (dl *DocumentList) Size() int {
	// TODO: Possibly should just return `int(dl.TotalRows)`
	return len(dl.Rows)
}

// RawDocuments returns the Document from each Row of this DocumentList as a json.RawMessage
// which can then be Unmarshalled into the correct type.
func (dl *DocumentList) RawDocuments() []json.RawMessage {
	var docs []json.RawMessage

	for _, row := range dl.Rows {
		if row.HasDocument() {
			docs = append(docs, row.Document)
		} else {
			docs = append(docs, nil)
		}
	}

	return docs
}

// MapDocuments returns a list containing the document from each row in this DocumentList
// Unmarshalled into a map[string]interface{}.
func (dl *DocumentList) MapDocuments() ([]map[string]interface{}, error) {
	var docs []map[string]interface{}

	for _, row := range dl.Rows {
		if row.HasDocument() {
			var rowDoc map[string]interface{}

			err := json.Unmarshal(row.Document, &rowDoc)
			if err != nil {
				return nil, err
			}

			docs = append(docs, rowDoc)
		}
	}

	return docs, nil
}

// UnmarshalDocuments is a convenience function which unmarshalls just the list of Documents from
// this. Due to a Marshal/Unmarhal cycle this method may be slow. Sometimes slow is fine though.
func (dl *DocumentList) UnmarshalDocuments(docs interface{}) error {
	rawDocs := dl.RawDocuments()

	jsonAllDocs, err := json.Marshal(rawDocs)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonAllDocs, docs)
}

// Row represents a row in a CouchDB DocumentList. These are mainly returned by views but
// may also be sent to the server as part of the bulk documents API.
type Row struct {
	ID       string          `json:"id,omitempty"`
	Key      interface{}     `json:"key,omitempty"`
	Value    interface{}     `json:"value,omitempty"`
	Document json.RawMessage `json:"doc,omitempty"`
}

// HasDocument checks the current Row for the presence of the Document attached to it.
func (r Row) HasDocument() bool {
	if r.Document == nil {
		return false
	}
	return true
}
