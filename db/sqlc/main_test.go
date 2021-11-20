package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/sajir-dev/goMasterClass/utils"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var config utils.Config
	var err error
	config, err = utils.LoadConfig("../../")
	if err != nil {
		log.Fatalln("could not load config", err)
	}
	// passing dbDriver and dbSource for creating db
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("unable to connect to the db")
	}

	// passing queries for new connection
	testQueries = New(testDB)
	os.Exit(m.Run())
}
