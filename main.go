package main

import (
	"database/sql"
	"fmt"
	"log"

	"example.com/simple_bank/api"
	"example.com/simple_bank/db/sqlc"
	"example.com/simple_bank/db/util"

	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load env variables")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	fmt.Println(err)

	if err != nil {
		log.Fatal("unable to connect to DB. Error Details: ", err)
	}

	store := sqlc.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("unable to start server. err details:", err)
	}

	err = server.StartServer(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
