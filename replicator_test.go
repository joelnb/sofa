package sofa

import (
	"fmt"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
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

	repl := con.NewReplication("_myrepl", "mydb", "yourdb")
	repl.CreateTarget = true
	repl.Continuous = true

	rev, err := con.Database("_replicator").Put(repl)
	st.Assert(t, err, nil)

	st.Assert(t, rev, "1-801609c9fdb4c6d196820c5b1f3c26c9")
}

func TestReplicatorReal(t *testing.T) {
	con := globalTestConnections.Version1(t, false)

	getRepl, err := con.Replication("my_repl", "")
	if err != nil {
		if !ErrorStatus(err, 404) {
			t.Fatal(err)
		}
	} else {
		err := con.DeleteReplication(getRepl)
		if err != nil {
			t.Fatal(err)
		}
	}

	putRepl := con.NewReplication("my_repl", "test_db", "yourdb")
	putRepl.CreateTarget = true
	putRepl.Continuous = true

	_, err = con.Database("_replicator").Put(putRepl)
	st.Assert(t, err, nil)

	putRepl, err = con.Replication("my_repl", "")
	st.Assert(t, err, nil)

	err = con.DeleteReplication(putRepl)
	st.Assert(t, err, nil)
}
