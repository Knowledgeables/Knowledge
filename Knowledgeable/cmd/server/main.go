package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"knowledgeable/internal/users"
)


func main() {
	
	
db, err := sql.Open("sqlite3", "knowledge.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	userRepo := users.NewRepository(db)

	_ = users.NewService(userRepo)

	// user handler below that takes userservice as an argument.

	log.Println("Dependencies wired successfully")




	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, map[string]string{
			"Title": "Knowledgeable",
		})
	})

	http.ListenAndServe(":8080", nil)
}
