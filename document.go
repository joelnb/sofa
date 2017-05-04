package sofa

import "encoding/json"

// Document is the interface which represents a CouchDB document. Contains methods to retrieve the
// metadata and the database containing this document.
type Document interface {
	Metadata() DocumentMetadata
}

// DocumentMetadata is the minimum amount of content which all documents have. It is the
// identity and revision of the document.
type DocumentMetadata struct {
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
}

func (md DocumentMetadata) Metadata() DocumentMetadata {
	return md
}

// GenericDocument implements the Document API and can be used to represent any type
// of document. For all but simple cases it is better to implement this yourself ona struct
// for easier access to the unmarshalled data.
type GenericDocument struct {
	document map[string]interface{}
}

// Metadata gets the DocumentMetadata for this Document
func (gen *GenericDocument) Metadata() DocumentMetadata {
	return DocumentMetadata{
		ID:  gen.document["_id"].(string),
		Rev: gen.document["_rev"].(string),
	}
}

// MarshalJSON provides an implementation of json.Marshaler by marshaling the stored map document
// into JSON data.
func (gen *GenericDocument) MarshalJSON() ([]byte, error) {
	return json.Marshal(gen.document)
}

// UnmarshalJSON provides an implementation of json.Unmarshaler by unmarshaling the provides
// data into the stored map.
func (gen *GenericDocument) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &gen.document)
}
