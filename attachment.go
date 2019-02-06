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

func (db *Database) GetAttachment(docid, name string) ([]byte, error) {
	path := urlConcat(db.DocumentPath(docid), name)

	resp, err := db.con.Get(path, NewURLOptions())
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func (db *Database) PutAttachment(docid, name string, doc io.Reader) (string, error) {
	path := urlConcat(db.DocumentPath(docid), name)

	resp, err := db.con.Put(path, NewURLOptions(), doc)
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
