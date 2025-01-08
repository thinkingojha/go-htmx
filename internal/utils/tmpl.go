package utils

import (
	"html/template"
	"path/filepath"
)

type TemplatesStruct struct {
	Templates *template.Template
	BasePath  string
}

var Templates TemplatesStruct

func ParseTemplates(basePath string) error {
	var err error
	
	templates := template.New("t")
	if templates, err = templates.ParseGlob(filepath.Join(basePath, "common", "*.html")); err != nil {
		return err
	}

	Templates = TemplatesStruct{
		Templates:templates,
		BasePath: basePath,
	}

	// if Templates, err = Templates.ParseGlob(filepath.Join(basePath, "info", "*.html")); err != nil {
	// 	return err
	// }

	// if Templates, err = Templates.ParseGlob(filepath.Join(basePath, "blog", "*.html")); err != nil {
	// 	return err
	// }

	// if Templates, err  = Templates.ParseGlob(filepath.Join(basePath, "products", "*.html")); err != nil {
	// 	return err
	// }
	return nil
}
