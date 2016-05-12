package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/bbengfort/scribo/scribo"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// TODO: accept command line arguments
// TODO: be smarter about sending queries to the database
func main() {

	// Create a list of files in the migrations directory
	migrations, err := filepath.Glob("migrations/[0-9][0-9][0-9][0-9]-*.sql")
	check(err)

	// Read the latest SQL statement into memory
	latest := migrations[len(migrations)-1]
	data, err := ioutil.ReadFile(latest)
	check(err)

	query := string(data)

	db := scribo.ConnectDB()

	// Verify that the connection is open
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Execute a query against the database
	result, err := db.Exec(query)
	check(err)

	rows, err := result.RowsAffected()
	check(err)

	fmt.Printf("Migration changed %d rows\n", rows)
}
