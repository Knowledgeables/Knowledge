package db

import (
	"database/sql"
	"log"
	"os"
)

func Init(dbPath string, schemaPath string) *sql.DB {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Using database:", dbPath)

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if schemaPath != "" {

		log.Println("Applying schema:", schemaPath)

		schema, err := os.ReadFile(schemaPath)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := db.Exec(string(schema)); err != nil {
			log.Fatal(err)
		}

	}

	return db
}
