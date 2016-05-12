package scribo

import (
	"database/sql"
	"errors"
	"time"
)

// Node is a model that represents a participant in the network
type Node struct {
	ID      int64     `json:"id"`      // Unique ID of the node
	Name    string    `json:"name"`    // Name/DNS of the node
	Address string    `json:"address"` // IP Address of the node
	DNS     string    `json:"dns"`     // DNS Lookup for the node
	Key     string    `json:"-"`       // Authentication key of the node
	Created time.Time `json:"created"` // Datetime the node was created
	Updated time.Time `json:"updated"` // Datetime the node was updated
}

// Ping is a model that represents a latency report.
type Ping struct {
	ID      int64     `json:"id"`      // Unique ID of the ping
	Source  int64     `json:"source"`  // The ID of the source node
	Target  int64     `json:"target"`  // The ID of the target node
	Payload int       `json:"payload"` // The size in bytes of the payload
	Latency float64   `json:"latency"` // The time in ms of the round trip
	Timeout bool      `json:"timeout"` // Whether or not the request timed out
	Created time.Time `json:"created"` // Datetime the ping was created
	Updated time.Time `json:"updated"` // Datetime the ping was updated
}

// Nodes is a collection of node items for use elsewhere.
type Nodes []Node

// Pings is a collection of latency reports for use elsewhere.
type Pings []Ping

// Save a node struct to the database. This function checks if the node has an
// ID or not. If it does, it will execute a SQL UPDATE, otherwise it will
// execute a SQL INSERT. Returns a boolean if the node was created (INSERT) or
// False if the node was simply updated in the normal manner. This method also
// handles setting the Created and Updated timestamps on the node.
// TODO: Transform this into a prepared statement that we can run.
func (node *Node) Save(db *sql.DB) (bool, error) {
	if node.ID > 0 {
		// This is the UPDATE method, so return false.
		// Update the updated timestamp on the Node.
		node.Updated = time.Now()

		// Execute the query against the database
		query := "UPDATE nodes SET name=$1, address=$2, dns=$3, key=$4, updated=$5 WHERE id = $6"
		_, err := db.Exec(query, node.Name, node.Address, node.DNS, node.Key, node.Updated, node.ID)

		return false, err

	}

	// This is the INSERT method, so return true.
	// Set the Created and Updated timestamps on the Node.
	node.Created = time.Now()
	node.Updated = time.Now()

	// Execute the INSERT query against the database
	query := "INSERT INTO nodes (name, address, dns, key, created, updated) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	row := db.QueryRow(query, node.Name, node.Address, node.DNS, node.Key, node.Created, node.Updated)
	err := row.Scan(&node.ID)

	if err != nil {
		return false, err
	}

	return true, err
}

// Delete a node from the database. This method is obviously destructive and
// returns true if the number of rows affected is 1 or false otherwise.
func (node *Node) Delete(db *sql.DB) (bool, error) {
	if node.ID == 0 {
		return false, errors.New("The node doesn't have an ID accessible by the database")
	}

	query := "DELETE FROM nodes WHERE id=$1"
	res, err := db.Exec(query, node.ID)

	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()

	switch {
	case err != nil:
		return false, err
	case rows > 1:
		return false, errors.New("Multiple deletions from the database?!")
	case rows == 1:
		return true, nil
	case rows < 1:
		return false, nil
	default:
		return false, errors.New("Unknown case in node deletion")
	}

}

// Save a ping struct to the database. This function checks if the ping has an
// ID or not. If it does, it will execute a SQL UPDATE, otherwise it will
// execute a SQL INSERT. Returns a boolean if the ping was created (INSERT) or
// False if the ping was simply updated in the normal manner. This method also
// handles setting the Created and Updated timestamps on the ping.
// TODO: Transform this into a prepared statement that we can run.
func (ping *Ping) Save(db *sql.DB) (bool, error) {
	if ping.ID > 0 {
		// This is the UPDATE method, so return false.
		// Update the updated timestamp on the Node.
		ping.Updated = time.Now()

		// Execute the query against the database
		query := "UPDATE pings SET source_id=$1, target_id=$2, payload=$3, latency=$4, timeout=$5, updated=$6 WHERE id = $7"
		_, err := db.Exec(query, ping.Source, ping.Target, ping.Payload, ping.Latency, ping.Timeout, ping.Updated, ping.ID)

		return false, err

	}

	// This is the INSERT method, so return true.
	// Set the Created and Updated timestamps on the Node.
	ping.Created = time.Now()
	ping.Updated = time.Now()

	// Execute the INSERT query against the database
	query := "INSERT INTO pings (source_id, target_id, payload, latency, timeout, created, updated) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	row := db.QueryRow(query, ping.Source, ping.Target, ping.Payload, ping.Latency, ping.Timeout, ping.Created, ping.Updated)
	err := row.Scan(&ping.ID)

	if err != nil {
		return false, err
	}

	return true, err
}

// Delete a ping from the database. This method is obviously destructive and
// returns true if the number of rows affected is 1 or false otherwise.
func (ping *Ping) Delete(db *sql.DB) (bool, error) {
	if ping.ID == 0 {
		return false, errors.New("The ping doesn't have an ID accessible by the database")
	}

	query := "DELETE FROM pings WHERE id=$1"
	res, err := db.Exec(query, ping.ID)

	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()

	switch {
	case err != nil:
		return false, err
	case rows > 1:
		return false, errors.New("Multiple deletions from the database?!")
	case rows == 1:
		return true, nil
	case rows < 1:
		return false, nil
	default:
		return false, errors.New("Unknown case in ping deletion")
	}

}
