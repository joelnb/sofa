package sofa

import (
	"fmt"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
)

// TestConnectionDatabase tests that a Database can successfully be retrieved from a Connection
func TestConnectionDatabase(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", TestMockHost)).
		Get("/test_db").
		Reply(200).
		JSON(map[string]interface{}{
			"db_name":              "test_db",
			"doc_count":            639,
			"doc_del_count":        10,
			"update_seq":           901,
			"purge_seq":            0,
			"compact_running":      false,
			"disk_size":            3427544,
			"data_size":            1563132,
			"instance_start_time":  "1484039376767413",
			"disk_format_version":  6,
			"committed_update_seq": 2913,
		})

	con := defaultMockTestConnection(t)
	db := con.Database("test_db")

	st.Assert(t, db.metadata, ((*DatabaseMetadata)(nil)))

	mdata, err := db.Metadata()
	st.Assert(t, err, nil)

	st.Reject(t, db.metadata, ((*DatabaseMetadata)(nil)))

	st.Assert(t, mdata.Name, "test_db")
	st.Assert(t, mdata.DocCount, 639)
	st.Assert(t, mdata.InstanceStartTime, "1484039376767413")
}

func TestConnectionEnsureDatabase(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", TestMockHost)).
		Get("/test_db").
		Reply(200).
		JSON(map[string]interface{}{
			"db_name":              "test_db",
			"doc_count":            639,
			"doc_del_count":        10,
			"update_seq":           901,
			"purge_seq":            0,
			"compact_running":      false,
			"disk_size":            3427544,
			"data_size":            1563132,
			"instance_start_time":  "1484039376767413",
			"disk_format_version":  6,
			"committed_update_seq": 2913,
		})

	con := defaultMockTestConnection(t)

	db, err := con.EnsureDatabase("test_db")
	st.Assert(t, err, nil)

	st.Reject(t, db.metadata, ((*DatabaseMetadata)(nil)))

	mdata, err := db.Metadata()
	st.Assert(t, err, nil)

	st.Assert(t, mdata.Name, "test_db")
}

func TestConnectionPing(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", TestMockHost)).
		Head("/").
		Reply(200)

	con := defaultMockTestConnection(t)

	if err := con.Ping(); err != nil {
		t.Fatal(err)
	}
}

func TestConnectionServerInfo(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", TestMockHost)).
		Get("/").
		Reply(200).
		JSON(map[string]interface{}{
			"couchdb": "Welcome",
			"uuid":    "038bbc5ae438344067d6ab90fe30ed05",
			"version": "1.6.1",
			"vendor":  map[string]string{"name": "The Apache Software Foundation", "version": "1.6.1"},
		})

	con := defaultMockTestConnection(t)

	info, err := con.ServerInfo()
	st.Assert(t, err, nil)

	st.Assert(t, info.CouchDB, "Welcome")
	st.Assert(t, info.UUID, "038bbc5ae438344067d6ab90fe30ed05")
	st.Assert(t, info.Version, "1.6.1")
	st.Assert(t, info.Vendor["name"].(string), "The Apache Software Foundation")
	st.Assert(t, info.Vendor["version"].(string), "1.6.1")
}

func TestConnectionListDatabases(t *testing.T) {
	defer gock.Off()

	testDBs := []string{
		"fruits",
		"testdb",
	}

	gock.New(fmt.Sprintf("https://%s", TestMockHost)).
		Get("/_all_dbs").
		Reply(200).
		JSON(testDBs)

	con := defaultMockTestConnection(t)

	serverDBs, err := con.ListDatabases()
	st.Assert(t, err, nil)

	for i, name := range testDBs {
		st.Assert(t, serverDBs[i], name)
	}
}

func TestConnectionReal(t *testing.T) {
	serverRequired(t)
	con := defaultTestConnection(t)

	info, err := con.ServerInfo()
	st.Assert(t, err, nil)

	st.Assert(t, info.CouchDB, "Welcome")

	db := initTestDB(t, con, "test_db")

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

	_, err = db.Get(doc, doc.ID, "")
	st.Assert(t, err, nil)

	_, err = db.Delete(doc)
	st.Assert(t, err, nil)

	cleanupTestDB(t, con, "test_db")
}
