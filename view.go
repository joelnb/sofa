package sofa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
)

// BooleanParameter is a special type of boolean created to have a zero value
// where it is not included in URL parameter output. This is useful for taking
// the default values of a parameter.
type BooleanParameter string

func (b BooleanParameter) String() string {
	return string(b)
}

const (
	// Empty is the zero value for the BooleanParameter type. It is the default
	// type and values of this type are not included in the URL parameters.
	Empty BooleanParameter = ""
	// True is the BooleanParameter equivalent of true. It will always be
	// included in a query string.
	True BooleanParameter = "true"
	// False is the BooleanParameter equivalent of true. It will always be
	// included in a query string.
	False BooleanParameter = "false"
)

// InterfaceParameter is a wrapper for an empty interface which ensures that it is correctly formatted when passed
// as a query parameter to CouchDB. The reason for this is that strings passed are usually required not to be quoted
// but in these fields which can also take JSON objects of other types the quotes seem required
type InterfaceParameter struct {
	innerVal interface{}
}

// NewInterfaceParameter returns a pointer to an InterfaceParameter wrapping the provided value. All new
// InterfaceParameters must be created through this function.
func NewInterfaceParameter(iface interface{}) *InterfaceParameter {
	return &InterfaceParameter{
		innerVal: iface,
	}
}

// MarshalJSON simply marshals the internal value, to ensure these objects are always included using the
// JSON-formatted representation of just the inner value.
func (i InterfaceParameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.innerVal)
}

func (i InterfaceParameter) EncodeValues(key string, v *url.Values) error {
	bytes, err := json.Marshal(i.innerVal)
	if err != nil {
		return err
	}

	v.Set(key, string(bytes))

	return nil
}

type InterfaceListParameter struct {
	innerVal []*InterfaceParameter
}

func NewInterfaceListParameter(ifaces []interface{}) *InterfaceListParameter {
	iList := InterfaceListParameter{}
	iList.innerVal = []*InterfaceParameter{}

	for _, i := range ifaces {
		iList.innerVal = append(iList.innerVal, NewInterfaceParameter(i))
	}

	return &iList
}

// MarshalJSON simply marshals the internal value.
func (i InterfaceListParameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.innerVal)
}

// ViewParams provides a type-safe implementation of the paramaters which may
// be passed to an execution of a CouchDB view function:
//  - conflicts (boolean) – Includes conflicts information in response.
//    Ignored if include_docs isn’t true. Default is false
//  - descending (boolean) – Return the documents in descending by key
//    order. Default is false
//  - endkey (json) – Stop returning records when the specified key is
//    reached. Optional
//  - end_key (json) – Alias for endkey param
//  - endkey_docid (string) – Stop returning records when the specified
//    document ID is reached. Requires endkey to be specified for this to
//    have any effect. Optional
//  - end_key_doc_id (string) – Alias for endkey_docid param
//  - group (boolean) – Group the results using the reduce function to a
//    group or single row. Default is false
//  - group_level (number) – Specify the group level to be used. Optional
//  - include_docs (boolean) – Include the associated document with each row.
//    Default is false.
//  - attachments (boolean) – Include the Base64-encoded content of
//    attachments in the documents that are included if include_docs is true.
//    Ignored if include_docs isn’t true. Default is false.
//  - att_encoding_info (boolean) – Include encoding information in
//    attachment stubs if include_docs is true and the particular attachment
//    is compressed. Ignored if include_docs isn’t true. Default is false.
//  - inclusive_end (boolean) – Specifies whether the specified end key
//    should be included in the result. Default is true
//  - key (json) – Return only documents that match the specified key.
//    Optional
//  - keys (json-array) – Return only documents where the key matches one of
//    the keys specified in the array. Optional
//  - limit (number) – Limit the number of the returned documents to the
//    specified number. Optional
//  - reduce (boolean) – Use the reduction function. Default is true
//  - skip (number) – Skip this number of records before starting to return
//    the results. Default is 0
//  - sorted (boolean) – Sort returned rows (see Sorting Returned Rows).
//    Setting this to false offers a performance boost. The total_rows and
//    offset fields are not available when this is set to false.
//    Default is true
//  - stale (string) – Allow the results from a stale view to be used.
//    Supported values: ok and update_after. Optional
//  - startkey (json) – Return records starting with the specified key.
//    Optional
//  - start_key (json) – Alias for startkey param
//  - startkey_docid (string) – Return records starting with the specified
//    document ID. Requires startkey to be specified for this to have any
//    effect. Optional
//  - start_key_doc_id (string) – Alias for startkey_docid param
//  - update_seq (boolean) – Response includes an update_seq value
//    indicating which sequence id of the database the view reflects.
//    Default is false
type ViewParams struct {
	Conflicts              BooleanParameter        `url:"conflicts,omitempty"`
	Descending             BooleanParameter        `url:"descending,omitempty"`
	EndKey                 *InterfaceParameter     `url:"endkey,omitempty"`
	EndKeyDocID            string                  `url:"endkey_docid,omitempty"`
	Group                  BooleanParameter        `url:"group,omitempty"`
	GroupLevel             float64                 `url:"group_level,omitempty"`
	IncludeDocs            BooleanParameter        `url:"include_docs,omitempty"`
	Attachments            BooleanParameter        `url:"attachments,omitempty"`
	AttachmentEncodingInfo BooleanParameter        `url:"att_encoding_info,omitempty"`
	InclusiveEnd           BooleanParameter        `url:"inclusive_end,omitempty"`
	Key                    *InterfaceParameter     `url:"key,omitempty"`
	Keys                   *InterfaceListParameter `url:"keys,omitempty"`
	Limit                  float64                 `url:"limit,omitempty"`
	Reduce                 BooleanParameter        `url:"reduce,omitempty"`
	Skip                   float64                 `url:"skip,omitempty"`
	Sorted                 BooleanParameter        `url:"sorted,omitempty"`
	Stale                  string                  `url:"stale,omitempty"`
	StartKey               *InterfaceParameter     `url:"startkey,omitempty"`
	StartKeyDocID          string                  `url:"startkey_docid,omitempty"`
	UpdateSeq              BooleanParameter        `url:"update_seq,omitempty"`
}

func (v ViewParams) Values() (url.Values, error) {
	return query.Values(v)
}

// View is an interface representing the way that views are executed and
// their results returned.
type View interface {
	Execute(Options) (DocumentList, error)
}

// TemporaryView is a type of view which can be created & accessed on the fly.
// Temporary views are good for debugging purposed but should never be used in
// production as they are slow for any large number of documents.
type TemporaryView struct {
	Map    string `json:"map,omitempty"`
	Reduce string `json:"reduce,omitempty"`

	db *Database
}

// TemporaryView creates a temporary view for this database. Only the map function is
// required but other parameters canbe added to the resulting TemporaryView if
// required.
func (d *Database) TemporaryView(mapFunc string) TemporaryView {
	return TemporaryView{
		Map: mapFunc,

		db: d,
	}
}

// Execute implements View for TemporaryView.
func (v TemporaryView) Execute(params ViewParams) (DocumentList, error) {
	jsString, err := json.Marshal(v)
	if err != nil {
		return DocumentList{}, err
	}

	opts, err := params.Values()
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

// NamedView represents a view stored on a design document in the database.
// It must be accessed with both the name of the design document and the
// name of the view.
type NamedView struct {
	DesignDoc string
	Name      string

	db *Database
}

// NamedView creates a new NamedView for this database. This can then be used
// to access the current results of the permanent view on the design document.
func (d *Database) NamedView(design, name string) NamedView {
	return NamedView{
		DesignDoc: design,
		Name:      name,

		db: d,
	}
}

// Execute implements View for NamedView.
func (v NamedView) Execute(params ViewParams) (DocumentList, error) {
	opts, err := params.Values()
	if err != nil {
		return DocumentList{}, err
	}

	var docs DocumentList
	if _, err := v.db.con.unmarshalRequest("GET", v.db.ViewPath(v.Path()), opts, nil, &docs); err != nil {
		return DocumentList{}, err
	}

	return docs, nil
}

// Path gets the path of the NamedView relative to the database root.
func (v NamedView) Path() string {
	return fmt.Sprintf("_design/%s/_view/%s", v.DesignDoc, v.Name)
}

// FullPath gets the path of the NamedView relative to the server root.
func (v NamedView) FullPath() string {
	return urlConcat(v.db.Name(), v.Path())
}
