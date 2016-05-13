package scribo

import (
	"database/sql"
	"log"
	"os"

	// Loads the driver for database/sql
	_ "github.com/jackc/pgx/stdlib"
)

// ConnectDB establishes a connection to the PostgreSQL database
func ConnectDB() *sql.DB {
	dbURL := os.Getenv("DATABASE_URL")
	log.Printf("Connecting to database at %s", dbURL)

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// GetNode by ID, attempts to return the node or an error otherwise.
func GetNode(db *sql.DB, id int64) (Node, error) {
	var n Node

	row := db.QueryRow("SELECT * FROM nodes WHERE id = $1", id)
	err := row.Scan(&n.ID, &n.Name, &n.Address, &n.DNS, &n.Key, &n.Created, &n.Updated)

	switch {
	case err == sql.ErrNoRows:
		return Node{}, nil
	case err != nil:
		return n, err
	default:
		return n, nil
	}
}

// GetNodeByName attempts to return the node from a name or an error otherwise.
func GetNodeByName(db *sql.DB, name string) (Node, error) {
	var n Node

	row := db.QueryRow("SELECT * FROM nodes WHERE name = $1", name)
	err := row.Scan(&n.ID, &n.Name, &n.Address, &n.DNS, &n.Key, &n.Created, &n.Updated)

	switch {
	case err == sql.ErrNoRows:
		return Node{}, nil
	case err != nil:
		return n, err
	default:
		return n, nil
	}

}

// NodeExists tests if the given ID is associated with a node.
func NodeExists(db *sql.DB, id int64) (bool, error) {
	var exists bool
	query := "select exists(select 1 from nodes where id=$1)"
	row := db.QueryRow(query, id)
	err := row.Scan(&exists)

	return exists, err
}

// NodeExistsByName tests if the given name is associated with a node.
func NodeExistsByName(db *sql.DB, name string) (bool, error) {
	var exists bool
	query := "select exists(select 1 from nodes where name=$1)"
	row := db.QueryRow(query, name)
	err := row.Scan(&exists)

	return exists, err
}

// FetchNodes returns a collection of nodes, ordered by the updated timestamp.
// This function expects you to limit the size of the collection by specifying
// the maximum number of nodes to return in the Nodes collection.
func FetchNodes(db *sql.DB, limit int) (Nodes, error) {
	var nodes Nodes

	rows, err := db.Query("SELECT * FROM nodes ORDER BY updated DESC LIMIT $1", limit)
	if err != nil {
		return nodes, err
	}

	for rows.Next() {
		var n Node
		if err := rows.Scan(&n.ID, &n.Name, &n.Address, &n.DNS, &n.Key, &n.Created, &n.Updated); err != nil {
			return nodes, err
		}

		nodes = append(nodes, n)
	}

	rows.Close()
	return nodes, nil
}

// GetPing by ID, attempts to return the ping or an error otherwise.
func GetPing(db *sql.DB, id int64) (Ping, error) {
	var p Ping

	row := db.QueryRow("SELECT * FROM pings WHERE id = $1", id)
	err := row.Scan(&p.ID, &p.Source, &p.Target, &p.Payload, &p.Latency, &p.Timeout, &p.Created, &p.Updated)

	switch {
	case err == sql.ErrNoRows:
		return Ping{}, nil
	case err != nil:
		return p, err
	default:
		return p, nil
	}

}

// PingExists tests if the given ID is associated with a ping.
func PingExists(db *sql.DB, id int64) (bool, error) {
	var exists bool
	query := "select exists(select 1 from pings where id=$1)"
	row := db.QueryRow(query, id)
	err := row.Scan(&exists)

	return exists, err
}

// FetchPings returns a collection of pings, ordered by the created timestamp.
// This function expects you to limit the size of the collection by specifying
// the maximum number of pings to return in the Pings collection.
func FetchPings(db *sql.DB, limit int) (Pings, error) {
	var pings Pings

	rows, err := db.Query("SELECT * FROM pings ORDER BY created DESC LIMIT $1", limit)
	if err != nil {
		return pings, err
	}

	for rows.Next() {
		var p Ping
		if err := rows.Scan(&p.ID, &p.Source, &p.Target, &p.Payload, &p.Latency, &p.Timeout, &p.Created, &p.Updated); err != nil {
			return pings, err
		}

		pings = append(pings, p)
	}

	rows.Close()
	return pings, nil
}
