package sofa

import (
	"encoding/json"
)

// UserDocument contains all of the fields used by CouchDB to represent a
// user on the server.
type UserDocument struct {
	DocumentMetadata

	Name     string   `json:"name,omitempty"`
	Password string   `json:"password,omitempty"`
	Roles    []string `json:"roles,omitempty"`
	TheType  string   `json:"type,omitempty"`
}

// Users gets a list of all users currently active on this CouchDB server
func (con *Connection) Users() ([]UserDocument, error) {
	db := con.Database("_users")

	table, err := db.AllDocuments()
	if err != nil {
		return nil, err
	}

	users := []UserDocument{}
	for _, rawUser := range table.RawDocuments() {
		user := UserDocument{}
		if err := json.Unmarshal(rawUser, &user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// User gets a particular UserDocument from the server by ID. A revision can
// also be specified to retrieve a particular revision of the document.
func (con *Connection) User(id string, rev string) (UserDocument, error) {
	db := con.Database("_users")

	user := UserDocument{}
	rev, err := db.Get(&user, id, rev)
	if err != nil {
		return user, err
	}

	return user, nil
}
