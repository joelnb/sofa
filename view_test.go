package sofa

import (
	"fmt"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
)

const TestViewFunc = `function(doc) {
    if (doc.type !== undefined && doc.name !== undefined) {
        emit(doc.type, doc.name)
    }
}`

func TestTemporaryView(t *testing.T) {
	defer gock.Off()

	// TODO: Check the data which gets sent
	gock.New(fmt.Sprintf("https://%s", TestMockHost)).
		Post("/view_test_db/_temp_view").
		BodyString(`{"map":"function(doc) {\n    if (doc.type !== undefined \u0026\u0026 doc.name !== undefined) {\n        emit(doc.type, doc.name)\n    }\n}"}`).
		Reply(200).
		JSON(map[string]interface{}{
			"total_rows": 2,
			"offset":     0,
			"rows": []map[string]interface{}{
				map[string]interface{}{
					"id":    "fruit1",
					"key":   "fruit",
					"value": "apple",
				},
				map[string]interface{}{
					"id":    "fruit2",
					"key":   "fruit",
					"value": "apple",
				},
			},
		})

	con := defaultMockTestConnection(t)
	db := con.Database("view_test_db")

	view := db.TemporaryView(TestViewFunc)
	result, err := view.Execute(NewURLOptions())
	st.Assert(t, err, nil)

	st.Assert(t, result.TotalRows, float64(2))
	st.Assert(t, result.Offset, float64(0))
}

func TestNamedView(t *testing.T) {
	defer gock.Off()

	gock.New(fmt.Sprintf("https://%s", TestMockHost)).
		Get("/view_test_db/_design/things/_view/byType").
		Reply(200).
		JSON(map[string]interface{}{
			"total_rows": 2,
			"offset":     0,
			"rows": []map[string]interface{}{
				map[string]interface{}{
					"id":    "fruit1",
					"key":   "fruit",
					"value": "apple",
				},
				map[string]interface{}{
					"id":    "fruit2",
					"key":   "fruit",
					"value": "apple",
				},
			},
		})

	con := defaultMockTestConnection(t)
	db := con.Database("view_test_db")
	view := db.NamedView("things", "byType")

	result, err := view.Execute(NewURLOptions())
	st.Assert(t, err, nil)

	st.Assert(t, result.TotalRows, float64(2))
	st.Assert(t, result.Offset, float64(0))
}

// TODO: Currently only tests a temporary view. Would be nice to also create a
//       named view and test that too.
func TestViewReal(t *testing.T) {
	serverRequired(t)
	con := defaultTestConnection(t)

	db := initTestDB(t, con, "view_test_db")
	defer cleanupTestDB(t, con, "view_test_db")

	appleDoc := &struct {
		DocumentMetadata
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		DocumentMetadata: DocumentMetadata{
			ID: "fruit1",
		},
		Name: "apple",
		Type: "fruit",
	}

	kiwiDoc := &struct {
		DocumentMetadata
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		DocumentMetadata: DocumentMetadata{
			ID: "fruit2",
		},
		Name: "kiwi",
		Type: "fruit",
	}

	var err error
	_, err = db.Put(appleDoc)
	st.Assert(t, err, nil)

	_, err = db.Put(kiwiDoc)
	st.Assert(t, err, nil)

	view := db.TemporaryView(TestViewFunc)
	result, err := view.Execute(NewURLOptions())
	st.Assert(t, err, nil)

	st.Assert(t, result.TotalRows, float64(2))
	st.Assert(t, result.Offset, float64(0))
}
