package sofa

import (
	"fmt"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

func TestDocumentList(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", globalTestConnections.Version1MockHost)).
		Post("/test_db/_all_docs").
		MatchParam("include_docs", "^true$").
		Reply(200).
		JSON(map[string]interface{}{
			"total_rows": 3,
			"offset":     0,
			"rows": []map[string]interface{}{
				{"id": "joel", "key": "joel", "value": map[string]string{
					"rev": "1-f75c2efff6baeed9e202d0cf58b8bd60",
				}},
				{"id": "bob", "key": "bob", "value": map[string]string{
					"rev": "1-2de76086df68b0d3840f20dff114a000",
				}},
				{"id": "jane", "key": "jane", "value": map[string]string{
					"rev": "10-6bf653103ec443a04119fde156189078",
				}},
			},
		})

	con := globalTestConnections.Version1(t, true)
	db := con.Database("test_db")

	allDocs, err := db.Documents()
	if err != nil {
		t.Fatal(err)
	}

	st.Assert(t, allDocs.TotalRows, float64(3))
	st.Assert(t, allDocs.Offset, float64(0))

	st.Assert(t, allDocs.Rows[0].ID, "joel")
	st.Assert(t, allDocs.Rows[1].ID, "bob")
	st.Assert(t, allDocs.Rows[2].ID, "jane")

	st.Assert(t, allDocs.Rows[0].Key, "joel")
	st.Assert(t, allDocs.Rows[1].Key, "bob")
	st.Assert(t, allDocs.Rows[2].Key, "jane")
}
