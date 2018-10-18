package sofa

import (
	"testing"

	"github.com/nbio/st"
)

func TestConnectionAllStatsRealVersion1(t *testing.T) {
	con := globalTestConnections.Version1(t, false)

	stats, err := con.AllStatistics()
	st.Assert(t, err, nil)

	st.Refute(t, stats.HTTPD, nil)
	st.Refute(t, stats.HTTPD.Requests, nil)
}

func TestConnectionStatRealVersion1(t *testing.T) {
	con := globalTestConnections.Version1(t, false)

	stats, err := con.Statistic("httpd", "requests")
	st.Assert(t, err, nil)

	st.Refute(t, stats.HTTPD, nil)
	st.Refute(t, stats.HTTPD.Requests, nil)
}

func TestConnectionAllStatsRealVersion2(t *testing.T) {
	con := globalTestConnections.Version2(t, false)

	stats, err := con.AllStatistics()
	st.Assert(t, err, nil)

	st.Refute(t, stats.CouchDB.HTTPD, nil)
	st.Refute(t, stats.CouchDB.HTTPD.Requests, nil)
}

func TestConnectionStatRealVersion2(t *testing.T) {
	con := globalTestConnections.Version2(t, false)

	stats, err := con.Statistic("couchdb", "httpd/requests")
	st.Assert(t, err, nil)

	st.Refute(t, stats.CouchDB.HTTPD, nil)
	st.Refute(t, stats.CouchDB.HTTPD.Requests, nil)
}
