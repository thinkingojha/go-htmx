package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/thinkingojha/go-htmx/internal/config"
	"github.com/thinkingojha/go-htmx/internal/handlers"
	"github.com/thinkingojha/go-htmx/internal/logger"
	"github.com/thinkingojha/go-htmx/internal/middleware"
)

type Server struct {
	config     *config.Config
	httpServer *http.Server
	router     *mux.Router
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
		router: mux.NewRouter(),
	}
}

func (s *Server) setupRoutes() {
	// Health check endpoint (no middleware for monitoring systems)
	s.router.HandleFunc("/health", middleware.HealthCheck()).Methods("GET")

	// API routes with middleware
	api := s.router.PathPrefix("/").Subrouter()

	// Apply middleware in order
	api.Use(middleware.Recovery)
	api.Use(middleware.SecurityHeaders)
	api.Use(middleware.RequestLogger)
	api.Use(middleware.RateLimiter(s.config.Security.RateLimitRPM))
	api.Use(middleware.CORS(s.config))
	api.Use(middleware.Timeout(30 * time.Second))

	// Static file serving
	staticHandler := http.StripPrefix("/static/",
		http.FileServer(http.Dir(s.config.App.StaticDir)))
	api.PathPrefix("/static/").Handler(staticHandler)

	// Application routes
	api.HandleFunc("/", s.makeHTTPHandlerFunc(handlers.HomeHandler)).Methods("GET")
	api.HandleFunc("/info", s.makeHTTPHandlerFunc(handlers.ExpHandler)).Methods("GET")

	api.HandleFunc("/products", s.makeHTTPHandlerFunc(handlers.ProductHandler)).Methods("GET")

	// Blog routes
	api.HandleFunc("/blog", s.makeHTTPHandlerFunc(handlers.WritingsHandler)).Methods("GET")
	api.HandleFunc("/blog/filter", s.makeHTTPHandlerFunc(handlers.BlogFilterHandler)).Methods("GET")
	api.HandleFunc("/blog/rss", s.makeHTTPHandlerFunc(handlers.BlogRSSHandler)).Methods("GET")
	api.HandleFunc("/blog/{slug}", s.makeHTTPHandlerFunc(handlers.BlogPostHandler)).Methods("GET")

	api.HandleFunc("/write", s.makeHTTPHandlerFunc(handlers.MarkdownHandler)).Methods("GET", "POST")

	// Add 404 handler
	api.NotFoundHandler = http.HandlerFunc(s.notFoundHandler)
}

func (s *Server) makeHTTPHandlerFunc(handlerFunc HTTPHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handlerFunc(w, r); err != nil {
			logger.Errorf("Handler error for %s %s: %v", r.Method, r.URL.Path, err)
			s.handleError(w, r, err)
		}
	}
}

func (s *Server) handleError(w http.ResponseWriter, r *http.Request, err error) {
	// Check if request accepts HTML or JSON
	acceptHeader := r.Header.Get("Accept")

	if contains(acceptHeader, "application/json") || contains(acceptHeader, "application/*") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"Internal Server Error","message":"%s"}`, err.Error())
	} else {
		// Return HTML error page for browsers
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Error</title></head>
<body>
	<h1>Something went wrong</h1>
	<p>Please try again later.</p>
</body>
</html>`)
	}
}

func (s *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Page Not Found</title></head>
<body>
	<h1>Page Not Found</h1>
	<p>The page you're looking for doesn't exist.</p>
	<a href="/">Go Home</a>
</body>
</html>`)
}

func (s *Server) Run() error {
	s.setupRoutes()

	// Create HTTP server with proper timeouts
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port),
		Handler:      s.router,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.config.Server.IdleTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Starting server on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(s.config.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
		return err
	}

	logger.Info("Server exited")
	return nil
}

// HTTPHandlerFunc represents a handler function that can return an error
type HTTPHandlerFunc func(w http.ResponseWriter, r *http.Request) error

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)+1] == substr+";" ||
			s[:len(substr)+1] == substr+"," || s[:len(substr)+1] == substr+" ")))
}
