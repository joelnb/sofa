package sofa

import (
	"fmt"
)

// Statistics represents all of the statistics available from the CouchDB server.
// This is the format returned when all statistics are being accessed.
type Statistics struct {
	CouchDB struct {
		AuthCacheHits   *Statistic `json:"auth_cache_hits,omitempty"`
		AuthCacheMisses *Statistic `json:"auth_cache_misses,omitempty"`
		DatabaseReads   *Statistic `json:"database_reads,omitempty"`
		DatabaseWrites  *Statistic `json:"database_writes,omitempty"`
		OpenDatabases   *Statistic `json:"open_databases,omitempty"`
		OpenOsFiles     *Statistic `json:"open_os_files,omitempty"`
		RequestTime     *Statistic `json:"request_time,omitempty"`
	} `json:"couchdb,omitempty"`
	HTTPD struct {
		BulkRequests             *Statistic `json:"bulk_requests,omitempty"`
		ClientsRequestingChanges *Statistic `json:"clients_requesting_changes,omitempty"`
		Requests                 *Statistic `json:"requests,omitempty"`
		TemporaryViewReads       *Statistic `json:"temporary_view_reads,omitempty"`
		ViewReads                *Statistic `json:"view_reads,omitempty"`
	} `json:"httpd,omitempty"`
	HTTPDRequestMethods struct {
		Copy   *Statistic `json:"COPY,omitempty"`
		Delete *Statistic `json:"DELETE,omitempty"`
		Get    *Statistic `json:"GET,omitempty"`
		Head   *Statistic `json:"HEAD,omitempty"`
		Post   *Statistic `json:"POST,omitempty"`
		Put    *Statistic `json:"PUT,omitempty"`
	} `json:"httpd_request_methods,omitempty"`
	HTTPDStatusCodes struct {
		Two00   *Statistic `json:"200,omitempty"`
		Two01   *Statistic `json:"201,omitempty"`
		Two02   *Statistic `json:"202,omitempty"`
		Three01 *Statistic `json:"301,omitempty"`
		Three04 *Statistic `json:"304,omitempty"`
		Four00  *Statistic `json:"400,omitempty"`
		Four01  *Statistic `json:"401,omitempty"`
		Four03  *Statistic `json:"403,omitempty"`
		Four04  *Statistic `json:"404,omitempty"`
		Four05  *Statistic `json:"405,omitempty"`
		Four09  *Statistic `json:"409,omitempty"`
		Four12  *Statistic `json:"412,omitempty"`
		Five00  *Statistic `json:"500,omitempty"`
	} `json:"httpd_status_codes,omitempty"`
}

// Statistic represents the format which each statistic returned from CouchDB has.
type Statistic struct {
	Current     *float64 `json:"current"`
	Description string   `json:"description"`
	Max         *float64 `json:"max"`
	Mean        *float64 `json:"mean"`
	Min         *float64 `json:"min"`
	Stddev      *float64 `json:"stddev"`
	Sum         *float64 `json:"sum"`
}

// AllStatistics gets all of the available statistics from the server.
func (con *Connection) AllStatistics() (Statistics, error) {
	var stats Statistics
	_, err := con.unmarshalRequest("GET", "/_stats", NewURLOptions(), nil, &stats)
	if err != nil {
		return Statistics{}, err
	}
	return stats, nil
}

// Statistic loads a single specific statistic from the server by category & name.
func (con *Connection) Statistic(category, name string) (Statistics, error) {
	var stats Statistics
	_, err := con.unmarshalRequest("GET", fmt.Sprintf("/_stats/%s/%s", category, name), NewURLOptions(), nil, &stats)
	if err != nil {
		return Statistics{}, err
	}
	return stats, nil
}
