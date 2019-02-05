package sofa

import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

type changesFeedType string

const (
	// FeedPolling represents a "normal" feed in CouchDB which regulargly polls
	// the server for updates.
	FeedPolling changesFeedType = "normal"
	// FeedLongPolling represents the type of feed which uses long-polling to reduce
	// polling frequency & therefore requests to the server.
	FeedLongPolling changesFeedType = "longpoll"
	// FeedContinuous represents the type of feed which holds open a connection to
	// the server and receives a stream of events.
	FeedContinuous changesFeedType = "continuous"
	// FeedEventSource represents the CouchDB feed type of "eventsource". This type is not
	// yet implemented in this library.
	FeedEventSource changesFeedType = "eventsource"
)

// ChangesFeed is an interface which is implemented by all types of changes feed which are
// available on the server.
type ChangesFeed interface {
	Next(ChangesFeedParams) (ChangesFeedUpdate, error)
}

// ChangesFeedUpdate is a single update from a changes feed.
type ChangesFeedUpdate struct {
	LastSeq int64               `json:"last_seq"`
	Pending int64               `json:"pending"`
	Results []ChangesFeedChange `json:"results"`
}

// ChangesFeedChange is a single change to a document in the database. One of
// more of these will be included in updates where any changes were actually made.
type ChangesFeedChange struct {
	Changes []struct {
		Rev string `json:"rev"`
	} `json:"changes"`
	Deleted bool   `json:"deleted"`
	ID      string `json:"id"`
	Seq     int64  `json:"seq"`
}

// ChangesFeedParams includes all parameters which can be provided to
// control the output of a changes feed from a database.
type ChangesFeedParams struct {
	DocumentIDs            []string         `url:"doc_ids,omitempty"`
	Conflicts              BooleanParameter `url:"conflicts,omitempty"`
	Descending             BooleanParameter `url:"descending,omitempty"`
	Feed                   string           `url:"feed,omitempty"`
	Filter                 string           `url:"filter,omitempty"`
	Heartbeat              int64            `url:"heartbeat,omitempty"`
	IncludeDocs            BooleanParameter `url:"include_docs,omitempty"`
	Attachments            BooleanParameter `url:"attachments,omitempty"`
	AttachmentEncodingInfo BooleanParameter `url:"att_encoding_info,omitempty"`
	LastEventID            int64            `url:"last-event-id,omitempty"`
	Limit                  int64            `url:"limit,omitempty"`
	Since                  int64            `url:"since,omitempty"`
	Style                  string           `url:"style,omitempty"`
	Timeout                int64            `url:"timeout,omitempty"`
	View                   string           `url:"view,omitempty"`
}

func (params ChangesFeedParams) Values() (url.Values, error) {
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	if params.DocumentIDs != nil {
		jBytes, err := json.Marshal(params.DocumentIDs)
		if err != nil {
			return nil, err
		}

		v.Set("doc_ids", string(jBytes))
	}

	return v, nil
}

// PollingChangesFeed can be used for either type of changes feed which polls
// the database for information ("normal" and "longpoll").
type PollingChangesFeed struct {
	db       *Database
	feedType string
}

// Next polls for the next update from the database. This may block until a timeout is
// reached if there are no updates available.
func (f PollingChangesFeed) Next(params ChangesFeedParams) (ChangesFeedUpdate, error) {
	params.Feed = f.feedType

	v, err := params.Values()
	if err != nil {
		return ChangesFeedUpdate{}, err
	}

	var u ChangesFeedUpdate
	_, err = f.db.con.unmarshalRequest("GET", f.db.ViewPath("_changes"), v, nil, &u)
	return u, err
}

// ContinuousChangesFeed maintains a connection to the database and receives continuous
// updates as they arrive.
type ContinuousChangesFeed struct {
	db     *Database
	params ChangesFeedParams

	resp    *http.Response
	scanner *bufio.Scanner
}

// Next gets the next available item from the changes feed. This will block until an item
// becomes available.
func (f *ContinuousChangesFeed) Next() (ChangesFeedChange, error) {
	if f.resp == nil {
		f.params.Feed = string(FeedContinuous)

		v, err := f.params.Values()
		if err != nil {
			return ChangesFeedChange{}, err
		}

		resp, err := f.db.con.urlRequest("GET", f.db.con.URL(f.db.ViewPath("_changes")), v, nil, false)
		if err != nil {
			return ChangesFeedChange{}, err
		}

		f.resp = resp
		f.scanner = bufio.NewScanner(f.resp.Body)
	}

	if !f.scanner.Scan() {
		return ChangesFeedChange{}, f.scanner.Err()
	}

	// Swallow an extra newline if needed
	if len(f.scanner.Bytes()) == 0 {
		if !f.scanner.Scan() {
			return ChangesFeedChange{}, f.scanner.Err()
		}
	}

	var u ChangesFeedChange
	if err := json.Unmarshal(f.scanner.Bytes(), &u); err != nil {
		return ChangesFeedChange{}, err
	}

	return u, nil
}
