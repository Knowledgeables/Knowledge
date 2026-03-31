// @title Knowledge API
// @version 1.0
// @description API for Knowledge service
// @host localhost:8080
// @BasePath /
package main

import (
	"database/sql"
	"html/template"
	"knowledgeable/internal/auth"
	"knowledgeable/internal/pages"
	"knowledgeable/internal/users"
	"log"
	"net/http"
	"os"

	_ "knowledgeable/docs"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "modernc.org/sqlite"
)

func main() {

	// db setup

	dbPath := os.Getenv("DB_PATH")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Using database:", dbPath)

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

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

	// seed
	if os.Getenv("APP_ENV") == "dev" {

		seed, err := os.ReadFile("seed-dev.sql")
		if err != nil {
			log.Fatal(err)
		}

		if _, err := db.Exec(string(seed)); err != nil {
			log.Fatal(err)
		}
	}
	// Swagger UI
	http.Handle("/swagger/", httpSwagger.Handler())

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

	http.HandleFunc("/page", pageHandler.ViewPage)

	http.HandleFunc("/search", pageHandler.Search)

	http.HandleFunc("/api/search", pageHandler.SearchAPI)

	http.HandleFunc("/register", userHandler.Register)
	http.HandleFunc("/api/register", userHandler.RegisterAPI)

	http.HandleFunc("/logout", authHandler.Logout)

	http.HandleFunc("/login", authHandler.Login)

	http.HandleFunc("/api/login", authHandler.LoginAPI)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "dashboard.html", nil); err != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
		}
	})

	// Metrics endpoint used by Prometheus and visualized in Grafana
	http.Handle("/metrics", promhttp.Handler())

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
