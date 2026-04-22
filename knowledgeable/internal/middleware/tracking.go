package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"
	"os"
	"github.com/google/uuid"
)

type contextKey string

// key for sving tracking_id ind request context
const sessionKey contextKey = "tracking_id"

// cookie name used for tracking and used in browser
const cookieName = "tracking_id"

// tracking middleware ensures that there is ALWAYS a tracking_id either from a cookie or a completely new tracking id.
// It also injects tracking id in the request context
func Tracking(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var trackingID string

		// attempts to collect the already existd ing cookie
		cookie, err := r.Cookie(cookieName)
		// if cookie doesnt exist generate a new one
		if err != nil || cookie.Value == "" {
			trackingID = uuid.New().String()

			// save tracking_id in browser cookie with httponly
			http.SetCookie(w, &http.Cookie{ // #nosec G124
				Name:     cookieName,
				Value:    trackingID,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				Secure:   os.Getenv("APP_ENV") == "production",
			})
		} else { // if cookie already exists in browser from before. use it again.
			trackingID = cookie.Value
		}

		// put t racking id into request context
		ctx := context.WithValue(r.Context(), sessionKey, trackingID)

		// update the requiest with new context
		r = r.WithContext(ctx)

		// continue to next handler
		next.ServeHTTP(w, r)
	})
}

// gettrackingid collects tracking_id from a request context
func GetTrackingID(r *http.Request) string {
	if v := r.Context().Value(sessionKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func PageView(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
        next.ServeHTTP(wrapped, r)
        duration := time.Since(start)
        slog.Info("pageview", // #nosec G706 -- JSON handler escapes all values
            "event", "pageview",
            "tracking_id", GetTrackingID(r),
            "path", r.URL.Path,
            "method", r.Method,
            "status", wrapped.statusCode,
            "duration_ms", duration.Milliseconds(),
        )
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
