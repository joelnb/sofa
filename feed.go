package sofa

import (
	"bufio"
	"encoding/json"
	"net/http"

	"github.com/google/go-querystring/query"
)

type changesFeedType string

const (
	FeedPolling     changesFeedType = "normal"
	FeedLongPolling changesFeedType = "longpoll"
	FeedContinuous  changesFeedType = "continuous"
	FeedEventSource changesFeedType = "eventsource"
)

type ChangesFeed interface {
	Next(ChangesFeedParams) (ChangesFeedUpdate, error)
}

type ChangesFeedUpdate struct {
	LastSeq int64               `json:"last_seq"`
	Pending int64               `json:"pending"`
	Results []ChangesFeedChange `json:"results"`
}

type ChangesFeedChange struct {
	Changes []struct {
		Rev string `json:"rev"`
	} `json:"changes"`
	Deleted bool   `json:"deleted"`
	ID      string `json:"id"`
	Seq     int64  `json:"seq"`
}

type ChangesFeedParams struct {
	DocumentIDs            []string `url:"doc_ids,omitempty"`
	Conflicts              bool     `url:"conflicts,omitempty"`
	Descending             bool     `url:"descending,omitempty"`
	Feed                   string   `url:"feed,omitempty"`
	Filter                 string   `url:"filter,omitempty"`
	Heartbeat              int64    `url:"heartbeat,omitempty"`
	IncludeDocs            bool     `url:"include_docs,omitempty"`
	Attachments            bool     `url:"attachments,omitempty"`
	AttachmentEncodingInfo bool     `url:"att_encoding_info,omitempty"`
	LastEventID            int64    `url:"last-event-id,omitempty"`
	Limit                  int64    `url:"limit,omitempty"`
	Since                  int64    `url:"since,omitempty"`
	Style                  string   `url:"style,omitempty"`
	Timeout                int64    `url:"timeout,omitempty"`
	View                   string   `url:"view,omitempty"`
}

type PollingChangesFeed struct {
	db       *Database
	feedType string
}

func (f PollingChangesFeed) Next(params ChangesFeedParams) (ChangesFeedUpdate, error) {
	params.Feed = f.feedType

	v, err := query.Values(params)
	if err != nil {
		return ChangesFeedUpdate{}, err
	}

	var u ChangesFeedUpdate
	_, err = f.db.con.unmarshalRequest("GET", f.db.ViewPath("_changes"), v, nil, &u)
	return u, err
}

type ContinuousChangesFeed struct {
	db     *Database
	params ChangesFeedParams

	resp    *http.Response
	scanner *bufio.Scanner
}

func (f *ContinuousChangesFeed) Next() (ChangesFeedChange, error) {
	if f.resp == nil {
		f.params.Feed = string(FeedContinuous)

		v, err := query.Values(f.params)
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

	var u ChangesFeedChange
	if err := json.Unmarshal(f.scanner.Bytes(), &u); err != nil {
		return ChangesFeedChange{}, err
	}

	return u, nil
}
