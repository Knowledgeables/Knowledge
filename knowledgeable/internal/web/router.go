package web

import (
	"net/http"

	"knowledgeable/internal/auth"
	"knowledgeable/internal/pages"
	"knowledgeable/internal/users"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(
	userHandler *users.Handler,
	pageHandler *pages.Handler,
	authHandler *auth.Handler,
	dashboardHandler http.HandlerFunc,
) {

	http.Handle("/swagger/", httpSwagger.Handler())

	http.HandleFunc("/page", pageHandler.ViewPage)
	http.HandleFunc("/search", pageHandler.Search)
	http.HandleFunc("/api/search", pageHandler.SearchAPI)

	http.HandleFunc("/register", userHandler.Register)
	http.HandleFunc("/api/register", userHandler.RegisterAPI)

	http.HandleFunc("/logout", authHandler.Logout)
	http.HandleFunc("/login", authHandler.Login)
	http.HandleFunc("/api/login", authHandler.LoginAPI)

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", dashboardHandler)
}