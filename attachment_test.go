package sofa

import (
	"strings"
	"testing"

	"github.com/nbio/st"
)

var (
	AttachmentTestDB   = "attachment_test_db"
	AttachmentFileName = "test.txt"
	AttachmentContent  = "Hello Couchy McCouchface"
)

func TestAttachmentReal(t *testing.T) {
	con := globalTestConnections.Version2(t, false)

	// Cleanup the DB if it existed already
	if err := con.DeleteDatabase(AttachmentTestDB); err != nil && !ErrorStatus(err, 404) {
		t.Fatal(err)
	}

	// Create a new database
	db, err := con.CreateDatabase(AttachmentTestDB)
	st.Assert(t, err, nil)

	doc := &struct {
		DocumentMetadata
		Name string
		Type string
	}{
		DocumentMetadata: DocumentMetadata{
			ID: "fruit1",
		},
		Name: "apple",
		Type: "fruit",
	}

	rev, err := db.Put(doc)
	st.Assert(t, err, nil)

	attRev, err := db.PutAttachment(doc.DocumentMetadata.ID, AttachmentFileName, strings.NewReader(AttachmentContent), rev)
	st.Assert(t, err, nil)

	attBytes, err := db.GetAttachment(doc.DocumentMetadata.ID, AttachmentFileName, "")
	st.Assert(t, err, nil)

	attString := string(attBytes)

	if attString != AttachmentContent {
		t.Fatalf("attachment response did not contain expected content: %s", attString)
	}

	_, err = db.DeleteAttachment(doc.DocumentMetadata.ID, AttachmentFileName, attRev)
	st.Assert(t, err, nil)

	_, err = db.GetAttachment(doc.DocumentMetadata.ID, AttachmentFileName, "")
	if !ErrorStatus(err, 404) {
		t.Fatalf("expected a 404 error getting attachment after deletion but got: %s", err)
	}

	err = con.DeleteDatabase(AttachmentTestDB)
	st.Assert(t, err, nil)
}
