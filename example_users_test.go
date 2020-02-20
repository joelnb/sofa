package sofa_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/joelnb/sofa"
)

func Example_createUser() {
	conn, err := sofa.NewConnection("http://localhost:5984", 10*time.Second, sofa.NullAuthenticator())
	if err != nil {
		panic(err)
	}

	user := sofa.UserDocument{
		Name:     "Couchy McCouchface",
		Password: "example",
		Roles:    []string{"boat"},
		TheType:  "user",
	}

	rev, err := conn.CreateUser(&user)
	if err != nil {
		panic(err)
	}

	// Show the revision which has been created
	fmt.Println(strings.Split(rev, "-")[0])

	// Retrieve the document which was created
	realUser, err := conn.User("Couchy McCouchface", rev)
	if err != nil {
		panic(err)
	}

	// Show the whole retrieved user
	// fmt.Printf("%+v\n", realUser)

	// Modify the roles for the user
	realUser.Roles = []string{"issue_creator"}

	// Save the modified user
	updateRev, err := conn.UpdateUser(&realUser)
	if err != nil {
		panic(err)
	}

	// Ensure the document has the latest revision so the delete works
	realUser.DocumentMetadata.Rev = updateRev

	// Delete the user
	delrev, err := conn.DeleteUser(&realUser)
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Split(delrev, "-")[0])
	// Output:
	// 1
	// 3
}
