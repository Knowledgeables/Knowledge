package main

import (
	"context"
	"html/template"
	_ "knowledgeable/docs"
	"knowledgeable/internal/auth"
	"knowledgeable/internal/db"
	"knowledgeable/internal/middleware"
	"knowledgeable/internal/pages"
	"knowledgeable/internal/users"
	"knowledgeable/internal/web"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		resp, err := http.Get("http://localhost:8080/health")
		if err != nil {
			os.Exit(1)
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		os.Exit(0)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	log.Println("[APP] Starting application")

	// db setup
	database := db.InitPostgres()

	// templates
	var tmplLoader func() *template.Template

	if os.Getenv("APP_ENV") == "dev" {
		tmplLoader = func() *template.Template {
			return template.Must(template.ParseGlob("templates/*.html"))
		}
	} else {
		tmpl := template.Must(template.ParseGlob("templates/*.html"))
		tmplLoader = func() *template.Template {
			return tmpl
		}
	}

	// serves static files such as css, js and images.
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	// user
	userRepo := users.NewRepository(database)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService, tmplLoader, auth.Create)

	// pages
	pageRepo := pages.NewRepository(database)
	pageService := pages.NewService(pageRepo)
	pageHandler := pages.NewHandler(pageService, tmplLoader)

	// auth
	authHandler := auth.NewHandler(userService, tmplLoader)

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

			tmpl := tmplLoader()

			if err := tmpl.ExecuteTemplate(w, "home.html", nil); err != nil {
				http.Error(w, "template error", http.StatusInternalServerError)
			}
		},
	)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      middleware.Tracking(middleware.PageView(http.DefaultServeMux)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start server async
	go func() {
		log.Println("[APP] Server running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[APP] Server failed: %v", err)
		}
	}()

	// shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("[APP] Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[APP] Forced shutdown: %v", err)
	}

	if err := database.Close(); err != nil {
		log.Printf("failed to close db: %v", err)
	}
}
