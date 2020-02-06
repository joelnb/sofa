package sofa

import (
	"os"
	"testing"
	"time"

	"gopkg.in/h2non/gock.v1"
)

var globalTestConnections *TestConnections

type TestConnections struct {
	Version1Host string
	Version2Host string

	Version1MockHost string
	Version2MockHost string
}

func NewTestConnections() *TestConnections {
	return &TestConnections{
		Version1Host: os.Getenv("SOFA_TEST_HOST_1"),
		Version2Host: os.Getenv("SOFA_TEST_HOST_2"),

		Version1MockHost: "couchdb1.local",
		Version2MockHost: "couchdb2.local",
	}
}

func (tc *TestConnections) Version1(t *testing.T, mock bool) *CouchDB1Connection {
	host := ""
	if mock {
		host = tc.Version1MockHost
	} else {
		host = tc.Version1Host
		if host == "" {
			t.Skip("skipping - $SOFA_TEST_HOST_1 not set")
			return nil
		}
	}

	con, err := newConnection(host, 10*time.Second, NullAuthenticator())
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if mock {
		gock.InterceptClient(con.http)
	}

	return &CouchDB1Connection{con}
}

func (tc *TestConnections) Version2(t *testing.T, mock bool) *CouchDB2Connection {
	host := ""
	if mock {
		host = tc.Version2MockHost
	} else {
		host = tc.Version2Host
		if host == "" {
			t.Skip("skipping - $SOFA_TEST_HOST_2 not set")
			return nil
		}
	}

	con, err := newConnection(host, 10*time.Second, NullAuthenticator())
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if mock {
		gock.InterceptClient(con.http)
	}

	return &CouchDB2Connection{con}
}

func TestMain(m *testing.M) {
	// Setup
	globalTestConnections = NewTestConnections()

	// Run
	retCode := m.Run()

	// Teardown
	os.Exit(retCode)
}
