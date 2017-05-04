package sofa

import (
	"os"
	"testing"
	"time"

	"github.com/nbio/st"
)

var (
	TestServerAvailable = false
	TestServer          = "127.0.0.1:5984"
	TestMockHost        = "couchdb.local"
)

func TestMain(m *testing.M) {
	// Setup
	if s := os.Getenv("SOFA_TEST_HOST"); s != "" {
		TestServerAvailable = true
		TestServer = s
	}
	// Run
	retCode := m.Run()
	// Teardown
	os.Exit(retCode)
}

func initTestDB(t *testing.T, con *Connection, dbname string) *Database {
	// Delete the DB if it currently exists
	if err := con.DeleteDatabase(dbname); err != nil {
		if !ErrorStatus(err, 404) {
			t.Fatal(err)
		}
	}

	db, err := con.CreateDatabase(dbname)
	st.Assert(t, err, nil)
	return db
}

func cleanupTestDB(t *testing.T, con *Connection, dbname string) {
	if err := con.DeleteDatabase(dbname); err != nil {
		t.Fatal(err)
	}
}

// serverRequired skips the current test if the address of a CouchDB server has not
// been provided by the user. This is to prevent all of these tests failing if the
// server is not available.
func serverRequired(t *testing.T) {
	if !TestServerAvailable {
		t.Skip("skipping - $SOFA_TEST_HOST not set")
	}
}

func defaultTestConnection(t *testing.T) *Connection {
	conn, err := NewConnection(TestServer, 10*time.Second, NullAuthenticator())
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	return conn
}

func defaultMockTestConnection(t *testing.T) *Connection {
	conn, err := NewConnection(TestMockHost, 10*time.Second, NullAuthenticator())
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	return conn
}
