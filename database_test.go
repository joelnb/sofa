package sofa

import (
	"fmt"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

func TestDatabaseGet(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Get("/test_db/somedoc").
		Reply(200).
		JSON(map[string]string{
			"_id":     "somedoc",
			"_rev":    "51-f6d5f19782344f44876a9d775376c549",
			"content": "Hello, Tester!",
		}).
		SetHeader("Etag", `"51-f6d5f19782344f44876a9d775376c549"`)

	con := globalTestConnections.Version1(t, true)
	db := con.Database("test_db")

	doc := struct {
		DocumentMetadata
		Content string `json:"content"`
	}{}

	rev, err := db.Get(&doc, "somedoc", "")
	st.Assert(t, err, nil)

	st.Assert(t, rev, "51-f6d5f19782344f44876a9d775376c549")
}

func TestDatabasePut(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Put("/test_db/newdoc").
		Reply(201).
		JSON(map[string]interface{}{
			"ok":  true,
			"id":  "newdoc",
			"rev": "2-801609c9fdb4c6d196820c5b1f3c26c9",
		}).
		SetHeader("Etag", `"2-801609c9fdb4c6d196820c5b1f3c26c9"`)

	con := globalTestConnections.Version1(t, true)
	db := con.Database("test_db")

	doc := &struct {
		DocumentMetadata
		Content string `json:"content"`
	}{
		DocumentMetadata: DocumentMetadata{
			ID:  "newdoc",
			Rev: "1-5bfa2c99eefe2b2eb4962db50aa3cfd4",
		},
		Content: "Not included",
	}

	rev, err := db.Put(doc)
	st.Assert(t, err, nil)

	st.Assert(t, rev, "2-801609c9fdb4c6d196820c5b1f3c26c9")
}

func TestDatabaseList(t *testing.T) {
	defer gock.Off()

	expectedRows := []Row{
		{
			ID:    "newdoc",
			Key:   "newdoc",
			Value: map[string]interface{}{"rev": "1-87ae7d46fb1561c570cc5d7ce5c80c1e"},
		},
		{
			ID:    "olddoc",
			Key:   "olddoc",
			Value: map[string]interface{}{"rev": "1-40bf1b639fdb5a345dcf519399a431f0"},
		},
	}

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Get("/list-db_23456/_all_docs").
		Reply(200).
		JSON(map[string]interface{}{
			"total_rows": 2,
			"offset":     0,
			"rows": []map[string]interface{}{
				{
					"id":    "newdoc",
					"key":   "newdoc",
					"value": map[string]string{"rev": "1-87ae7d46fb1561c570cc5d7ce5c80c1e"},
				},
				{
					"id":    "olddoc",
					"key":   "olddoc",
					"value": map[string]string{"rev": "1-40bf1b639fdb5a345dcf519399a431f0"},
				},
			},
		})

	con := globalTestConnections.Version1(t, true)
	db := con.Database("list-db_23456")

	docs, err := db.ListDocuments()
	st.Assert(t, err, nil)

	if len(expectedRows) != len(docs.Rows) {
		t.Fatalf("length of expected list & actual list not equal: %d != %d", len(expectedRows), len(docs.Rows))
	}

	for i, row := range docs.Rows {
		expected := expectedRows[i]

		st.Assert(t, row.ID, expected.ID)
		st.Assert(t, row.Key, expected.Key)
		st.Assert(t, row.Value, expected.Value)
		st.Assert(t, row.Document, expected.Document)
	}
}

func TestDatabaseReal1(t *testing.T) {
	con := globalTestConnections.Version1(t, false)

	cleanTestDB(t, con.Connection, "test_db", true)

	db, err := con.CreateDatabase("test_db")
	st.Assert(t, err, nil)

	metadata, err := db.Metadata()
	st.Assert(t, err, nil)
	st.Assert(t, metadata.Name, "test_db")

	cleanTestDB(t, con.Connection, "test_db", false)
}

func TestDatabaseReal2(t *testing.T) {
	con := globalTestConnections.Version2(t, false)

	cleanTestDB(t, con.Connection, "test_db", true)

	db, err := con.CreateDatabase("test_db")
	st.Assert(t, err, nil)

	metadata, err := db.Metadata()
	st.Assert(t, err, nil)
	st.Assert(t, metadata.Name, "test_db")

	cleanTestDB(t, con.Connection, "test_db", false)
}

func TestDatabaseReal3(t *testing.T) {
	con := globalTestConnections.Version3(t, false)

	cleanTestDB(t, con.Connection, "test_db", true)

	db, err := con.CreateDatabase("test_db")
	st.Assert(t, err, nil)

	metadata, err := db.Metadata()
	st.Assert(t, err, nil)
	st.Assert(t, metadata.Name, "test_db")

	cleanTestDB(t, con.Connection, "test_db", false)
}
