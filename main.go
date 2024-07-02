package main

import (
	"bank/api"
	db "bank/db/sqlc"
	"bank/util"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot laod config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to the db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
