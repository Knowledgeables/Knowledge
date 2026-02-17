package main

import (
	"knowledgeable/internal/pages"
	"database/sql"
	"knowledgeable/internal/users"
	"log"
	"net/http"
	"os"
	_ "modernc.org/sqlite"

)

func main() {

	// db setup
	db, err := sql.Open("sqlite", "whoknows.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	schema, err := os.ReadFile("knowledge.sql")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		log.Fatal(err)
	}

	// dependency injection
	userRepo := users.NewRepository(db)

	userService := users.NewService(userRepo)

	userHandler := users.NewHandler(userService)

	pageRepo := pages.NewRepository(db)
	pageService := pages.NewService(pageRepo)
	pageHandler := pages.NewHandler(pageService)


	// user handler below that takes userservice as an argument.

	log.Println("Dependencies wired successfully")


	http.HandleFunc("/", pageHandler.Search)

	
	// user handler below that takes userservice as an argument.
	http.HandleFunc("/users", userHandler.GetAll)
	
	// page handler below that takes pageservice as an argument.
	http.HandleFunc("/pages", pageHandler.GetAll)


	log.Fatal(http.ListenAndServe(":8080", nil))
}
