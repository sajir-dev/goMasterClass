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

func TestMain(m *testing.M) {
	// passing dbDriver and dbSource for creating db
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("unable to connect to the db")
	}

	// passing queries for new connection
	testQueries = New(conn)
	os.Exit(m.Run())
}
