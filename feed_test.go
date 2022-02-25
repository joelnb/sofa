package sofa

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

func getFeedTestOtherDoc() Document {
	return &struct {
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
}

func TestPollingFeed(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Post("/test_db/_changes")
}

func TestFeedPollingReal1(t *testing.T) {
	con := globalTestConnections.Version1(t, false)

	defer func() {
		cleanTestDB(t, con.Connection, "feed_test_db", false)
	}()

	for _, lp := range []bool{true, false} {
		// Delete the DB if it currently exists
		cleanTestDB(t, con.Connection, "feed_test_db", true)

		db, err := con.CreateDatabase("feed_test_db")
		st.Assert(t, err, nil)

		feed := db.PollingChangesFeed(lp)

		// This is not true if we are doing long polling - it would just block
		// until a change is made (which would be never).
		if !lp {
			emptyUpdate, err := feed.Next(&ChangesFeedParams1{})
			st.Assert(t, err, nil)

			st.Expect(t, emptyUpdate.LastSeq, AlwaysString("0"))
			st.Expect(t, emptyUpdate.Pending, int64(0))
			st.Expect(t, len(emptyUpdate.Results), 0)
		}

		doc := getDefaultTestDoc()

		_, err = db.Put(doc)
		st.Assert(t, err, nil)

		middleUpdate, err := feed.Next(&ChangesFeedParams1{})
		st.Assert(t, err, nil)

		st.Expect(t, middleUpdate.LastSeq, AlwaysString("1"))
		st.Assert(t, middleUpdate.Pending, int64(0))

		st.Assert(t, middleUpdate.Results[0].Deleted, false)
		st.Assert(t, middleUpdate.Results[0].ID, "fruit1")
		st.Assert(t, middleUpdate.Results[0].Seq, AlwaysString("1"))

		_, err = db.Get(doc, doc.Metadata().ID, "")
		st.Assert(t, err, nil)

		_, err = db.Delete(doc)
		st.Assert(t, err, nil)

		lastUpdate, err := feed.Next(&ChangesFeedParams1{})
		st.Assert(t, err, nil)

		st.Assert(t, lastUpdate.LastSeq, AlwaysString("2"))
		st.Assert(t, lastUpdate.Pending, int64(0))

		st.Assert(t, lastUpdate.Results[0].Deleted, true)
		st.Assert(t, lastUpdate.Results[0].ID, "fruit1")
		st.Assert(t, lastUpdate.Results[0].Seq, AlwaysString("2"))

		since, err := strconv.Atoi(string(middleUpdate.LastSeq))
		st.Assert(t, err, nil)

		updateSince, err := feed.Next(&ChangesFeedParams1{
			Since: int64(since),
		})
		st.Assert(t, err, nil)

		st.Assert(t, lastUpdate.LastSeq, updateSince.LastSeq)
	}
}

func TestFeedContinuousReal1(t *testing.T) {
	con := globalTestConnections.Version1(t, false)

	// Delete the DB if it currently exists
	cleanTestDB(t, con.Connection, "feed_test_db", true)

	db, err := con.CreateDatabase("feed_test_db")
	st.Assert(t, err, nil)

	defer func() {
		cleanTestDB(t, con.Connection, "feed_test_db", false)
	}()

	feed := db.ContinuousChangesFeed(&ChangesFeedParams1{})

	doc := getDefaultTestDoc()

	_, err = db.Put(doc)
	st.Assert(t, err, nil)

	middleUpdate, err := feed.Next()
	st.Assert(t, err, nil)

	st.Assert(t, middleUpdate.Deleted, false)
	st.Assert(t, middleUpdate.ID, "fruit1")
	st.Assert(t, middleUpdate.Seq, AlwaysString("1"))

	_, err = db.Get(doc, doc.Metadata().ID, "")
	st.Assert(t, err, nil)

	_, err = db.Delete(doc)
	st.Assert(t, err, nil)

	lastUpdate, err := feed.Next()
	st.Assert(t, err, nil)

	st.Assert(t, lastUpdate.Deleted, true)
	st.Assert(t, lastUpdate.ID, "fruit1")
	st.Assert(t, lastUpdate.Seq, AlwaysString("2"))

	otherDoc := getFeedTestOtherDoc()

	go func() {
		time.Sleep(time.Millisecond * 500)
		_, err := db.Put(otherDoc)
		st.Assert(t, err, nil)
	}()

	updateAsync, err := feed.Next()
	st.Assert(t, err, nil)

	st.Assert(t, updateAsync.Deleted, false)
	st.Assert(t, updateAsync.ID, "fruit2")
	st.Assert(t, updateAsync.Seq, AlwaysString("3"))
}

func TestFeedPollingReal2(t *testing.T) {
	con := globalTestConnections.Version2(t, false)

	defer func() {
		cleanTestDB(t, con.Connection, "feed_test_db", false)
	}()

	for _, lp := range []bool{true, false} {
		// Delete the DB if it currently exists
		cleanTestDB(t, con.Connection, "feed_test_db", true)

		db, err := con.CreateDatabase("feed_test_db")
		st.Assert(t, err, nil)

		feed := db.PollingChangesFeed(lp)

		// This is not true if we are doing long polling - it would just block
		// until a change is made (which would be never).
		if !lp {
			emptyUpdate, err := feed.Next(&ChangesFeedParams2{})
			st.Assert(t, err, nil)

			assertPrefix(t, string(emptyUpdate.LastSeq), "0-")
			st.Expect(t, emptyUpdate.Pending, int64(0))
			st.Expect(t, len(emptyUpdate.Results), 0)
		}

		doc := getDefaultTestDoc()

		_, err = db.Put(doc)
		st.Assert(t, err, nil)

		middleUpdate, err := feed.Next(&ChangesFeedParams2{})
		st.Assert(t, err, nil)

		assertPrefix(t, string(middleUpdate.LastSeq), "1-")
		st.Assert(t, middleUpdate.Pending, int64(0))

		st.Assert(t, middleUpdate.Results[0].Deleted, false)
		st.Assert(t, middleUpdate.Results[0].ID, "fruit1")
		assertPrefix(t, string(middleUpdate.Results[0].Seq), "1-")

		_, err = db.Get(doc, doc.Metadata().ID, "")
		st.Assert(t, err, nil)

		_, err = db.Delete(doc)
		st.Assert(t, err, nil)

		lastUpdate, err := feed.Next(&ChangesFeedParams2{})
		st.Assert(t, err, nil)

		assertPrefix(t, string(lastUpdate.LastSeq), "2-")
		st.Assert(t, lastUpdate.Pending, int64(0))

		st.Assert(t, lastUpdate.Results[0].Deleted, true)
		st.Assert(t, lastUpdate.Results[0].ID, "fruit1")
		assertPrefix(t, string(lastUpdate.Results[0].Seq), "2-")

		updateSince, err := feed.Next(&ChangesFeedParams2{
			Since: string(middleUpdate.LastSeq),
		})
		st.Assert(t, err, nil)

		st.Assert(t, lastUpdate.LastSeq, updateSince.LastSeq)
	}
}

func TestFeedContinuousReal2(t *testing.T) {
	con := globalTestConnections.Version2(t, false)

	// Delete the DB if it currently exists
	cleanTestDB(t, con.Connection, "feed_test_db", true)

	db, err := con.CreateDatabase("feed_test_db")
	st.Assert(t, err, nil)

	defer func() {
		cleanTestDB(t, con.Connection, "feed_test_db", false)
	}()

	feed := db.ContinuousChangesFeed(&ChangesFeedParams2{})

	doc := getDefaultTestDoc()

	_, err = db.Put(doc)
	st.Assert(t, err, nil)

	middleUpdate, err := feed.Next()
	st.Assert(t, err, nil)

	st.Assert(t, middleUpdate.Deleted, false)
	st.Assert(t, middleUpdate.ID, "fruit1")
	assertPrefix(t, string(middleUpdate.Seq), "1-")

	_, err = db.Get(doc, doc.Metadata().ID, "")
	st.Assert(t, err, nil)

	_, err = db.Delete(doc)
	st.Assert(t, err, nil)

	lastUpdate, err := feed.Next()
	st.Assert(t, err, nil)

	st.Assert(t, lastUpdate.Deleted, true)
	st.Assert(t, lastUpdate.ID, "fruit1")
	assertPrefix(t, string(lastUpdate.Seq), "2-")

	otherDoc := getFeedTestOtherDoc()

	go func() {
		time.Sleep(time.Millisecond * 500)
		_, err := db.Put(otherDoc)
		st.Assert(t, err, nil)
	}()

	updateAsync, err := feed.Next()
	st.Assert(t, err, nil)

	st.Assert(t, updateAsync.Deleted, false)
	st.Assert(t, updateAsync.ID, "fruit2")
	assertPrefix(t, string(updateAsync.Seq), "3-")
}

func TestFeedPollingReal3(t *testing.T) {
	con := globalTestConnections.Version3(t, false)

	defer func() {
		cleanTestDB(t, con.Connection, "feed_test_db", false)
	}()

	for _, lp := range []bool{true, false} {
		// Delete the DB if it currently exists
		cleanTestDB(t, con.Connection, "feed_test_db", true)

		db, err := con.CreateDatabase("feed_test_db")
		st.Assert(t, err, nil)

		feed := db.PollingChangesFeed(lp)

		// This is not true if we are doing long polling - it would just block
		// until a change is made (which would be never).
		if !lp {
			emptyUpdate, err := feed.Next(&ChangesFeedParams2{})
			st.Assert(t, err, nil)

			assertPrefix(t, string(emptyUpdate.LastSeq), "0-")
			st.Expect(t, emptyUpdate.Pending, int64(0))
			st.Expect(t, len(emptyUpdate.Results), 0)
		}

		doc := getDefaultTestDoc()

		_, err = db.Put(doc)
		st.Assert(t, err, nil)

		middleUpdate, err := feed.Next(&ChangesFeedParams2{})
		st.Assert(t, err, nil)

		assertPrefix(t, string(middleUpdate.LastSeq), "1-")
		st.Assert(t, middleUpdate.Pending, int64(0))

		st.Assert(t, middleUpdate.Results[0].Deleted, false)
		st.Assert(t, middleUpdate.Results[0].ID, "fruit1")
		assertPrefix(t, string(middleUpdate.Results[0].Seq), "1-")

		_, err = db.Get(doc, doc.Metadata().ID, "")
		st.Assert(t, err, nil)

		_, err = db.Delete(doc)
		st.Assert(t, err, nil)

		lastUpdate, err := feed.Next(&ChangesFeedParams2{})
		st.Assert(t, err, nil)

		assertPrefix(t, string(lastUpdate.LastSeq), "2-")
		st.Assert(t, lastUpdate.Pending, int64(0))

		st.Assert(t, lastUpdate.Results[0].Deleted, true)
		st.Assert(t, lastUpdate.Results[0].ID, "fruit1")
		assertPrefix(t, string(lastUpdate.Results[0].Seq), "2-")

		updateSince, err := feed.Next(&ChangesFeedParams2{
			Since: string(middleUpdate.LastSeq),
		})
		st.Assert(t, err, nil)

		st.Assert(t, lastUpdate.LastSeq, updateSince.LastSeq)
	}
}

func TestFeedContinuousReal3(t *testing.T) {
	con := globalTestConnections.Version3(t, false)

	// Delete the DB if it currently exists
	cleanTestDB(t, con.Connection, "feed_test_db", true)

	db, err := con.CreateDatabase("feed_test_db")
	st.Assert(t, err, nil)

	defer func() {
		cleanTestDB(t, con.Connection, "feed_test_db", false)
	}()

	feed := db.ContinuousChangesFeed(&ChangesFeedParams3{})

	doc := getDefaultTestDoc()

	_, err = db.Put(doc)
	st.Assert(t, err, nil)

	middleUpdate, err := feed.Next()
	st.Assert(t, err, nil)

	st.Assert(t, middleUpdate.Deleted, false)
	st.Assert(t, middleUpdate.ID, "fruit1")
	assertPrefix(t, string(middleUpdate.Seq), "1-")

	_, err = db.Get(doc, doc.Metadata().ID, "")
	st.Assert(t, err, nil)

	_, err = db.Delete(doc)
	st.Assert(t, err, nil)

	lastUpdate, err := feed.Next()
	st.Assert(t, err, nil)

	st.Assert(t, lastUpdate.Deleted, true)
	st.Assert(t, lastUpdate.ID, "fruit1")
	assertPrefix(t, string(lastUpdate.Seq), "2-")

	otherDoc := getFeedTestOtherDoc()

	go func() {
		time.Sleep(time.Millisecond * 500)
		_, err := db.Put(otherDoc)
		st.Assert(t, err, nil)
	}()

	updateAsync, err := feed.Next()
	st.Assert(t, err, nil)

	st.Assert(t, updateAsync.Deleted, false)
	st.Assert(t, updateAsync.ID, "fruit2")
	assertPrefix(t, string(updateAsync.Seq), "3-")
}
