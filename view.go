package sofa

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// conflicts (boolean) – Includes conflicts information in response. Ignored if include_docs isn’t true. Default is false
// descending (boolean) – Return the documents in descending by key order. Default is false
// endkey (json) – Stop returning records when the specified key is reached. Optional
// end_key (json) – Alias for endkey param
// endkey_docid (string) – Stop returning records when the specified document ID is reached. Requires endkey to be specified for this to have any effect. Optional
// end_key_doc_id (string) – Alias for endkey_docid param
// group (boolean) – Group the results using the reduce function to a group or single row. Default is false
// group_level (number) – Specify the group level to be used. Optional
// include_docs (boolean) – Include the associated document with each row. Default is false.
// attachments (boolean) – Include the Base64-encoded content of attachments in the documents that are included if include_docs is true. Ignored if include_docs isn’t true. Default is false.
// att_encoding_info (boolean) – Include encoding information in attachment stubs if include_docs is true and the particular attachment is compressed. Ignored if include_docs isn’t true. Default is false.
// inclusive_end (boolean) – Specifies whether the specified end key should be included in the result. Default is true
// key (json) – Return only documents that match the specified key. Optional
// keys (json-array) – Return only documents where the key matches one of the keys specified in the array. Optional
// limit (number) – Limit the number of the returned documents to the specified number. Optional
// reduce (boolean) – Use the reduction function. Default is true
// skip (number) – Skip this number of records before starting to return the results. Default is 0
// sorted (boolean) – Sort returned rows (see Sorting Returned Rows). Setting this to false offers a performance boost. The total_rows and offset fields are not available when this is set to false. Default is true
// stale (string) – Allow the results from a stale view to be used. Supported values: ok and update_after. Optional
// startkey (json) – Return records starting with the specified key. Optional
// start_key (json) – Alias for startkey param
// startkey_docid (string) – Return records starting with the specified document ID. Requires startkey to be specified for this to have any effect. Optional
// start_key_doc_id (string) – Alias for startkey_docid param
// update_seq (boolean) – Response includes an update_seq value indicating which sequence id of the database the view reflects. Default is false
type ViewParams struct {
	Conflicts              bool          `url:"conflicts,omitempty"`
	Descending             bool          `url:"descending,omitempty"`
	EndKey                 interface{}   `url:"endkey,omitempty"`
	EndKeyDocID            string        `url:"endkey_docid,omitempty"`
	Group                  bool          `url:"group,omitempty"`
	GroupLevel             float64       `url:"group_level,omitempty"`
	IncludeDocs            bool          `url:"include_docs,omitempty"`
	Attachments            bool          `url:"attachments,omitempty"`
	AttachmentEncodingInfo bool          `url:"att_encoding_info,omitempty"`
	InclusiveEnd           bool          `url:"inclusive_end,omitempty"`
	Key                    interface{}   `url:"key,omitempty"`
	Keys                   []interface{} `url:"keys,omitempty"`
	Limit                  float64       `url:"limit,omitempty"`
	Reduce                 bool          `url:"reduce,omitempty"`
	Skip                   float64       `url:"skip,omitempty"`
	Sorted                 bool          `url:"sorted,omitempty"`
	Stale                  string        `url:"stale,omitempty"`
	StartKey               interface{}   `url:"startkey,omitempty"`
	StartKeyDocID          string        `url:"startkey_docid,omitempty"`
	UpdateSeq              bool          `url:"update_seq,omitempty"`
}

type View interface {
	Execute(Options) (DocumentList, error)
}

type TemporaryView struct {
	Map    string `json:"map,omitempty"`
	Reduce string `json:"reduce,omitempty"`

	db *Database
}

func (d *Database) TemporaryView(mapFunc string) TemporaryView {
	return TemporaryView{
		Map: mapFunc,

		db: d,
	}
}

func (v TemporaryView) Execute(opts URLOptions) (DocumentList, error) {
	jsString, err := json.Marshal(v)
	if err != nil {
		return DocumentList{}, err
	}

	var docs DocumentList
	_, err = v.db.con.unmarshalRequest("POST", v.db.ViewPath("_temp_view"), opts, bytes.NewBuffer(jsString), &docs)
	if err != nil {
		return DocumentList{}, err
	}

	return docs, nil
}

type NamedView struct {
	DesignDoc string
	Name      string

	db *Database
}

func (d *Database) NamedView(design, name string) NamedView {
	return NamedView{
		DesignDoc: design,
		Name:      name,

		db: d,
	}
}

func (v NamedView) Execute(opts URLOptions) (DocumentList, error) {
	var docs DocumentList
	if _, err := v.db.con.unmarshalRequest("GET", v.db.ViewPath(v.Path()), opts, nil, &docs); err != nil {
		return DocumentList{}, err
	}

	return docs, nil
}

func (v NamedView) Path() string {
	return fmt.Sprintf("_design/%s/_view/%s", v.DesignDoc, v.Name)
}

func (v NamedView) FullPath() string {
	return urlConcat(v.db.Name(), v.Path())
}
