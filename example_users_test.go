package sofa_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/joelnb/sofa"
)

func Example_CreateUser() {
	conn, err := sofa.NewConnection("http://localhost:5984", 10*time.Second, sofa.NullAuthenticator())
	if err != nil {
		panic(err)
	}

	db := conn.Database("_users")

	user := sofa.UserDocument{
		DocumentMetadata: sofa.DocumentMetadata{
			ID: "org.couchdb.user:Couchy McCouchface",
		},
		Name:     "Couchy McCouchface",
		Password: "example",
		Roles:    []string{"boat"},
		TheType:  "user",
	}

	rev, err := db.Put(&user)
	if err != nil {
		panic(err)
	}

	// Show the revision which has been created
	fmt.Println(strings.Split(rev, "-")[0])

	// Retrieve the document which was created
	realUser, err := conn.User("org.couchdb.user:Couchy McCouchface", rev)
	if err != nil {
		panic(err)
	}

	delrev, err := db.Delete(&realUser)
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Split(delrev, "-")[0])

	// Output:
	// 1
	// 2
}
