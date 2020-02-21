package sofa

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Attachment represents files or other data which is attached to documents in the
// Database.
type Attachment struct {
	ContentType   string  `json:"content_type,omitempty"`
	Data          string  `json:"data,omitempty"`
	Digest        string  `json:"digest,omitempty"`
	EncodedLength float64 `json:"encoded_length,omitempty"`
	Encoding      string  `json:"encoding,omitempty"`
	Length        int64   `json:"length,omitempty"`
	RevPos        float64 `json:"revpos,omitempty"`
	Stub          bool    `json:"stub,omitempty"`
	Follows       bool    `json:"follows,omitempty"`
}

type attachmentPutResponse struct {
	ID  string `json:"id"`
	Rev string `json:"rev"`
	OK  bool   `json:"ok"`
}

// GetAttachment gets the current attachment and returns
func (db *Database) GetAttachment(docid, name, rev string) ([]byte, error) {
	path := urlConcat(db.DocumentPath(docid), name)

	opts := NewURLOptions()
	if rev != "" {
		if err := opts.Set("rev", rev); err != nil {
			return nil, err
		}
	}

	resp, err := db.con.Get(path, opts)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

// PutAttachment replaces the content of the attachment with new content read from an
// io.Reader or creates it if it does not exist. If the provided rev is not the most
// recent then an error will be returned from CouchDB.
func (db *Database) PutAttachment(docid, name string, doc io.Reader, rev string) (string, error) {
	path := urlConcat(db.DocumentPath(docid), name)

	opts := NewURLOptions()
	if rev != "" {
		if err := opts.Set("rev", rev); err != nil {
			return "", err
		}
	}

	resp, err := db.con.Put(path, opts, doc)
	if err != nil {
		return "", err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ar := attachmentPutResponse{}
	if err := json.Unmarshal(respBytes, &ar); err != nil {
		return "", err
	}

	return ar.Rev, nil
}

// DeleteAttachment removes an attachment from a document in CouchDB.
func (db *Database) DeleteAttachment(docid, name, rev string) (string, error) {
	path := urlConcat(db.DocumentPath(docid), name)

	opts := NewURLOptions()
	if err := opts.Set("rev", rev); err != nil {
		return "", err
	}

	resp, err := db.con.Delete(path, opts)
	if err != nil {
		return "", err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ar := attachmentPutResponse{}
	if err := json.Unmarshal(respBytes, &ar); err != nil {
		return "", err
	}

	return ar.Rev, nil
}
