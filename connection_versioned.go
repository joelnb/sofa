package sofa

import (
    "time"
)

// CouchDB1Connection is a connection specifically for a version 1 server.
type CouchDB1Connection struct {
    *Connection
}

// CouchDB2Connection is a connection specifically for a version 2 server.
type CouchDB2Connection struct {
    *Connection
}

// CouchDB3Connection is a connection specifically for a version 3 server.
type CouchDB3Connection struct {
    *CouchDB2Connection
}

// NewConnection creates a new CouchDB1Connection which can be used to interact with a single CouchDB server.
// Any query parameters passed in the serverUrl are discarded before creating the connection.
func NewConnection(serverURL string, timeout time.Duration, auth Authenticator) (*CouchDB1Connection, error) {
    con, err := newConnection(serverURL, timeout, auth)
    if err != nil {
        return nil, err
    }

    return &CouchDB1Connection{con}, nil
}

// NewConnection2 creates a new CouchDB2Connection which can be used to interact with a single CouchDB server.
// Any query parameters passed in the serverUrl are discarded before creating the connection.
func NewConnection2(serverURL string, timeout time.Duration, auth Authenticator) (*CouchDB2Connection, error) {
    con, err := newConnection(serverURL, timeout, auth)
    if err != nil {
        return nil, err
    }

    return &CouchDB2Connection{con}, nil
}

// NewConnection3 creates a new CouchDB3Connection which can be used to interact with a single CouchDB server.
// Any query parameters passed in the serverUrl are discarded before creating the connection.
func NewConnection3(serverURL string, timeout time.Duration, auth Authenticator) (*CouchDB3Connection, error) {
    con, err := NewConnection2(serverURL, timeout, auth)
    if err != nil {
        return nil, err
    }

    return &CouchDB3Connection{con}, nil
}

// ServerInfo gets the information about this CouchDB instance returned when accessing the root
// page
func (con *CouchDB1Connection) ServerInfo() (ServerDetails1, error) {
    d := ServerDetails1{}
    _, err := con.unmarshalRequest("GET", "/", NewURLOptions(), nil, &d)
    return d, err
}

// ServerInfo gets the information about this CouchDB instance returned when accessing the root
// page
func (con *CouchDB2Connection) ServerInfo() (ServerDetails2, error) {
    d := ServerDetails2{}
    _, err := con.unmarshalRequest("GET", "/", NewURLOptions(), nil, &d)
    return d, err
}
