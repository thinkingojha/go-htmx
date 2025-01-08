package utils

import (
	"html/template"
	"path/filepath"
)

var Templates *template.Template

func ParseTemplates(basePath string) error {
	var err error
	Templates = template.New("t")

	if Templates, err  = Templates.ParseGlob(filepath.Join(basePath, "common", "*.html")); err != nil {
		return err
	}

	if Templates, err  = Templates.ParseGlob(filepath.Join(basePath, "home", "*.html")); err != nil {
		return err
	}

	// if Templates, err  = Templates.ParseGlob(filepath.Join(basePath, "info", "*.html")); err != nil {
	// 	return err
	// }

	// if Templates, err  = Templates.ParseGlob(filepath.Join(basePath, "blog", "*.html")); err != nil {
	// 	return err
	// }

	// if Templates, err  = Templates.ParseGlob(filepath.Join(basePath, "products", "*.html")); err != nil {
	// 	return err
	// }
	return nil
}