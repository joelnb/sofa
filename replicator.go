package sofa

import (
	"encoding/json"
)

// Replication represents a single document in the CouchDB '_replicator' database.
// Each one of these documents represents a replication between two databases. The full API
// documentation is located at: http://docs.couchdb.org/en/2.0.0/replication/replicator.html
type Replication struct {
	DocumentMetadata

	CreateTarget bool `json:"create_target,omitempty"`
	Continuous   bool `json:"continuous,omitempty"`

	Connections     float64 `json:"connection_timeout,omitempty"`
	Retries         float64 `json:"retries_per_request,omitempty"`
	HTTPConnections float64 `json:"http_connections,omitempty"`

	Target string `json:"target,omitempty"`
	Source string `json:"source,omitempty"`

	Owner          string `json:"owner,omitempty"`
	AdditionalType string `json:"additionalType,omitempty"`

	ReplicationState     string                 `json:"_replication_state,omitempty"`
	ReplicationStateTime string                 `json:"_replication_state_time,omitempty"`
	ReplicationID        string                 `json:"_replication_id,omitempty"`
	ReplicationStats     map[string]interface{} `json:"_replication_stats,omitempty"`

	UserContext map[string]interface{} `json:"user_ctx,omitempty"`
	QueryParams map[string]interface{} `json:"query_params,omitempty"`
}

// NewReplication creates a new replication on the CouchDB server. The source and target
// may either be a local database or a database on a remote server.
func (con *Connection) NewReplication(id, source, target string) Replication {
	return Replication{
		DocumentMetadata: DocumentMetadata{ID: id},
		Source:           source,
		Target:           target,
	}
}

// Replications returns the full list of replications currently active on the server.
func (con *Connection) Replications() ([]Replication, error) {
	db := con.Database("_replicator")

	table, err := db.AllDocuments()
	if err != nil {
		return nil, err
	}

	repls := []Replication{}
	for _, rawRepl := range table.RawDocuments() {
		repl := Replication{}
		if err := json.Unmarshal(rawRepl, &repl); err != nil {
			return nil, err
		}
		repls = append(repls, repl)
	}
	return repls, nil
}

// Replication gets a particular Replication from the server by ID. A revision can
// also be specified to retrieve a particular revision of the document.
func (con *Connection) Replication(id, rev string) (Replication, error) {
	db := con.Database("_replicator")

	repl := Replication{}
	_, err := db.Get(&repl, id, rev)
	if err != nil {
		return repl, err
	}

	return repl, nil
}

// DeleteReplication removes a replication from the replicator database and cancels it.
func (con *Connection) DeleteReplication(repl Replication) error {
	db := con.Database("_replicator")

	_, err := db.Delete(repl)
	if err != nil {
		return err
	}

	return nil
}
