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
    // Make the connection with no authentication
    conn, err := sofa.NewConnection("http://localhost:5984", 10*time.Second, sofa.NullAuthenticator())
    if err != nil {
        panic(err)
    }

    // Create a new database
    db, err := conn.CreateDatabase("example_db")
    if err != nil {
        panic(err)
    }

    // Show the format of the Database struct
    fmt.Printf("%+v\n", db)
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

type Fruit struct {
    sofa.DocumentMetadata
    Name string `json:"name"`
    Type string `json:"type"`
}

func main() {
    conn, err := sofa.NewConnection("http://localhost:5984", 10*time.Second, sofa.NullAuthenticator())
    if err != nil {
        panic(err)
    }

    // Open the previously-created database
    db := conn.Database("example_db")

    // Create a document to save
    doc := &Fruit{
        DocumentMetadata: sofa.DocumentMetadata{
            ID: "fruit1",
        },
        Name: "apple",
        Type: "fruit",
    }

    // Create the document on the CouchDB server
    rev, err := db.Put(doc)
    if err != nil {
        panic(err)
    }

    // Show the revision which has been created
    fmt.Println(rev)

    // Retrieve the document which was created
    getDoc := &Fruit{}
    getRev, err := db.Get(getDoc, "fruit1", "")
    if err != nil {
        panic(err)
    }

    // Is the document we put there the same revision we got back?
    if getRev != rev {
        panic("someone changed the document while this was running")
    }

    // Show the struct which was returned from the DB
    fmt.Printf("%+v\n", getDoc)
}
```

## Limitations

Large parts of the [CouchDB API](http://docs.couchdb.org/en/2.0.0/api/) are covered but not all functionality is currently implemented. Pull requests for any missing functionality would be welcomed!

## Contributing

Contributions of all sizes are welcomed. Simply make a pull request and I will be happy to discuss. If you don't have time to write the code please consider at least creating an issue so that I can ensure the issue gets sorted eventually.

### Running tests

The basic tests can be run using a simple `go test`. To run a more complete set of tests which access a real database you will need a temporary CouchDB instance. The simplest way to create this is using docker:

    docker run -d --name couchdb -p 5984:5984 couchdb:1

To run all the tests you will also need a version 2 server:

    docker run -d --name couchdb -p 5985:5984 couchdb:2

You can then set `SOFA_TEST_HOST_1` and `SOFA_TEST_HOST_2` to set the connection to each server:

    SOFA_TEST_HOST_1=http://127.0.0.1:5984 SOFA_TEST_HOST_2=http://127.0.0.1:5985 go test -v

If you have chosen to only start a single version then only include the appropriate environment variable to ensure tests for the other version are not run.
