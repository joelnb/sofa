package sofa

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

// Connection is a connection to a CouchDB server.
type Connection struct {
	auth Authenticator
	http *http.Client

	url     *url.URL
	timeout time.Duration
}

// NewConnection creates a new Connection which can be used to interact with a single CouchDB server.
// Any query parameters passed in the serverUrl are discarded before creating the connection.
func NewConnection(serverURL string, timeout time.Duration, auth Authenticator) (*Connection, error) {
	hasURLScheme, err := regexp.MatchString("^https?://.*", serverURL)
	if err != nil {
		return nil, err
	}

	if !hasURLScheme {
		serverURL = fmt.Sprintf("https://%s", serverURL)
	}

	surl, err := url.Parse(serverURL)
	if err != nil {
		return nil, nil
	}

	// Ensure no query parameters are set at this point
	surl.RawQuery = ""

	con := &Connection{
		auth: auth,

		url:     surl,
		timeout: timeout,
	}

	client, err := auth.Client()
	if err != nil {
		return nil, err
	}
	con.http = client

	if err := auth.Setup(con); err != nil {
		return nil, err
	}

	return con, nil
}

// URL returns the URL of the server with a path appended.
func (con *Connection) URL(path string) url.URL {
	durl := *con.url
	durl.Path = urlConcat(durl.Path, path)
	return durl
}

// Database creates a new Database object. No validation or contact with the couchdb
// server is performed in this method so it is possible to create Database objects for
// databases which do not exist
func (con *Connection) Database(name string) *Database {
	return &Database{
		name: name,
		con:  con,
	}
}

// EnsureDatabase creates a new Database object & then requests the metadata to ensure
// that the Database actually exists on the server. The Database is returned with the
// metadata already available.
func (con *Connection) EnsureDatabase(name string) (*Database, error) {
	db := con.Database(name)
	_, err := db.Metadata()
	return db, err
}

// Request performs a request and returns the http.Response which results from that request.
// Request also checks the response status and returns a ResponseError if an error HTTP
// statuscode is received.
func (con *Connection) Request(method, path string, opts Options, body io.Reader) (resp *http.Response, err error) {
	return con.urlRequest(method, con.URL(path), opts, body, true)
}

func (con *Connection) urlRequest(method string, durl url.URL, opts Options, body io.Reader, doTimeout bool) (resp *http.Response, err error) {
	durl.RawQuery = opts.Encode()

	req, err := newRequest(method, durl.String(), body)
	if err != nil {
		return nil, err
	}

	// Let the Authenticator add info to the request.
	con.auth.Authenticate(req)

	// Set timeout on the request if it was needed (so not for long-polling etc.)
	if doTimeout {
		ctx, cancel := context.WithTimeout(req.Context(), con.timeout)
		defer func() {
			if ctx.Err() != context.DeadlineExceeded {
				// TODO: I think this call should be here but it causes errors to bubble up
				// cancel()
				_ = cancel
			}
		}()

		req = req.WithContext(ctx)
	}

	c := make(chan error, 1)
	go func() {
		resp, err = con.http.Do(req)
		c <- err
	}()

	select {
	case err := <-c:
		if err != nil {
			return nil, err
		}
	}

	if resp.StatusCode >= 400 {
		return nil, httpResponseError(resp)
	}

	return resp, nil
}

// unmarshalRequest performs a request and then attempts to unmarshal the result into the
// provided value.
func (con *Connection) unmarshalRequest(method, path string, opts Options, body io.Reader, res interface{}) (*http.Response, error) {
	resp, err := con.Request(method, path, opts, body)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(res); err != nil {
		return nil, err
	}
	return resp, nil
}

// Get sends a GET request to the provided path on the CouchDB server.
func (con *Connection) Get(path string, opts Options) (resp *http.Response, err error) {
	return con.Request("GET", path, opts, nil)
}

// Head sends a HEAD request to the provided path on the CouchDB server.
func (con *Connection) Head(path string, opts Options) (resp *http.Response, err error) {
	return con.Request("HEAD", path, opts, nil)
}

// Put sends a PUT request to the provided path on the CouchDB server. The contents of the provided
// io.Reader is sent as the body of the request.
func (con *Connection) Put(path string, opts Options, body io.Reader) (resp *http.Response, err error) {
	return con.Request("PUT", path, opts, body)
}

// Patch sends a PATCH request to the provided path on the CouchDB server. The contents of the provided
// io.Reader is sent as the body of the request.
func (con *Connection) Patch(path string, opts Options, body io.Reader) (resp *http.Response, err error) {
	return con.Request("PATCH", path, opts, body)
}

// Post sends a POST request to the provided path on the CouchDB server. The contents of the provided
// io.Reader is sent as the body of the request.
func (con *Connection) Post(path string, opts Options, body io.Reader) (resp *http.Response, err error) {
	return con.Request("POST", path, opts, body)
}

// Delete sends a DELETE request to the provided path on the CouchDB server.
func (con *Connection) Delete(path string, opts Options) (resp *http.Response, err error) {
	return con.Request("DELETE", path, opts, nil)
}

// Ping tests basic connection to CouchDB by making a HEAD request for one of the databases
func (con *Connection) Ping() error {
	_, err := con.Head("/", NewURLOptions())
	return err
}

// ServerInfo gets the information about this CouchDB instance returned when accessing the root
// page
func (con *Connection) ServerInfo() (ServerDetails, error) {
	d := ServerDetails{}
	_, err := con.unmarshalRequest("GET", "/", NewURLOptions(), nil, &d)
	return d, err
}

// ListDatabases returns the list of all database names on the server. Internal
// couchdb databases _replicator and _users are excluded as they are always
// present & accessed using special methods
func (con *Connection) ListDatabases() (databases []string, err error) {
	if _, err = con.unmarshalRequest("GET", "/_all_dbs", NewURLOptions(), nil, &databases); err != nil {
		return nil, err
	}

	w := 0
	internalNames := []string{"_replicator", "_users"}

loop:
	for _, dbname := range databases {
		for i, iname := range internalNames {
			if iname == dbname {
				internalNames = append(internalNames[:i], internalNames[i+1:]...)
				continue loop
			}
		}
		databases[w] = dbname
		w++
	}

	return databases[:w], nil
}

// Databases returns a Database object for every database on the server,
// excluding CouchDB internal databases as there are special methods for
// accessing them.
func (con *Connection) Databases() (databases []*Database, err error) {
	dbnames, err := con.ListDatabases()
	if err != nil {
		return nil, err
	}

	for _, dbname := range dbnames {
		databases = append(databases, con.Database(dbname))
	}

	return databases, nil
}

func (con *Connection) CreateDatabase(name string) (*Database, error) {
	resp, err := con.Put(name, NewURLOptions(), nil)
	if err != nil {
		return nil, err
	}

	res := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return con.Database(name), nil
}

func (con *Connection) DeleteDatabase(name string) error {
	resp, err := con.Delete(name, NewURLOptions())
	if err != nil {
		return err
	}

	res := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	return nil
}
