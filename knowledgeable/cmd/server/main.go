package main

import (
	"html/template"
	"knowledgeable/internal/auth"
	"knowledgeable/internal/db"
	"knowledgeable/internal/pages"
	"knowledgeable/internal/users"
	"knowledgeable/internal/web"
	"log"
	"net/http"
	"os"

	_ "knowledgeable/docs"

	_ "modernc.org/sqlite"
)

func main() {

	// db setup
	database := db.Init(os.Getenv("DB_PATH"), "knowledge.sql")
	defer database.Close()

	// seed
	if os.Getenv("APP_ENV") == "dev" {
		log.Println("Seeding database (dev)")

		seed, err := os.ReadFile("seed-dev.sql")
		if err != nil {
			log.Fatal(err)
		}

		if _, err := database.Exec(string(seed)); err != nil {
			log.Fatal(err)
		}
	}

	// templates
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	// user
	userRepo := users.NewRepository(database)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService, tmpl)

	// pages
	pageRepo := pages.NewRepository(database)
	pageService := pages.NewService(pageRepo)
	pageHandler := pages.NewHandler(pageService)

	// auth
	authHandler := auth.NewHandler(userService, tmpl)

	// routes
	web.SetupRoutes(
		userHandler,
		pageHandler,
		authHandler,
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}

			if err := tmpl.ExecuteTemplate(w, "dashboard.html", nil); err != nil {
				http.Error(w, "template error", http.StatusInternalServerError)
			}
		},
	)

	log.Println("Server running on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
