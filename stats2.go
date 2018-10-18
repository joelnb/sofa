package sofa

import (
	"fmt"
)

// Statistics2 represents all of the statistics available from the CouchDB version 2 server.
// This is the format returned when all statistics are being accessed.
type Statistics2 struct {
	CouchLog struct {
		Level struct {
			Alert     SimpleStatistic `json:"alert"`
			Critical  SimpleStatistic `json:"critical"`
			Debug     SimpleStatistic `json:"debug"`
			Emergency SimpleStatistic `json:"emergency"`
			Error     SimpleStatistic `json:"error"`
			Info      SimpleStatistic `json:"info"`
			Notice    SimpleStatistic `json:"notice"`
			Warning   SimpleStatistic `json:"warning"`
		} `json:"level"`
	} `json:"couch_log"`
	CouchReplicator struct {
		ChangesManagerDeaths SimpleStatistic `json:"changes_manager_deaths"`
		ChangesQueueDeaths   SimpleStatistic `json:"changes_queue_deaths"`
		ChangesReadFailures  SimpleStatistic `json:"changes_read_failures"`
		ChangesReaderDeaths  SimpleStatistic `json:"changes_reader_deaths"`
		Checkpoints          struct {
			Failure SimpleStatistic `json:"failure"`
			Success SimpleStatistic `json:"success"`
		} `json:"checkpoints"`
		ClusterIsStable SimpleStatistic `json:"cluster_is_stable"`
		Connection      struct {
			Acquires      SimpleStatistic `json:"acquires"`
			Closes        SimpleStatistic `json:"closes"`
			Creates       SimpleStatistic `json:"creates"`
			OwnerCrashes  SimpleStatistic `json:"owner_crashes"`
			Releases      SimpleStatistic `json:"releases"`
			WorkerCrashes SimpleStatistic `json:"worker_crashes"`
		} `json:"connection"`
		DBScans SimpleStatistic `json:"db_scans"`
		Docs    struct {
			CompletedStateUpdates SimpleStatistic `json:"completed_state_updates"`
			DBChanges             SimpleStatistic `json:"db_changes"`
			DbsCreated            SimpleStatistic `json:"dbs_created"`
			DbsDeleted            SimpleStatistic `json:"dbs_deleted"`
			DbsFound              SimpleStatistic `json:"dbs_found"`
			FailedStateUpdates    SimpleStatistic `json:"failed_state_updates"`
		} `json:"docs"`
		FailedStarts SimpleStatistic `json:"failed_starts"`
		Jobs         struct {
			Adds          SimpleStatistic `json:"adds"`
			Crashed       SimpleStatistic `json:"crashed"`
			Crashes       SimpleStatistic `json:"crashes"`
			DuplicateAdds SimpleStatistic `json:"duplicate_adds"`
			Pending       SimpleStatistic `json:"pending"`
			Removes       SimpleStatistic `json:"removes"`
			Running       SimpleStatistic `json:"running"`
			Starts        SimpleStatistic `json:"starts"`
			Stops         SimpleStatistic `json:"stops"`
			Total         SimpleStatistic `json:"total"`
		} `json:"jobs"`
		Requests  SimpleStatistic `json:"requests"`
		Responses struct {
			Failure SimpleStatistic `json:"failure"`
			Success SimpleStatistic `json:"success"`
		} `json:"responses"`
		StreamResponses struct {
			Failure SimpleStatistic `json:"failure"`
			Success SimpleStatistic `json:"success"`
		} `json:"stream_responses"`
		WorkerDeaths   SimpleStatistic `json:"worker_deaths"`
		WorkersStarted SimpleStatistic `json:"workers_started"`
	} `json:"couch_replicator"`
	CouchDB struct {
		AuthCacheHits      SimpleStatistic   `json:"auth_cache_hits"`
		AuthCacheMisses    SimpleStatistic   `json:"auth_cache_misses"`
		CollectResultsTime DetailedStatistic `json:"collect_results_time"`
		CouchServer        struct {
			LruSkip SimpleStatistic `json:"lru_skip"`
		} `json:"couch_server"`
		DatabaseReads   SimpleStatistic   `json:"database_reads"`
		DatabaseWrites  SimpleStatistic   `json:"database_writes"`
		DBOpenTime      DetailedStatistic `json:"db_open_time"`
		DBInfo          DetailedStatistic `json:"dbinfo"`
		DocumentInserts SimpleStatistic   `json:"document_inserts"`
		DocumentWrites  SimpleStatistic   `json:"document_writes"`
		HTTPD           struct {
			AbortedRequests          SimpleStatistic   `json:"aborted_requests"`
			BulkDocs                 DetailedStatistic `json:"bulk_docs"`
			BulkRequests             SimpleStatistic   `json:"bulk_requests"`
			ClientsRequestingChanges SimpleStatistic   `json:"clients_requesting_changes"`
			Requests                 SimpleStatistic   `json:"requests"`
			TemporaryViewReads       SimpleStatistic   `json:"temporary_view_reads"`
			ViewReads                SimpleStatistic   `json:"view_reads"`
		} `json:"httpd"`
		HTTPDRequestMethods struct {
			Copy    SimpleStatistic `json:"COPY"`
			Delete  SimpleStatistic `json:"DELETE"`
			Get     SimpleStatistic `json:"GET"`
			Head    SimpleStatistic `json:"HEAD"`
			Options SimpleStatistic `json:"OPTIONS"`
			Post    SimpleStatistic `json:"POST"`
			Put     SimpleStatistic `json:"PUT"`
		} `json:"httpd_request_methods"`
		HTTPDStatusCodes struct {
			Two00   SimpleStatistic `json:"200"`
			Two01   SimpleStatistic `json:"201"`
			Two02   SimpleStatistic `json:"202"`
			Two04   SimpleStatistic `json:"204"`
			Two06   SimpleStatistic `json:"206"`
			Three01 SimpleStatistic `json:"301"`
			Three02 SimpleStatistic `json:"302"`
			Three04 SimpleStatistic `json:"304"`
			Four00  SimpleStatistic `json:"400"`
			Four01  SimpleStatistic `json:"401"`
			Four03  SimpleStatistic `json:"403"`
			Four04  SimpleStatistic `json:"404"`
			Four05  SimpleStatistic `json:"405"`
			Four06  SimpleStatistic `json:"406"`
			Four09  SimpleStatistic `json:"409"`
			Four12  SimpleStatistic `json:"412"`
			Four13  SimpleStatistic `json:"413"`
			Four14  SimpleStatistic `json:"414"`
			Four15  SimpleStatistic `json:"415"`
			Four16  SimpleStatistic `json:"416"`
			Four17  SimpleStatistic `json:"417"`
			Five00  SimpleStatistic `json:"500"`
			Five01  SimpleStatistic `json:"501"`
			Five03  SimpleStatistic `json:"503"`
		} `json:"httpd_status_codes"`
		LocalDocumentWrites SimpleStatistic `json:"local_document_writes"`
		MrView              struct {
			Emits  SimpleStatistic `json:"emits"`
			MapDoc SimpleStatistic `json:"map_doc"`
		} `json:"mrview"`
		OpenDatabases SimpleStatistic `json:"open_databases"`
		OpenOsFiles   SimpleStatistic `json:"open_os_files"`
		QueryServer   struct {
			VduProcessTime DetailedStatistic `json:"vdu_process_time"`
			VduRejects     SimpleStatistic   `json:"vdu_rejects"`
		} `json:"query_server"`
		RequestTime DetailedStatistic `json:"request_time"`
	} `json:"couchdb"`
	DDocCache struct {
		Hit      SimpleStatistic `json:"hit"`
		Miss     SimpleStatistic `json:"miss"`
		Recovery SimpleStatistic `json:"recovery"`
	} `json:"ddoc_cache"`
	Fabric struct {
		DocUpdate struct {
			Errors            SimpleStatistic `json:"errors"`
			MismatchedErrors  SimpleStatistic `json:"mismatched_errors"`
			WriteQuorumErrors SimpleStatistic `json:"write_quorum_errors"`
		} `json:"doc_update"`
		OpenShard struct {
			Timeouts SimpleStatistic `json:"timeouts"`
		} `json:"open_shard"`
		ReadRepairs struct {
			Failure SimpleStatistic `json:"failure"`
			Success SimpleStatistic `json:"success"`
		} `json:"read_repairs"`
		Worker struct {
			Timeouts SimpleStatistic `json:"timeouts"`
		} `json:"worker"`
	} `json:"fabric"`
	GlobalChanges struct {
		DBWrites               SimpleStatistic `json:"db_writes"`
		EventDocConflict       SimpleStatistic `json:"event_doc_conflict"`
		ListenerPendingUpdates SimpleStatistic `json:"listener_pending_updates"`
		Rpcs                   SimpleStatistic `json:"rpcs"`
		ServerPendingUpdates   SimpleStatistic `json:"server_pending_updates"`
	} `json:"global_changes"`
	Mem3 struct {
		ShardCache struct {
			Eviction SimpleStatistic `json:"eviction"`
			Hit      SimpleStatistic `json:"hit"`
			Miss     SimpleStatistic `json:"miss"`
		} `json:"shard_cache"`
	} `json:"mem3"`
	Pread struct {
		ExceedEOF   SimpleStatistic `json:"exceed_eof"`
		ExceedLimit SimpleStatistic `json:"exceed_limit"`
	} `json:"pread"`
	Rexi struct {
		Buffered SimpleStatistic `json:"buffered"`
		Down     SimpleStatistic `json:"down"`
		Dropped  SimpleStatistic `json:"dropped"`
		Streams  struct {
			Timeout struct {
				InitStream SimpleStatistic `json:"init_stream"`
				Stream     SimpleStatistic `json:"stream"`
				WaitForAck SimpleStatistic `json:"wait_for_ack"`
			} `json:"timeout"`
		} `json:"streams"`
	} `json:"rexi"`
}

// SimpleStatistic is the type used by CouchDB 2 to represent a statistic
// which does not have detailed statistical data associated with it meaning
// that it is just a single value.
type SimpleStatistic struct {
	Desc  string  `json:"desc"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

// DetailedStatistic is the type used by CouchDB 2 to represent a statistic
// with included statistical data such as the mean, max, min, median,
// kurtosis etc.
type DetailedStatistic struct {
	Desc  string `json:"desc"`
	Type  string `json:"type"`
	Value struct {
		ArithmeticMean    float64     `json:"arithmetic_mean"`
		GeometricMean     float64     `json:"geometric_mean"`
		HarmonicMean      float64     `json:"harmonic_mean"`
		Histogram         [][]float64 `json:"histogram"`
		Kurtosis          float64     `json:"kurtosis"`
		Max               float64     `json:"max"`
		Median            float64     `json:"median"`
		Min               float64     `json:"min"`
		N                 float64     `json:"n"`
		Percentile        [][]float64 `json:"percentile"`
		Skewness          float64     `json:"skewness"`
		StandardDeviation float64     `json:"standard_deviation"`
		Variance          float64     `json:"variance"`
	} `json:"value"`
}

// AllStatistics gets all of the available statistics from the server.
func (con *CouchDB2Connection) AllStatistics() (Statistics2, error) {
	var stats Statistics2
	_, err := con.unmarshalRequest("GET", "/_node/_local/_stats", NewURLOptions(), nil, &stats)
	if err != nil {
		return Statistics2{}, err
	}
	return stats, nil
}

// Statistic loads a single specific statistic from the server by category & name.
func (con *CouchDB2Connection) Statistic(category, name string) (Statistics2, error) {
	var stats Statistics2
	_, err := con.unmarshalRequest("GET", fmt.Sprintf("/_node/_local/_stats/%s/%s", category, name), NewURLOptions(), nil, &stats)
	if err != nil {
		return Statistics2{}, err
	}
	return stats, nil
}

// AllClusterStatistics gets all of the available statistics from the server.
func (con *Connection) AllClusterStatistics() (Statistics2, error) {
	var stats Statistics2
	_, err := con.unmarshalRequest("GET", "/_stats", NewURLOptions(), nil, &stats)
	if err != nil {
		return Statistics2{}, err
	}
	return stats, nil
}

// ClusterStatistic loads a single specific statistic from the server by category & name.
func (con *Connection) ClusterStatistic(category, name string) (Statistics2, error) {
	var stats Statistics2
	_, err := con.unmarshalRequest("GET", fmt.Sprintf("/_stats/%s/%s", category, name), NewURLOptions(), nil, &stats)
	if err != nil {
		return Statistics2{}, err
	}
	return stats, nil
}
