package sofa

import (
	"testing"
	"time"

	"github.com/nbio/st"
)

func TestFutonURL(t *testing.T) {
	conn := globalTestConnections.Version1(t, true)
	baseURL := conn.URL("/")

	uA := conn.FutonURL("/database_name")
	resA := urlConcat(baseURL.String(), "_utils/database.html?database_name")
	st.Expect(t, resA, uA.String())

	uB := conn.FutonURL("/database_name/document_name")
	resB := urlConcat(baseURL.String(), "_utils/document.html?database_name/document_name")
	st.Expect(t, resB, uB.String())

	uC := conn.FutonURL("/database_name/document_name/")
	resC := urlConcat(baseURL.String(), "_utils/document.html?database_name/document_name/")
	st.Expect(t, resC, uC.String())

	uD := conn.FutonURL("/database_name/document_name/attachment_name")
	resD := urlConcat(baseURL.String(), "_utils/document.html?database_name/document_name/attachment_name")
	st.Expect(t, resD, uD.String())

	uE := conn.FutonURL("/database_name/_all_docs")
	resE := urlConcat(baseURL.String(), "_utils/database.html?database_name/_all_docs")
	st.Expect(t, resE, uE.String())

	uF := conn.FutonURL("/database_name/_design/myDesign/_theDocument")
	resF := urlConcat(baseURL.String(), "_utils/database.html?database_name/_design/myDesign/_theDocument")
	st.Expect(t, resF, uF.String())

	altServerURL := "https://the-couchdb-server:6984/couchdb"
	conn2, err := NewConnection(altServerURL, 10*time.Second, NullAuthenticator())
	if err != nil {
		t.Logf("%v\n", err)
		t.Fail()
	}

	u1 := conn2.FutonURL("/database_name/document_name/")
	res1 := urlConcat(altServerURL, "_utils/document.html?database_name/document_name/")
	st.Expect(t, res1, u1.String())

	u2 := conn2.FutonURL("/database_name/document_name")
	res2 := urlConcat(altServerURL, "_utils/document.html?database_name/document_name")
	st.Expect(t, res2, u2.String())

	u3 := conn2.FutonURL("/database_name/document_name/")
	res3 := urlConcat(altServerURL, "_utils/document.html?database_name/document_name/")
	st.Expect(t, res3, u3.String())

	u4 := conn2.FutonURL("/database_name/document_name/attachment_name")
	res4 := urlConcat(altServerURL, "_utils/document.html?database_name/document_name/attachment_name")
	st.Expect(t, res4, u4.String())

	u5 := conn2.FutonURL("/database_name/_all_docs")
	res5 := urlConcat(altServerURL, "_utils/database.html?database_name/_all_docs")
	st.Expect(t, res5, u5.String())

	u6 := conn2.FutonURL("/database_name/_design/myDesign/_theDocument")
	res6 := urlConcat(altServerURL, "_utils/database.html?database_name/_design/myDesign/_theDocument")
	st.Expect(t, res6, u6.String())
}
