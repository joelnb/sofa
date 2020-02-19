package sofa_test

import (
	"fmt"
	"time"

	"github.com/joelnb/sofa"
)

// ExampleDeleteDatabase
func Example_deleteDatabase() {
	// Make the connection with no authentication
	conn, err := sofa.NewConnection("http://localhost:5984", 10*time.Second, sofa.NullAuthenticator())
	if err != nil {
		panic(err)
	}

	// Cleanup the database again
	if err := conn.DeleteDatabase("example_db"); err != nil {
		panic(err)
	}

	// Get the current list of databases and show that example_db has been removed
	dbs, err := conn.ListDatabases()
	if err != nil {
		panic(err)
	}
	fmt.Println(dbs)

	// Output:
	// []
}
