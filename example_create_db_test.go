package sofa_test

import (
	"fmt"
	"time"

	"github.com/joelnb/sofa"
)

// ExampleCreateDatabase
func Example_CreateDatabase() {
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

	// Show the DB details
	fmt.Println(db.Name())

	md, err := db.Metadata()
	if err != nil {
		panic(err)
	}

	fmt.Println(md.DocCount)
	fmt.Println(md.DelCount)

	// Output:
	// example_db
	// 0
	// 0
}
