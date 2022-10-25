package sofa

import (
	"fmt"
	"testing"

	"github.com/nbio/st"
	"github.com/h2non/gock"
)

func getMockDatabaseRoot() map[string]interface{} {
	return map[string]interface{}{
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
	}
}

// TestConnectionDatabase tests that a Database can successfully be retrieved from a Connection
func TestConnectionDatabase(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Get("/test_db").
		Reply(200).
		JSON(getMockDatabaseRoot())

	con := globalTestConnections.Version1(t, true)
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

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Get("/test_db").
		Reply(200).
		JSON(getMockDatabaseRoot())

	con := globalTestConnections.Version1(t, true)

	db, err := con.EnsureDatabase("test_db")
	st.Assert(t, err, nil)

	st.Reject(t, db.metadata, ((*DatabaseMetadata)(nil)))

	mdata, err := db.Metadata()
	st.Assert(t, err, nil)

	st.Assert(t, mdata.Name, "test_db")
}

func TestConnectionPing(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Head("/").
		Reply(200)

	con := globalTestConnections.Version1(t, true)

	if err := con.Ping(); err != nil {
		t.Fatal(err)
	}
}

func TestConnectionServerInfo1(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Get("/").
		Reply(200).
		JSON(map[string]interface{}{
			"couchdb": "Welcome",
			"uuid":    "038bbc5ae438344067d6ab90fe30ed05",
			"version": "1.6.1",
			"vendor":  map[string]string{"name": "The Apache Software Foundation", "version": "1.6.1"},
		})

	con := globalTestConnections.Version1(t, true)

	info, err := con.ServerInfo()
	st.Assert(t, err, nil)

	st.Assert(t, info.CouchDB, "Welcome")
	st.Assert(t, info.UUID, "038bbc5ae438344067d6ab90fe30ed05")
	st.Assert(t, info.Version, "1.6.1")
	st.Assert(t, info.Vendor["name"].(string), "The Apache Software Foundation")
	st.Assert(t, info.Vendor["version"].(string), "1.6.1")
}

func TestConnectionServerInfo2(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version2MockHost)).
		Get("/").
		Reply(200).
		JSON(map[string]interface{}{
			"couchdb":  "Welcome",
			"features": []string{"scheduler"},
			"version":  "2.1.2",
			"vendor":   map[string]string{"name": "The Apache Software Foundation"},
		})

	con := globalTestConnections.Version2(t, true)

	info, err := con.ServerInfo()
	st.Assert(t, err, nil)

	st.Assert(t, info.CouchDB, "Welcome")
	st.Assert(t, info.Features, []string{"scheduler"})
	st.Assert(t, info.Version, "2.1.2")
	st.Assert(t, info.Vendor["name"].(string), "The Apache Software Foundation")
}

func TestConnectionServerInfo3(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version3MockHost)).
		Get("/").
		Reply(200).
		JSON(map[string]interface{}{
			"couchdb":  "Welcome",
			"features": []string{"access-ready", "partitioned", "pluggable-storage-engines", "reshard", "scheduler"},
			"version":  "3.2.1",
			"uuid":     "8d84dd13ec80d945086dbf40a068f910",
			"git_sha":  "244d428af",
			"vendor":   map[string]string{"name": "The Apache Software Foundation"},
		})

	con := globalTestConnections.Version3(t, true)

	info, err := con.ServerInfo()
	st.Assert(t, err, nil)

	st.Assert(t, info.CouchDB, "Welcome")
	st.Assert(t, info.Features, []string{"access-ready", "partitioned", "pluggable-storage-engines", "reshard", "scheduler"})
	st.Assert(t, info.Version, "3.2.1")
	st.Assert(t, info.Vendor["name"].(string), "The Apache Software Foundation")
}

func TestConnectionListDatabases(t *testing.T) {
	defer gock.Off()

	testDBs := []string{
		"fruits",
		"testdb",
	}

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Get("/_all_dbs").
		Reply(200).
		JSON(testDBs)

	con := globalTestConnections.Version1(t, true)

	serverDBs, err := con.ListDatabases()
	st.Assert(t, err, nil)

	for i, name := range testDBs {
		st.Assert(t, serverDBs[i], name)
	}
}

func TestConnectionReal1(t *testing.T) {
	con := globalTestConnections.Version1(t, false)

	info, err := con.ServerInfo()
	st.Assert(t, err, nil)

	st.Assert(t, info.CouchDB, "Welcome")

	cleanTestDB(t, con.Connection, "test_db", true)

	db, err := con.CreateDatabase("test_db")
	st.Assert(t, err, nil)

	doc := getDefaultTestDoc()

	_, err = db.Put(doc)
	st.Assert(t, err, nil)

	_, err = db.Get(doc, doc.Metadata().ID, "")
	st.Assert(t, err, nil)

	_, err = db.Delete(doc)
	st.Assert(t, err, nil)

	cleanTestDB(t, con.Connection, "test_db", false)
}

func TestConnectionReal2(t *testing.T) {
	con := globalTestConnections.Version2(t, false)

	info, err := con.ServerInfo()
	st.Assert(t, err, nil)

	st.Assert(t, info.CouchDB, "Welcome")

	cleanTestDB(t, con.Connection, "test_db", true)

	db, err := con.CreateDatabase("test_db")
	st.Assert(t, err, nil)

	doc := getDefaultTestDoc()

	_, err = db.Put(doc)
	st.Assert(t, err, nil)

	_, err = db.Get(doc, doc.Metadata().ID, "")
	st.Assert(t, err, nil)

	_, err = db.Delete(doc)
	st.Assert(t, err, nil)

	cleanTestDB(t, con.Connection, "test_db", false)
}

func TestConnectionReal3(t *testing.T) {
	con := globalTestConnections.Version3(t, false)

	info, err := con.ServerInfo()
	st.Assert(t, err, nil)

	st.Assert(t, info.CouchDB, "Welcome")

	cleanTestDB(t, con.CouchDB2Connection.Connection, "test_db", true)

	db, err := con.CreateDatabase("test_db")
	st.Assert(t, err, nil)

	doc := getDefaultTestDoc()

	_, err = db.Put(doc)
	st.Assert(t, err, nil)

	_, err = db.Get(doc, doc.Metadata().ID, "")
	st.Assert(t, err, nil)

	_, err = db.Delete(doc)
	st.Assert(t, err, nil)

	cleanTestDB(t, con.CouchDB2Connection.Connection, "test_db", false)
}
