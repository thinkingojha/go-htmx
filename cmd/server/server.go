package server

import (
	"fmt"
	"net/http"
	"github.com/thinkingojha/go-htmx/internal/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	addrString string
}

func (s *Server) Run() error {
	fmt.Println("Server running at Port", s.addrString)
	r := mux.NewRouter()
	
	// Routes
	r.HandleFunc("/", makeHttpHandlerFunc( handlers.HomeHandler))

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
