package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/thinkingojha/go-htmx/cmd/server"
	"github.com/thinkingojha/go-htmx/internal/utils"
)

func main() {

	if err := utils.ParseTemplates(filepath.Join("internal", "template")); err != nil {
		log.Fatal(fmt.Errorf(" cannot parse templates: %w", err))
	}

	s := server.NewServer(":3000")
	if err := s.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
