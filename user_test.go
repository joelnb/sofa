package sofa

import (
	"fmt"
	"testing"

	"github.com/nbio/st"
	"github.com/h2non/gock"
)

func TestCreateUser(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Put("/_users/org.couchdb.user:Couchy McCouchface").
		Reply(201).
		JSON(map[string]interface{}{
			"ok":  true,
			"id":  "org.couchdb.user:Couchy McCouchface",
			"rev": DefaultFirstRev,
		}).
		SetHeader("Etag", fmt.Sprintf(`"%s"`, DefaultFirstRev))

	con := globalTestConnections.Version1(t, true)

	rev, err := con.CreateUser(&UserDocument{
		Name:     "Couchy McCouchface",
		Password: "example",
		Roles:    []string{"boat"},
		TheType:  "user",
	})
	if err != nil {
		t.Fatal(err)
	}

	st.Assert(t, rev, DefaultFirstRev)
}

func TestUpdateUser(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Get("/_users/org.couchdb.user:Couchy McCouchface").
		Reply(200).
		JSON(map[string]interface{}{
			"_id":             "org.couchdb.user:Couchy McCouchface",
			"_rev":            DefaultFirstRev,
			"password_scheme": "pbkdf2",
			"iterations":      10,
			"name":            "Couchy McCouchface",
			"roles":           []string{"boat"},
			"type":            "user",
			"derived_key":     "31aee83dcb20882eadcf455381164dbe605f29c0",
			"salt":            "b6be17ff6bae9d932428478de90eea05",
		}).
		SetHeader("Etag", fmt.Sprintf(`"%s"`, DefaultFirstRev))

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Put("/_users/org.couchdb.user:Couchy McCouchface").
		Reply(201).
		JSON(map[string]interface{}{
			"ok":  true,
			"id":  "org.couchdb.user:Couchy McCouchface",
			"rev": DefaultSecondRev,
		}).
		SetHeader("Etag", fmt.Sprintf(`"%s"`, DefaultSecondRev))

	con := globalTestConnections.Version1(t, true)

	user, err := con.User("Couchy McCouchface", "")
	if err != nil {
		t.Fatal(err)
	}

	user.Roles = append(user.Roles, "awesome")

	rev, err := con.UpdateUser(&user)
	if err != nil {
		t.Fatal(err)
	}

	st.Assert(t, rev, DefaultSecondRev)
}

func TestDeleteUser(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Delete("/_users/org.couchdb.user:Couchy McCouchface").
		Reply(201).
		JSON(map[string]interface{}{
			"ok":  true,
			"id":  "org.couchdb.user:Couchy McCouchface",
			"rev": DefaultThirdRev,
		}).
		SetHeader("Etag", fmt.Sprintf(`"%s"`, DefaultThirdRev))

	con := globalTestConnections.Version1(t, true)

	rev, err := con.DeleteUser(&UserDocument{
		DocumentMetadata: DocumentMetadata{
			Rev: DefaultFirstRev,
			ID:  "org.couchdb.user:Couchy McCouchface",
		},
		Name: "Couchy McCouchface",
	})
	if err != nil {
		t.Fatal(err)
	}

	st.Assert(t, rev, DefaultThirdRev)
}
