package main

import (
	"database/sql"
	"html/template"
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

	// user handler below that takes userservice as an argument.

	log.Println("Dependencies wired successfully")

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, map[string]string{
			"Title": "Knowledgeable",
		})
	})

	http.HandleFunc("/users", userHandler.GetAll)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
