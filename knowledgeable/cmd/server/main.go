package main

import (
	"database/sql"
	"html/template"
	"knowledgeable/internal/auth"
	"knowledgeable/internal/pages"
	"knowledgeable/internal/users"
	"log"
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

	// templates
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	// user
	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService, tmpl)

	// pages
	pageRepo := pages.NewRepository(db)
	pageService := pages.NewService(pageRepo)
	pageHandler := pages.NewHandler(pageService)

	// auth
	authHandler := auth.NewHandler(userService, tmpl)

	log.Println("Dependencies wired successfully")

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

		if _, ok := auth.Get(cookie.Value); ok {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	http.Handle("/page",
		auth.Middleware(http.HandlerFunc(pageHandler.ViewPage)),
	)
	http.Handle("/search",
		auth.Middleware(http.HandlerFunc(pageHandler.Search)),
	)
	http.Handle("/api/search",
	auth.Middleware(http.HandlerFunc(pageHandler.SearchAPI)),
	)

	http.Handle("/users",
		auth.Middleware(http.HandlerFunc(userHandler.GetAll)),
	)

	http.Handle("/register",
		auth.Middleware(http.HandlerFunc(userHandler.Register)),
	)
	http.HandleFunc("/api/register", userHandler.RegisterAPI)

	http.HandleFunc("/logout", authHandler.Logout)

	http.HandleFunc("/login", authHandler.Login)
	
	http.HandleFunc("/api/login", authHandler.LoginAPI)

	http.Handle("/dashboard",
		auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := tmpl.ExecuteTemplate(w, "dashboard.html", nil); err != nil {
				http.Error(w, "template error", http.StatusInternalServerError)
			}
		})),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
