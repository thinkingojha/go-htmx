package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) error {

	execDir, err := os.Executable()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	
	baseDir := filepath.Dir(execDir)
	templatePath := filepath.Join(baseDir, "../internal/template/home/*.html") 
	_, err = filepath.Glob(templatePath)
	if err != nil{
		log.Fatal(err.Error())
		return err
	}

	tmpl := template.Must(template.ParseGlob(templatePath))
	if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		return err
	}
	return nil
}
