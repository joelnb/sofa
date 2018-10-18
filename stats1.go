package sofa

import (
	"fmt"
)

// Statistics1 represents all of the statistics available from the CouchDB
// version 1 server. This is the format returned when all statistics are
// being accessed.
type Statistics1 struct {
	CouchDB struct {
		AuthCacheHits   *Statistic1 `json:"auth_cache_hits,omitempty"`
		AuthCacheMisses *Statistic1 `json:"auth_cache_misses,omitempty"`
		DatabaseReads   *Statistic1 `json:"database_reads,omitempty"`
		DatabaseWrites  *Statistic1 `json:"database_writes,omitempty"`
		OpenDatabases   *Statistic1 `json:"open_databases,omitempty"`
		OpenOsFiles     *Statistic1 `json:"open_os_files,omitempty"`
		RequestTime     *Statistic1 `json:"request_time,omitempty"`
	} `json:"couchdb,omitempty"`
	HTTPD struct {
		BulkRequests             *Statistic1 `json:"bulk_requests,omitempty"`
		ClientsRequestingChanges *Statistic1 `json:"clients_requesting_changes,omitempty"`
		Requests                 *Statistic1 `json:"requests,omitempty"`
		TemporaryViewReads       *Statistic1 `json:"temporary_view_reads,omitempty"`
		ViewReads                *Statistic1 `json:"view_reads,omitempty"`
	} `json:"httpd,omitempty"`
	HTTPDRequestMethods struct {
		Copy   *Statistic1 `json:"COPY,omitempty"`
		Delete *Statistic1 `json:"DELETE,omitempty"`
		Get    *Statistic1 `json:"GET,omitempty"`
		Head   *Statistic1 `json:"HEAD,omitempty"`
		Post   *Statistic1 `json:"POST,omitempty"`
		Put    *Statistic1 `json:"PUT,omitempty"`
	} `json:"httpd_request_methods,omitempty"`
	HTTPDStatusCodes struct {
		Two00   *Statistic1 `json:"200,omitempty"`
		Two01   *Statistic1 `json:"201,omitempty"`
		Two02   *Statistic1 `json:"202,omitempty"`
		Three01 *Statistic1 `json:"301,omitempty"`
		Three04 *Statistic1 `json:"304,omitempty"`
		Four00  *Statistic1 `json:"400,omitempty"`
		Four01  *Statistic1 `json:"401,omitempty"`
		Four03  *Statistic1 `json:"403,omitempty"`
		Four04  *Statistic1 `json:"404,omitempty"`
		Four05  *Statistic1 `json:"405,omitempty"`
		Four09  *Statistic1 `json:"409,omitempty"`
		Four12  *Statistic1 `json:"412,omitempty"`
		Five00  *Statistic1 `json:"500,omitempty"`
	} `json:"httpd_status_codes,omitempty"`
}

// Statistic1 represents the format which each statistic returned from CouchDB version 1 has.
type Statistic1 struct {
	Current     *float64 `json:"current"`
	Description string   `json:"description"`
	Max         *float64 `json:"max"`
	Mean        *float64 `json:"mean"`
	Min         *float64 `json:"min"`
	Stddev      *float64 `json:"stddev"`
	Sum         *float64 `json:"sum"`
}

// AllStatistics gets all of the available statistics from the server.
func (con *CouchDB1Connection) AllStatistics() (Statistics1, error) {
	var stats Statistics1
	_, err := con.unmarshalRequest("GET", "/_stats", NewURLOptions(), nil, &stats)
	if err != nil {
		return Statistics1{}, err
	}
	return stats, nil
}

// Statistic loads a single specific statistic from the server by category & name.
func (con *CouchDB1Connection) Statistic(category, name string) (Statistics1, error) {
	var stats Statistics1
	_, err := con.unmarshalRequest("GET", fmt.Sprintf("/_stats/%s/%s", category, name), NewURLOptions(), nil, &stats)
	if err != nil {
		return Statistics1{}, err
	}
	return stats, nil
}
