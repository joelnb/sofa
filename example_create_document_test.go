package sofa_test

import (
	"fmt"
	"time"

	"github.com/joelnb/sofa"
)

type Fruit struct {
	sofa.DocumentMetadata
	Name string `json:"name"`
	Type string `json:"type"`
}

// ExampleCreateDocument creates a
func Example_createDocument() {
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
	fmt.Printf("Name: %s\n", getDoc.Name)
	fmt.Printf("Type: %s\n", getDoc.Type)

	// Output:
	// Name: apple
	// Type: fruit
}
