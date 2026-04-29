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

	http.HandleFunc("/search", pageHandler.Search)
	http.HandleFunc("/api/search", pageHandler.SearchAPI)
	http.HandleFunc("/api/crawler/targets", pageHandler.CrawlerTargetsAPI)
	http.HandleFunc("/api/crawler/ingest", pageHandler.CrawlerIngestAPI)
	http.HandleFunc("/api/crawler/existing-urls", pageHandler.CrawlerExistingURLsAPI)

	http.HandleFunc("/register", userHandler.Register)
	http.HandleFunc("/api/register", userHandler.RegisterAPI)

	http.HandleFunc("/logout", authHandler.Logout)
	http.HandleFunc("/login", authHandler.Login)
	http.HandleFunc("/api/login", authHandler.LoginAPI)
	http.HandleFunc("/api/me", authHandler.Me)
	http.HandleFunc("/change-password", authHandler.ChangePassword)
	http.HandleFunc("/forgot-password", authHandler.ForgotPassword)
	http.HandleFunc("/reset-password", authHandler.ResetPassword)

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	http.HandleFunc("/", dashboardHandler)

}
