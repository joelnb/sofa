package sofa

import (
	"strings"
	"testing"
)

var (
	AttachmentTestDB   = "attachment_test_db"
	AttachmentFileName = "test.txt"
	AttachmentContent  = "Hello Couchy McCouchface"
)

func TestAttachmentReal(t *testing.T) {
	con := globalTestConnections.Version2(t, false)

	// Create a new database
	db, err := con.CreateDatabase(AttachmentTestDB)
	if err != nil {
		t.Fatal(err)
	}

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
	if err != nil {
		t.Fatal(err)
	}

	attRev, err := db.PutAttachment(doc.DocumentMetadata.ID, AttachmentFileName, strings.NewReader(AttachmentContent), rev)
	if err != nil {
		t.Fatal(err)
	}

	attBytes, err := db.GetAttachment(doc.DocumentMetadata.ID, AttachmentFileName, "")
	if err != nil {
		t.Fatal(err)
	}

	attString := string(attBytes)

	if attString != AttachmentContent {
		t.Fatalf("attachment response did not contain expected content: %s", attString)
	}

	_, err = db.DeleteAttachment(doc.DocumentMetadata.ID, AttachmentFileName, attRev)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GetAttachment(doc.DocumentMetadata.ID, AttachmentFileName, "")
	if !ErrorStatus(err, 404) {
		t.Fatalf("expected a 404 error getting attachment after deletion but got: %s", err)
	}

	if err := con.DeleteDatabase(AttachmentTestDB); err != nil {
		t.Fatal(err)
	}
}
