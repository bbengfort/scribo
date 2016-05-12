package scribo

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// Loads the driver for database/sql
	_ "github.com/jackc/pgx/stdlib"
)

// ConnectDB establishes a connection to the PostgreSQL database
func ConnectDB() *sql.DB {
	dbURL := os.Getenv("DATABASE_URL")
	fmt.Println(dbURL)

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// GetNode by ID, attempts to return the node or an error otherwise.
func GetNode(db *sql.DB, id uint64) (Node, error) {
	var n Node

	row := db.QueryRow("SELECT * FROM nodes WHERE id = $1", id)
	err := row.Scan(&n.ID, &n.Name, &n.Address, &n.DNS, &n.Key, &n.Created, &n.Updated)

	if err != nil {
		return Node{}, err
	}

	return n, nil
}

// GetNodeByName attempts to return the node from a name or an error otherwise.
func GetNodeByName(db *sql.DB, name string) (Node, error) {
	var n Node

	row := db.QueryRow("SELECT * FROM nodes WHERE name = $1", name)
	err := row.Scan(&n.ID, &n.Name, &n.Address, &n.DNS, &n.Key, &n.Created, &n.Updated)

	if err != nil {
		return Node{}, err
	}

	return n, nil
}

// FetchNodes returns a collection of nodes, ordered by the updated timestamp.
// This function expects you to limit the size of the collection by specifying
// the maximum number of nodes to return in the Nodes collection.
func FetchNodes(db *sql.DB, limit int) (Nodes, error) {
	var nodes Nodes

	rows, err := db.Query("SELECT * FROM nodes ORDER BY updated DESC LIMIT $1", limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var n Node
		if err := rows.Scan(&n.ID, &n.Name, &n.Address, &n.DNS, &n.Key, &n.Created, &n.Updated); err != nil {
			return nil, err
		}

		nodes = append(nodes, n)
	}

	return nodes, nil
}

// GetPing by ID, attempts to return the ping or an error otherwise.
func GetPing(db *sql.DB, id uint64) (Ping, error) {
	var p Ping

	row := db.QueryRow("SELECT * FROM pings WHERE id = $1", id)
	err := row.Scan(&p.ID, &p.Source, &p.Target, &p.Payload, &p.Latency, &p.Timeout, &p.Created, &p.Updated)

	if err != nil {
		return Ping{}, err
	}

	return p, nil
}

// FetchPings returns a collection of pings, ordered by the created timestamp.
// This function expects you to limit the size of the collection by specifying
// the maximum number of pings to return in the Pings collection.
func FetchPings(db *sql.DB, limit int) (Pings, error) {
	var pings Pings

	rows, err := db.Query("SELECT * FROM pings ORDER BY created DESC LIMIT $1", limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p Ping
		if err := rows.Scan(&p.ID, &p.Source, &p.Target, &p.Payload, &p.Latency, &p.Timeout, &p.Created, &p.Updated); err != nil {
			return nil, err
		}

		pings = append(pings, p)
	}

	return pings, nil
}
