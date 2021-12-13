package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/sajir-dev/goMasterClass/api"
	db "github.com/sajir-dev/goMasterClass/db/sqlc"
	"github.com/sajir-dev/goMasterClass/utils"
)

func main() {
	config, err := utils.LoadConfig("./")
	if err != nil {
		log.Fatalln("could not load config", config)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("could not connect with db")
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatalln("could not create server")
	}

	if err := server.Start(config.Port); err != nil {
		log.Fatalln(err)
	}
}
