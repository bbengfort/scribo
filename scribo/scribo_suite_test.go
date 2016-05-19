package scribo_test

import (
	"database/sql"
	"os"
	"strings"

	"github.com/bbengfort/scribo/scribo"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestScribo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scribo Suite")
}

var (
	db         *sql.DB
	app        *scribo.App
	migrations []string
	tables     []string
)

// Establish a connection to the database before tests are run.
var _ = BeforeSuite(func() {
	// Load the .env file if it exists
	godotenv.Load()

	var err error
	By("Connecting to a testing database")

	// Test the database url to make sure it ends in -test
	dbURL := os.Getenv("TEST_DATABASE_URL")
	Expect(strings.HasSuffix(dbURL, "-test")).To(BeTrue(), "The test database url should end in -test")

	// Establish the database connection
	db, err = connectTestDB()
	Expect(err).NotTo(HaveOccurred(), "Could not establish a connection to the database")

	// Expect that a database ping does not cause an error
	Expect(db.Ping()).NotTo(HaveOccurred(), "Could not ping the database")

	By("Loading the schema from migration files on disk")

	// Load the migrations from disk
	migrations = loadMigrations()

	// Execute the migrations to the database
	executeMigrations(migrations)

	// Load the tables from the database
	tables = listTables()

	By("Creating a test application")
	app = createTestApp()
})

// Clean up the database connections after the test suite is run.
var _ = AfterSuite(func() {
	// Drop all the tables from the database so we can test again!
	dropTables(tables)

	// Expect that an error on closing doesn't occur
	Ω(db.Close()).Should(Succeed(), "Could not close the database")
})
