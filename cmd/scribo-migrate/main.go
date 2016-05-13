package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bbengfort/scribo/scribo"
	"github.com/codegangsta/cli"
)

func main() {

	// Instantiate the command line application
	app := cli.NewApp()
	app.Name = "scribo-migrate"
	app.Usage = "executes the latest migration SQL to the database"
	app.Version = scribo.Version
	app.Author = "Benjamin Bengfort"
	app.Email = "benjamin@bengfort.com"
	app.Action = migrateDatabase

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "all",
			Usage: "run all migrations from the beginning",
		},
		cli.BoolTFlag{
			Name:  "latest",
			Usage: "run only the latest migration (default) and exit",
		},
	}

	// Run the command line application
	app.Run(os.Args)
}

func migrateDatabase(ctx *cli.Context) error {
	var rows int64

	// Create a list of files in the migrations directory
	migrations := loadMigrations()

	// Connect and verify that the connection is open
	db := scribo.ConnectDB()
	check(db.Ping())

	if ctx.NArg() > 1 {
		return cli.NewExitError("Only one migration can be specified at a time.", 1)
	}

	if ctx.NArg() == 1 {
		idx, err := strconv.Atoi(ctx.Args()[0])
		check(err)
		if idx >= len(migrations) {
			return cli.NewExitError(fmt.Sprintf("No migration #%d", idx), 1)
		}

		rows += executeMigration(migrations[idx], db)
	} else if ctx.Bool("all") {
		for _, migration := range migrations {
			rows += executeMigration(migration, db)
		}
	} else if ctx.BoolT("latest") {
		latest := migrations[len(migrations)-1]
		rows += executeMigration(latest, db)
	} else {
		return cli.NewExitError("No migrations specified to execute!", 1)
	}

	fmt.Printf("Migration changed %d rows\n", rows)
	return nil
}

// Load the schema from the migration files.
func loadMigrations() []string {

	// Create a list of files in the migrations directory
	files, err := filepath.Glob("migrations/[0-9][0-9][0-9][0-9]-*.sql")
	check(err)

	return files
}

// Execute the migrations to the database.
func executeMigration(path string, db *sql.DB) int64 {

	data, err := ioutil.ReadFile(path)
	check(err)

	result, err := db.Exec(string(data))
	check(err)

	rows, err := result.RowsAffected()
	check(err)

	return rows
}

func check(e error) {
	if e != nil {
		cli.NewExitError(e.Error(), 1)
	}
}
