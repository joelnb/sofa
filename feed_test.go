package sofa

import (
	"fmt"
	"testing"
	"time"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

func TestPollingFeed(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Post("/test_db/_changes")
}

func TestFeedPollingReal(t *testing.T) {
	con := globalTestConnections.Version1(t, false)

	defer func() {
		if err := con.DeleteDatabase("feed_test_db"); err != nil {
			t.Fatal(err)
		}
	}()

	for _, lp := range []bool{true, false} {
		// Delete the DB if it currently exists
		if err := con.DeleteDatabase("feed_test_db"); err != nil {
			if !ErrorStatus(err, 404) {
				t.Fatal(err)
			}
		}

		db, err := con.CreateDatabase("feed_test_db")
		st.Assert(t, err, nil)

		feed := db.PollingChangesFeed(lp)

		// This is not true if we are doing long polling - it would just block
		// until a change is made (which would be never).
		if !lp {
			emptyUpdate, err := feed.Next(ChangesFeedParams{})
			st.Assert(t, err, nil)

			st.Expect(t, emptyUpdate.LastSeq, int64(0))
			st.Expect(t, emptyUpdate.Pending, int64(0))
			st.Expect(t, len(emptyUpdate.Results), 0)
		}

		doc := &struct {
			DocumentMetadata
			Name string `json:"name"`
			Type string `json:"type"`
		}{
			DocumentMetadata: DocumentMetadata{
				ID: "fruit1",
			},
			Name: "apple",
			Type: "fruit",
		}

		_, err = db.Put(doc)
		st.Assert(t, err, nil)

		middleUpdate, err := feed.Next(ChangesFeedParams{})
		st.Assert(t, err, nil)

		st.Expect(t, middleUpdate.LastSeq, int64(1))
		st.Assert(t, middleUpdate.Pending, int64(0))

		st.Assert(t, middleUpdate.Results[0].Deleted, false)
		st.Assert(t, middleUpdate.Results[0].ID, "fruit1")
		st.Assert(t, middleUpdate.Results[0].Seq, int64(1))

		_, err = db.Get(doc, doc.ID, "")
		st.Assert(t, err, nil)

		_, err = db.Delete(doc)
		st.Assert(t, err, nil)

		lastUpdate, err := feed.Next(ChangesFeedParams{})
		st.Assert(t, err, nil)

		st.Assert(t, lastUpdate.LastSeq, int64(2))
		st.Assert(t, lastUpdate.Pending, int64(0))

		st.Assert(t, lastUpdate.Results[0].Deleted, true)
		st.Assert(t, lastUpdate.Results[0].ID, "fruit1")
		st.Assert(t, lastUpdate.Results[0].Seq, int64(2))

		updateSince, err := feed.Next(ChangesFeedParams{
			Since: middleUpdate.LastSeq,
		})
		st.Assert(t, err, nil)

		st.Assert(t, lastUpdate.LastSeq, updateSince.LastSeq)
	}
}

func TestFeedContinuousReal(t *testing.T) {
	con := globalTestConnections.Version1(t, false)

	// Delete the DB if it currently exists
	if err := con.DeleteDatabase("feed_test_db"); err != nil {
		if !ErrorStatus(err, 404) {
			t.Fatal(err)
		}
	}

	db, err := con.CreateDatabase("feed_test_db")
	st.Assert(t, err, nil)

	defer func() {
		if err := con.DeleteDatabase("feed_test_db"); err != nil {
			t.Fatal(err)
		}
	}()

	feed := db.ContinuousChangesFeed(ChangesFeedParams{})

	doc := &struct {
		DocumentMetadata
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		DocumentMetadata: DocumentMetadata{
			ID: "fruit1",
		},
		Name: "apple",
		Type: "fruit",
	}

	_, err = db.Put(doc)
	st.Assert(t, err, nil)

	middleUpdate, err := feed.Next()
	st.Assert(t, err, nil)

	st.Assert(t, middleUpdate.Deleted, false)
	st.Assert(t, middleUpdate.ID, "fruit1")
	st.Assert(t, middleUpdate.Seq, int64(1))

	_, err = db.Get(doc, doc.ID, "")
	st.Assert(t, err, nil)

	_, err = db.Delete(doc)
	st.Assert(t, err, nil)

	lastUpdate, err := feed.Next()
	st.Assert(t, err, nil)

	st.Assert(t, lastUpdate.Deleted, true)
	st.Assert(t, lastUpdate.ID, "fruit1")
	st.Assert(t, lastUpdate.Seq, int64(2))

	otherDoc := &struct {
		DocumentMetadata
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		DocumentMetadata: DocumentMetadata{
			ID: "fruit2",
		},
		Name: "papaya",
		Type: "fruit",
	}

	go func() {
		time.Sleep(time.Millisecond * 500)
		_, err := db.Put(otherDoc)
		st.Assert(t, err, nil)
	}()

	updateAsync, err := feed.Next()
	st.Assert(t, err, nil)

	st.Assert(t, updateAsync.Deleted, false)
	st.Assert(t, updateAsync.ID, "fruit2")
	st.Assert(t, updateAsync.Seq, int64(3))
}
