package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	// passing dbDriver and dbSource for creating db
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("unable to connect to the db")
	}

	// passing queries for new connection
	testQueries = New(testDB)
	os.Exit(m.Run())
}
