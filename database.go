package sofa

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strconv"
)

// DatabaseMetadata contains information about the database in CouchDB.
type DatabaseMetadata struct {
	CompactRunning     bool         `json:"compact_running"`
	Name               string       `json:"db_name"`
	DocCount           int          `json:"doc_count"`
	DelCount           int          `json:"doc_del_count"`
	InstanceStartTime  string       `json:"instance_start_time"`
	DataSize           int          `json:"data_size"`
	DiskSize           int          `json:"disk_size"`
	DiskFormatVersion  int          `json:"disk_format_version"`
	PurgeSeq           AlwaysString `json:"purge_seq"`
	UpdateSeq          AlwaysString `json:"update_seq"`
	CommittedUpdateSeq AlwaysString `json:"committed_update_seq"`
}

// Database represents a CouchDB database & provides methods to access documents in the database.
type Database struct {
	name     string
	metadata *DatabaseMetadata
	con      *Connection
}

// Get retrieves a single document from the database and unmarshals it into the
// provided interface.
func (d *Database) Get(document Document, id, rev string) (string, error) {
	path := d.DocumentPath(id)

	var opts = NewURLOptions()
	if rev != "" {
		if err := opts.Add("rev", rev); err != nil {
			return "", err
		}
	}

	resp, err := d.con.unmarshalRequest("GET", path, opts, nil, document)
	if err != nil {
		return "", err
	}

	return responseEtag(resp)
}

// Put marshals the provided document into JSON and the sends it to the CouchDB server
// with a PUT request. This allows modification of the document on the server.
func (d *Database) Put(document Document) (string, error) {
	docMeta := document.Metadata()
	path := d.DocumentPath(docMeta.ID)

	var opts = NewURLOptions()
	if docMeta.Rev != "" {
		if err := opts.Add("rev", docMeta.Rev); err != nil {
			return "", err
		}
	}

	b, err := json.Marshal(document)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(b)

	res := ServerResponse{}
	resp, err := d.con.unmarshalRequest("PUT", path, opts, buf, &res)
	if err != nil {
		return "", err
	}

	return responseEtag(resp)
}

// Delete removed a document from the Database.
func (d *Database) Delete(document Document) (string, error) {
	docMeta := document.Metadata()
	path := d.DocumentPath(docMeta.ID)

	var opts = NewURLOptions()
	if docMeta.Rev != "" {
		if err := opts.Add("rev", docMeta.Rev); err != nil {
			return "", err
		}
	}

	res := ServerResponse{}
	resp, err := d.con.unmarshalRequest("DELETE", path, opts, nil, &res)
	if err != nil {
		return "", err
	}

	return responseEtag(resp)
}

// Name returns the name of the Database.
func (d *Database) Name() string {
	return d.name
}

// Path returns the path for this database. The value (as with the other *Path functions)
// is in the correct for to be passed to Connection.URL()
func (d *Database) Path() string {
	return "/" + d.name
}

// ViewPath returns the path to a view in this database
func (d *Database) ViewPath(view string) string {
	return urlConcat(d.Path(), view)
}

// DocumentPath returns the path to a document in this database
func (d *Database) DocumentPath(id string) string {
	return urlConcat(d.Path(), id)
}

// ViewCleanup cleans old data from the database. From the CouchDB API documentation:
// "Old view output remains on disk until you explicitly run cleanup."
func (d *Database) ViewCleanup() error {
	path := urlConcat(d.Path(), "_view_cleanup")
	_, err := d.con.Post(path, NewURLOptions(), nil)
	return err
}

// CompactView compacts the stored data for a view, meaning that the view uses less space on
// the disk.
func (d *Database) CompactView(name string) error {
	path := urlConcat(urlConcat(d.Path(), "_compact"), name)
	_, err := d.con.Post(path, NewURLOptions(), nil)
	return err
}

// Metadata downloads the Metadata for the Database and saves it to the Database object. If
// there is already metadata stored then that will be returned without contacting the server.
func (d *Database) Metadata() (DatabaseMetadata, error) {
	if d.metadata != nil {
		return *d.metadata, nil
	}

	var metadata DatabaseMetadata
	if _, err := d.con.unmarshalRequest("GET", d.Path(), NewURLOptions(), nil, &metadata); err != nil {
		return DatabaseMetadata{}, err
	}

	d.metadata = &metadata
	return metadata, nil
}

// AllDocuments gets all documents from a database. All document content is included for each row.
func (d *Database) AllDocuments() (DocumentList, error) {
	resp, err := d.con.Get(d.ViewPath("_all_docs"), URLOptions{url.Values{"include_docs": []string{"true"}}})
	if err != nil {
		return DocumentList{}, err
	}

	var docs DocumentList
	err = json.NewDecoder(resp.Body).Decode(&docs)
	if err != nil {
		return DocumentList{}, err
	}
	return docs, nil
}

// ListDocuments gets all rows from the Database but does not include the content of the documents.
func (d *Database) ListDocuments() (DocumentList, error) {
	resp, err := d.con.Get(d.ViewPath("_all_docs"), NewURLOptions())
	if err != nil {
		return DocumentList{}, err
	}

	var docs DocumentList
	err = json.NewDecoder(resp.Body).Decode(&docs)
	if err != nil {
		return DocumentList{}, err
	}
	return docs, nil
}

// Documents gets a set of Documents from a database. All of the IDs requested will be downloaded &
// all document content is included for each row.
func (d *Database) Documents(ids ...string) (DocumentList, error) {
	body := map[string]interface{}{"keys": ids}
	bodyBytes, err := json.Marshal(&body)
	if err != nil {
		return DocumentList{}, err
	}

	bodyBuf := bytes.NewBuffer(bodyBytes)

	resp, err := d.con.Post(d.ViewPath("_all_docs"), URLOptions{url.Values{"include_docs": []string{"true"}}}, bodyBuf)
	if err != nil {
		return DocumentList{}, err
	}

	var docs DocumentList
	err = json.NewDecoder(resp.Body).Decode(&docs)
	if err != nil {
		return DocumentList{}, err
	}
	return docs, nil
}

// PollingChangesFeed gets a changes feed which will poll the server for changes to documents.
func (d *Database) PollingChangesFeed(long bool) PollingChangesFeed {
	var t = "normal"
	if long {
		t = "longpoll"
	}

	return PollingChangesFeed{
		db:       d,
		feedType: t,
	}
}

// ContinuousChangesFeed gets a changes feed with a continuous connection to the database. New
// changes are then pushed over the existing connection as they arrive.
func (d *Database) ContinuousChangesFeed(params ChangesFeedParams) ContinuousChangesFeed {
	return ContinuousChangesFeed{
		db:     d,
		params: params,
	}
}
