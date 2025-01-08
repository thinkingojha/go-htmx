package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/thinkingojha/go-htmx/cmd/server"
	"github.com/thinkingojha/go-htmx/internal/utils"
)

func main() {
	execDir, err := os.Executable()
	if err != nil {
		log.Fatal(fmt.Errorf(" something went wrong: %w", err))
	}
	baseDir, _ := filepath.Split(filepath.Dir(execDir))
	if err := utils.ParseTemplates(filepath.Join(baseDir, "internal", "template")); err != nil {
		log.Fatal(fmt.Errorf(" cannot parse templates: %w", err))
	}

	s := server.NewServer(":3000")
	if err := s.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
