package main

import (
	"html/template"
	"log"

	"github.com/thinkingojha/go-htmx/cmd/server"
)

var templates *template.Template

func main() {
	s := server.NewServer(":3000")

	if err := s.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
