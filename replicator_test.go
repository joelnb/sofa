package sofa

import (
	"fmt"
	"testing"
	"time"

	"github.com/nbio/st"
	"github.com/h2non/gock"
)

var (
	ReplicatorTestDB       = "replicator_test"
	ReplicatorTestTargetDB = "replicator_test_target"
	ReplicatorTestName     = "my_repl"
)

func TestReplicatorSave(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Put("/_replicator/_myrepl").
		Reply(201).
		JSON(map[string]interface{}{
			"ok":  true,
			"id":  "_myrepl",
			"rev": "1-801609c9fdb4c6d196820c5b1f3c26c9",
		}).
		SetHeader("Etag", `"1-801609c9fdb4c6d196820c5b1f3c26c9"`)

	con := globalTestConnections.Version1(t, true)

	repl := NewReplication("_myrepl", "mydb", "yourdb")
	repl.CreateTarget = true
	repl.Continuous = true

	rev, err := con.PutReplication(&repl)
	st.Assert(t, err, nil)

	st.Assert(t, rev, "1-801609c9fdb4c6d196820c5b1f3c26c9")
}

func TestReplicatorReal(t *testing.T) {
	con := globalTestConnections.Version1(t, false)

	// Run the cleanup as a best effort task - failures here do not case the test to fail
	defer func() {
		if err := con.DeleteDatabase(ReplicatorTestDB); err != nil {
			t.Logf("error cleaning up source DB: %s", err)
		}

		if err := con.DeleteDatabase(ReplicatorTestTargetDB); err != nil {
			t.Logf("error cleaning up target DB: %s", err)
		}
	}()

	// Create a new database
	_, err := con.CreateDatabase(ReplicatorTestDB)
	if err != nil {
		t.Fatal(err)
	}

	getRepl, err := con.Replication(ReplicatorTestName, "")
	if err != nil {
		if !ErrorStatus(err, 404) {
			t.Fatal(err)
		}
	} else {
		_, err := con.DeleteReplication(&getRepl)
		if err != nil {
			t.Fatal(err)
		}
	}

	putRepl := NewReplication(ReplicatorTestName, ReplicatorTestDB, ReplicatorTestTargetDB)
	putRepl.CreateTarget = true
	putRepl.Continuous = true
	putRepl.Owner = "admin"
	putRepl.UserContext = map[string]interface{}{"roles": []string{"_admin"}}

	_, err = con.PutReplication(&putRepl)
	st.Assert(t, err, nil)

	// This loop handles the case where CouchDB is slow to create the database after the replication
	// is created. This seems to get slower after the first time the replication is created but the
	// race condition is always present.
	getDBAttempts := 1
	for {
		_, err = con.EnsureDatabase(ReplicatorTestTargetDB)
		if err != nil {
			if getDBAttempts >= 10 {
				fmt.Println(time.Now())
				t.Fatalf("couldn't get the created database after %d attempts: %s", getDBAttempts, err)
			}

			time.Sleep(1 * time.Second)
			getDBAttempts++

			continue
		}

		break
	}

	// This loop is to handle the case where we were able to retrieve the DB before CouchDB had updated
	// the replicator document to account for this - in this case the replicator document can have the state
	// set to error with the reason: "db_not_found: could not open replicator_test_target"
	getReplAttempts := 1
	for {
		getRepl, err = con.Replication(ReplicatorTestName, "")
		st.Assert(t, err, nil)

		if getRepl.ReplicationState != "triggered" {
			if getReplAttempts >= 3 {
				fmt.Println(time.Now())
				t.Fatalf(
					"replication still in state '%s' after %d attempts: %s",
					getRepl.ReplicationState,
					getReplAttempts,
					getRepl.ReplicationStateReason,
				)
			}

			time.Sleep(1 * time.Second)
			getReplAttempts++

			continue
		}

		break
	}

	_, err = con.DeleteReplication(&getRepl)
	st.Assert(t, err, nil)
}
