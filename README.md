# sofa

Simply connect to CouchDB database servers using Go.

## Examples

View the [full documentation](https://godoc.org/github.com/joelnb/sofa).

Create a new database:

```go
package main

import (
    "time"
    "fmt"

    "github.com/joelnb/sofa"
)

func main() {
    conn, err := sofa.NewConnection("http://localhost:5984", 10*time.Second, sofa.NullAuthenticator())
    if err != nil {
        panic(err)
    }

    db, err := conn.CreateDatabase("example_db")
    if err != nil {
        panic(err)
    }

    fmt.Println(db)
}
```

Add a document to an existing database (and the retrieve it again):

```go
package main

import (
    "time"
    "fmt"

    "github.com/joelnb/sofa"
)

func main() {
    conn, err := sofa.NewConnection("http://localhost:5984", 10*time.Second, sofa.NullAuthenticator())
    if err != nil {
        panic(err)
    }

    db := conn.Database("example_db")

    doc := &struct {
        sofa.DocumentMetadata
        Name string `json:"name"`
        Type string `json:"type"`
    }{
        DocumentMetadata: sofa.DocumentMetadata{
            ID: "fruit1",
        },
        Name: "apple",
        Type: "fruit",
    }

    rev, err := db.Put(doc)
    if err != nil {
        panic(err)
    }

    fmt.Println(rev)

    getDoc := &struct {
        sofa.DocumentMetadata
        Name string `json:"name"`
        Type string `json:"type"`
    }{}

    getRev, err := db.Get(getDoc, "fruit1", "")
    if err != nil {
        panic(err)
    }

    fmt.Println(getRev)
    fmt.Println(getDoc.Metadata().Rev)
}
```

## Limitations

Large parts of the [CouchDB API](http://docs.couchdb.org/en/2.0.0/api/) are covered but not all functionality is currently implemented. Pull requests for any missing functionality would be welcomed!

## Contributing

Contributions of all sizes are welcomed. Simply make a pull request and I will be happy to discuss. If you don't have time to write the code please consider at least creating an issue so that I can ensure the issue gets sorted eventually.

### Running tests

The basic tests can be run using a simple `go test`. To run a more complete set of tests which access a real database you will need a temporary CouchDB instance. The simplest way to create this is using docker:

    docker run -d --name couchdb -p 5984:5984 couchdb:1

You can then set `SOFA_TEST_HOST` appropriately to use the server:

    SOFA_TEST_HOST=http://127.0.0.1:5984 go test -v
