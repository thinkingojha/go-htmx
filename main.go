package main

import (
	"os"
	"path/filepath"

	"github.com/thinkingojha/go-htmx/cmd/server"
	"github.com/thinkingojha/go-htmx/internal/config"
	"github.com/thinkingojha/go-htmx/internal/logger"
	"github.com/thinkingojha/go-htmx/internal/utils"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.Init(cfg.App.LogLevel, cfg.IsProduction())
	logger.Infof("Starting %s v%s in %s mode", cfg.App.Name, cfg.App.Version, cfg.App.Environment)

	// Determine base directory for templates
	var baseDir string
	if cfg.IsDevelopment() {
		// In development, use current working directory
		baseDir, err = os.Getwd()
		if err != nil {
			logger.Fatalf("Failed to get working directory: %v", err)
		}
	} else {
		// In production, use executable directory
		execDir, err := os.Executable()
		if err != nil {
			logger.Fatalf("Failed to get executable path: %v", err)
		}
		baseDir = filepath.Dir(execDir)
	}

	// Parse templates
	templateDir := filepath.Join(baseDir, cfg.App.TemplateDir)
	if err := utils.ParseTemplates(templateDir); err != nil {
		logger.Fatalf("Failed to parse templates: %v", err)
	}
	logger.Infof("Templates loaded from %s", templateDir)

	// Create and run server
	srv := server.NewServer(cfg)
	if err := srv.Run(); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
