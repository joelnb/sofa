package sofa

import (
	"fmt"
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

const (
	DefaultFirstRev  = "1-801609c9fdb4c6d196820c5b1f3c26c9"
	DefaultSecondRev = "2-754abe0a29104a287f0493d8c3009524"
	DefaultThirdRev  = "3-0865d3627f7ded46b07d836cf102c5e8"
)

func NewTestConnections() *TestConnections {
	return &TestConnections{
		Version1Host: os.Getenv("SOFA_TEST_HOST_1"),
		Version2Host: os.Getenv("SOFA_TEST_HOST_2"),

		Version1MockHost: "couchdb1.local",
		Version2MockHost: "couchdb2.local",
	}
}

func (tc *TestConnections) hostOrSkip(t *testing.T, mock bool, version int, mockHost, realHost string) string {
	if mock {
		return mockHost
	}

	if realHost == "" {
		t.Skip(fmt.Sprintf("skipping - $SOFA_TEST_HOST_%d not set", version))
	}

	return realHost
}

func (tc *TestConnections) testInnerConnection(t *testing.T, mock bool, version int, mockHost, realHost string, auth Authenticator) *Connection {
	host := tc.hostOrSkip(t, mock, version, mockHost, realHost)

	con, err := newConnection(host, 10*time.Second, auth)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if mock {
		gock.InterceptClient(con.http)
	}

	return con
}

func (tc *TestConnections) Version1(t *testing.T, mock bool) *CouchDB1Connection {
	con := tc.testInnerConnection(t, mock, 1, tc.Version1MockHost, tc.Version1Host, NullAuthenticator())

	return &CouchDB1Connection{con}
}

func (tc *TestConnections) Version2(t *testing.T, mock bool) *CouchDB2Connection {
	con := tc.testInnerConnection(t, mock, 2, tc.Version2MockHost, tc.Version2Host, NullAuthenticator())

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
