package utils

import (
	"database/sql"
	"fibo_go_server/config"
	"log"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	databaseUrl := config.GetDatabaseURL()
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal("Error opening database: ", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
		return nil, err
	}

	createSchemaAndTables(db)

	return db, nil
}
