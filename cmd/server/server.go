package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/thinkingojha/go-htmx/internal/handlers"
	"net/http"
)

type Server struct {
	addrString string
}

func (s *Server) Run() error {
	fmt.Println("Server running at Port", s.addrString)
	r := mux.NewRouter()

	staticDir := "internal/static"
	// Routes
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	r.HandleFunc("/", makeHttpHandlerFunc(handlers.HomeHandler)).Methods("GET")
	r.HandleFunc("/info", makeHttpHandlerFunc(handlers.ExpHandler)).Methods("GET")
	r.HandleFunc("/products", makeHttpHandlerFunc(handlers.ProductHandler)).Methods("GET")
	r.HandleFunc("/blog", makeHttpHandlerFunc(handlers.WritingsHandler)).Methods("GET")

	if err := http.ListenAndServe(s.addrString, r); err != nil {
		return err
	}
	return nil
}

func NewServer(addString string) *Server {
	return &Server{
		addrString: addString,
	}
}
