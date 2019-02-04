package sofa

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	Conflicts              BooleanParameter `url:"conflicts,omitempty"`
	Descending             BooleanParameter `url:"descending,omitempty"`
	EndKey                 interface{}      `url:"endkey,omitempty"`
	EndKeyDocID            string           `url:"endkey_docid,omitempty"`
	Group                  BooleanParameter `url:"group,omitempty"`
	GroupLevel             float64          `url:"group_level,omitempty"`
	IncludeDocs            BooleanParameter `url:"include_docs,omitempty"`
	Attachments            BooleanParameter `url:"attachments,omitempty"`
	AttachmentEncodingInfo BooleanParameter `url:"att_encoding_info,omitempty"`
	InclusiveEnd           BooleanParameter `url:"inclusive_end,omitempty"`
	Key                    interface{}      `url:"key,omitempty"`
	Keys                   []interface{}    `url:"keys,omitempty"`
	Limit                  float64          `url:"limit,omitempty"`
	Reduce                 BooleanParameter `url:"reduce,omitempty"`
	Skip                   float64          `url:"skip,omitempty"`
	Sorted                 BooleanParameter `url:"sorted,omitempty"`
	Stale                  string           `url:"stale,omitempty"`
	StartKey               interface{}      `url:"startkey,omitempty"`
	StartKeyDocID          string           `url:"startkey_docid,omitempty"`
	UpdateSeq              BooleanParameter `url:"update_seq,omitempty"`
}

func (v *ViewParams) URLOptions() (*URLOptions, error) {
	u := NewURLOptions()

	// Process booleans
	if v.Conflicts != Empty {
		if err := u.Set("conflicts", v.Conflicts); err != nil {
			return nil, err
		}
	}
	if v.Descending != Empty {
		if err := u.Set("descending", v.Descending); err != nil {
			return nil, err
		}
	}
	if v.Group != Empty {
		if err := u.Set("group", v.Group); err != nil {
			return nil, err
		}
	}
	if v.IncludeDocs != Empty {
		if err := u.Set("include_docs", v.IncludeDocs); err != nil {
			return nil, err
		}
	}
	if v.Attachments != Empty {
		if err := u.Set("attachments", v.Attachments); err != nil {
			return nil, err
		}
	}
	if v.AttachmentEncodingInfo != Empty {
		if err := u.Set("att_encoding_info", v.AttachmentEncodingInfo); err != nil {
			return nil, err
		}
	}
	if v.InclusiveEnd != Empty {
		if err := u.Set("inclusive_end", v.InclusiveEnd); err != nil {
			return nil, err
		}
	}
	if v.Reduce != Empty {
		if err := u.Set("reduce", v.Reduce); err != nil {
			return nil, err
		}
	}
	if v.Sorted != Empty {
		if err := u.Set("sorted", v.Sorted); err != nil {
			return nil, err
		}
	}
	if v.UpdateSeq != Empty {
		if err := u.Set("update_seq", v.UpdateSeq); err != nil {
			return nil, err
		}
	}

	// Process interfaces
	if v.StartKey != nil {
		val := v.StartKey
		if err := u.Set("startkey", val); err != nil {
			return nil, err
		}
	}
	if v.EndKey != nil {
		val := v.EndKey
		if err := u.Set("endkey", val); err != nil {
			return nil, err
		}
	}
	if v.Key != nil {
		val := v.Key
		if err := u.Set("key", val); err != nil {
			return nil, err
		}
	}

	// Process lists
	if v.Keys != nil {
		if err := u.Set("keys", v.Keys); err != nil {
			return nil, err
		}
	}

	// Process strings
	if v.EndKeyDocID != "" {
		if err := u.Set("endkey_docid", v.EndKeyDocID); err != nil {
			return nil, err
		}
	}
	if v.StartKeyDocID != "" {
		if err := u.Set("startkey_docid", v.StartKeyDocID); err != nil {
			return nil, err
		}
	}
	if v.Stale != "" {
		if err := u.Set("stale", v.Stale); err != nil {
			return nil, err
		}
	}

	// Process floats
	// TODO: Can something better be done that checking for zero?
	if v.GroupLevel != 0 {
		if err := u.Set("group_level", v.GroupLevel); err != nil {
			return nil, err
		}
	}
	if v.Limit != 0 {
		if err := u.Set("limit", v.Limit); err != nil {
			return nil, err
		}
	}
	if v.Skip != 0 {
		if err := u.Set("skip", v.Skip); err != nil {
			return nil, err
		}
	}

	return &u, nil
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

	opts, err := params.URLOptions()
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
	opts, err := params.URLOptions()
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
