package sofa

import (
	"testing"

	"github.com/nbio/st"
)

func TestConnectionAllStatsReal(t *testing.T) {
	serverRequired(t)
	con := defaultTestConnection(t)

	stats, err := con.AllStatistics()
	st.Assert(t, err, nil)

	st.Refute(t, stats.HTTPD, nil)
	st.Refute(t, stats.HTTPD.Requests, nil)
}

func TestConnectionStatReal(t *testing.T) {
	serverRequired(t)
	con := defaultTestConnection(t)

	stats, err := con.Statistic("httpd", "requests")
	st.Assert(t, err, nil)

	st.Refute(t, stats.HTTPD, nil)
	st.Refute(t, stats.HTTPD.Requests, nil)
}
