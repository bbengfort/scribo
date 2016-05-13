package scribo_test

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/bbengfort/scribo/scribo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// The test suite for the database module
var _ = Describe("Database", func() {

	It("should only connect to a database ending in -test", func() {
		dbURL := os.Getenv("TEST_DATABASE_URL")
		Ω(dbURL).Should(HaveSuffix("-test"))

	})

	It("should be able to connect to a test database", func() {
		Ω(db.Ping()).Should(Succeed())
	})

	It("should have only one migration file", func() {
		Ω(migrations).Should(HaveLen(1))
	})

	It("should only have two visible tables", func() {
		Ω(tables).Should(HaveLen(2))
	})

	AfterEach(func() {
		truncateTables(tables)
	})

	Describe("Node", func() {

		Context("when saving a node to the database", func() {

			It("should create a node if it doesn't have an ID", func() {
				exists, err := scribo.NodeExists(db, 1)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(exists).Should(BeFalse())

				node := &scribo.Node{
					Name:    "apollo",
					Address: "108.51.64.223",
				}

				threshold, err := time.ParseDuration("1s")
				Ω(err).ShouldNot(HaveOccurred())

				created, err := node.Save(db)
				Ω(created).Should(BeTrue())
				Ω(err).ShouldNot(HaveOccurred())
				Ω(node.ID).Should(BeNumerically(">", 0))
				Ω(node.Created).Should(BeTemporally("~", time.Now(), threshold))
				Ω(node.Updated).Should(BeTemporally("~", time.Now(), threshold))

				exists, err = scribo.NodeExists(db, node.ID)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(exists).Should(BeTrue())
			})

			It("should update a node if it already has an ID", func() {

				_, err := db.Exec("INSERT INTO nodes (name, address) VALUES ($1, $2)", "apollo", "108.51.64.223")
				Ω(err).ShouldNot(HaveOccurred())

				node, err := scribo.GetNodeByName(db, "apollo")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(node.ID).Should(Equal(int64(1)))

				node.DNS = "bryant.bengfort.com"
				created, err := node.Save(db)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(created).Should(BeFalse())
				Ω(node.Updated).Should(BeTemporally(">", node.Created))

				node2, err := scribo.GetNode(db, 1)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(node2.DNS).Should(Equal("bryant.bengfort.com"))

			})

			It("should be able to delete a node", func() {
				_, err := db.Exec("INSERT INTO nodes (name, address) VALUES ($1, $2)", "apollo", "108.51.64.223")
				Ω(err).ShouldNot(HaveOccurred())

				node, err := scribo.GetNodeByName(db, "apollo")
				Ω(err).ShouldNot(HaveOccurred())

				deleted, err := node.Delete(db)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(deleted).Should(BeTrue())

				exists, err := scribo.NodeExists(db, 1)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(exists).Should(BeFalse())
			})

		})

		Context("when fetching a collection of nodes from the database", func() {

			It("should return an empty list of nodes", func() {
				nodes, err := scribo.FetchNodes(db, 50)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(nodes).Should(HaveLen(0))
			})

			It("should return an ordered, limited list of nodes", func() {
				var nodes scribo.Nodes
				nodes = append(nodes, scribo.Node{Name: "test1", Address: "127.0.0.1"})
				nodes = append(nodes, scribo.Node{Name: "test2", Address: "127.0.0.2"})
				nodes = append(nodes, scribo.Node{Name: "test3", Address: "127.0.0.3"})

				// Create the nodes
				for _, node := range nodes {
					created, err := node.Save(db)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(created).Should(BeTrue())
				}

				// Fetch the nodes collection
				collection, err := scribo.FetchNodes(db, 2)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(collection).Should(HaveLen(2))
			})

		})

	})

	Describe("Ping", func() {

		Context("when saving a ping to the database", func() {

			It("should create a ping if it doesn't have an ID", func() {
				Skip("Ping database testing not implemented yet.")
			})

			It("should update a ping if it already has an ID", func() {
				Skip("Ping database testing not implemented yet.")
			})

			It("should be able to delete a ping", func() {
				Skip("Ping database testing not implemented yet.")
			})

		})

		Context("when fetching a collection of pings from the database", func() {

			It("should return an ordered, limited list of pings", func() {
				Skip("Ping database testing not implemented yet.")
			})

		})

	})

})

// Establish a connection to the test database.
func connectTestDB() (*sql.DB, error) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	fmt.Fprintln(GinkgoWriter, dbURL)

	return sql.Open("pgx", dbURL)
}

// Load the schema from the migration files.
func loadMigrations() []string {

	// Create a list of files in the migrations directory
	files, err := filepath.Glob("../migrations/[0-9][0-9][0-9][0-9]-*.sql")
	Expect(err).NotTo(HaveOccurred(), "Could not list migration files from migrations directory!")

	var queries []string
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		Expect(err).NotTo(HaveOccurred(), "Could not read migration file from migrations directory!")
		queries = append(queries, string(data))
	}

	return queries
}

// Execute the migrations to the database.
func executeMigrations(queries []string) {
	for _, query := range queries {
		_, err := db.Exec(query)
		Expect(err).NotTo(HaveOccurred(), "Could not execute a migration SQL file to database!")
	}
}

// Get a listing of all the tables in the database
func listTables() []string {
	var tables []string
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE'"
	rows, err := db.Query(query)
	Expect(err).NotTo(HaveOccurred(), "Could not query the information schema for table names!")

	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		Expect(err).NotTo(HaveOccurred(), "Could not append a table name from the query into tables!")

		tables = append(tables, table)
	}

	return tables
}

// Truncate all the tables in the database.
func truncateTables(tables []string) {

	// Start the truncation transaction
	txn, err := db.Begin()
	Expect(err).NotTo(HaveOccurred(), "Could not begin table truncate transaction!")

	// Defer the rollback after the function returns
	// If the transaction was commited, this will do nothing
	defer txn.Rollback()
	query := "TRUNCATE TABLE %s RESTART IDENTITY CASCADE"

	for _, table := range tables {
		_, err := txn.Exec(fmt.Sprintf(query, table))
		Expect(err).NotTo(HaveOccurred(), "Could not execute %s table truncation", table)
	}

	// Commit the transaction
	Expect(txn.Commit()).To(Succeed(), "Could not commit the truncate table transaction!")

}

// Drop all tables in the database.
func dropTables(tables []string) {

	// Start the truncation transaction
	txn, err := db.Begin()
	Expect(err).NotTo(HaveOccurred(), "Could not begin drop table transaction!")

	// Defer the rollback after the function returns
	// If the transaction was commited, this will do nothing
	defer txn.Rollback()
	query := "DROP TABLE %s CASCADE"

	for _, table := range tables {
		_, err := txn.Exec(fmt.Sprintf(query, table))
		Expect(err).NotTo(HaveOccurred(), "Could not execute drop table %s", table)
	}

	// Commit the transaction
	Expect(txn.Commit()).To(Succeed(), "Could not commit the drop table transaction!")

}
