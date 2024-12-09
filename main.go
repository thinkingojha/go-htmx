package main

import (
	"log"

	"github.com/thinkingojha/go-htmx/cmd/server"
)

func main() {
	s := server.NewServer(":3000")

	if err := s.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
