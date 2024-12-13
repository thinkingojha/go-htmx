package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) error {

	_, err := filepath.Glob("internal/template/home/*.html")
	if err != nil{
		log.Fatal(err.Error())
		return err
	}

	tmpl := template.Must(template.ParseGlob("internal/template/home/*.html"))
	if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		return err
	}
	return nil
}
