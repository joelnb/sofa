# sofa

Simply connect to CouchDB database servers using Go.

## Examples

Create a new database:

    conn, err := sofa.NewConnection("http://localhost:5984", 10*time.Second, NullAuthenticator())
    if err != nil {
        panic(err)
    }

    db, err := conn.CreateDatabase("example_db")
    if err != nil {
        panic(err)
    }

Add a document to an existing database:

    conn, err := sofa.NewConnection("http://localhost:5984", 10*time.Second, NullAuthenticator())
    if err != nil {
        panic(err)
    }

    db := conn.Database("example_db")

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

    rev, err := db.Put(doc)
    if err != nil {
        panic(err)
    }

## Limitations

Large parts of the [CouchDB API](http://docs.couchdb.org/en/2.0.0/api/) are covered but not all functionality is currently implemented. Pull requests for any missing functionality would be welcomed!

## Contributing

Contributions of all sizes are welcomed. Simply make a pull request and I will be happy to discuss. If you don't have time to write the code please consider at least creating an issue so that I can ensure the issue gets sorted eventually.

### Running tests

The basic tests can be run using a simple `go test`. To run a more complete set of tests which access a real database you will need a temporary CouchDB instance. The simplest way to create this is using docker:

    docker run -d --name couchdb -p 5984:5984 couchdb

You can then set `SOFA_TEST_HOST` appropriately to use the server:

    SOFA_TEST_HOST=http://127.0.0.1:5984 go test -v
