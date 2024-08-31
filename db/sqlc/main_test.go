package sqlc

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"example.com/simple_bank/db/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	var err error

	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("unable to fetch env variables")
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("unable to connect to DB. Error Details: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
