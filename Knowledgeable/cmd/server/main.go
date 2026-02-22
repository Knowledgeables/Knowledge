package main

import (
	"knowledgeable/internal/pages"
	"database/sql"
	"knowledgeable/internal/auth"
	"knowledgeable/internal/users"
	"log"
	"html/template"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
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

	// user
	userRepo := users.NewRepository(db)

	userService := users.NewService(userRepo)

	userHandler := users.NewHandler(userService)

	pageRepo := pages.NewRepository(db)
	pageService := pages.NewService(pageRepo)
	pageHandler := pages.NewHandler(pageService)

	// auth
	authHandler := auth.NewHandler(userService)

	// user handler below that takes userservice as an argument.

	log.Println("Dependencies wired successfully")


	http.Handle("/page",
		auth.Middleware(http.HandlerFunc(pageHandler.ViewPage)),
	)
	http.Handle("/search",
		auth.Middleware(http.HandlerFunc(pageHandler.Search)),
	)
	http.Handle("/api/search",
		auth.Middleware(http.HandlerFunc(pageHandler.SearchAPI)),
	)
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		_, ok := auth.Get(cookie.Value)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/search-engine", http.StatusSeeOther)
	})

	http.Handle("/users",
		auth.Middleware(http.HandlerFunc(userHandler.GetAll)),
	)

	http.HandleFunc("/logout", authHandler.Logout)
	http.HandleFunc("/login", authHandler.Login)

	http.Handle("/dashboard",
	auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))
		tmpl.Execute(w, nil)
	})),
)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
