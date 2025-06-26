package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/rs/cors"
	"github.com/thinkingojha/go-htmx/internal/config"
	"github.com/thinkingojha/go-htmx/internal/logger"
	"golang.org/x/time/rate"
)

// SecurityHeaders adds security headers to responses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// HTMX-specific headers
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		next.ServeHTTP(w, r)
	})
}

// RequestLogger logs all incoming requests with structured logging
func RequestLogger(next http.Handler) http.Handler {
	return handlers.LoggingHandler(logger.Logger.Out, next)
}

// Recovery middleware to handle panics gracefully
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Panic occurred: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RateLimiter creates a rate limiting middleware
func RateLimiter(requestsPerMinute int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(requestsPerMinute)/60, requestsPerMinute)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				logger.Warnf("Rate limit exceeded for %s", r.RemoteAddr)
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// CORS creates CORS middleware based on configuration
func CORS(cfg *config.Config) func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Configure based on your needs
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	})

	return c.Handler
}

// HealthCheck provides a simple health check endpoint
func HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	}
}

// Timeout middleware to prevent hanging requests
func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, timeout, "Request Timeout")
	}
}
