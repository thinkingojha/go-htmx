package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/thinkingojha/go-htmx/internal/utils"
)

func WritingsHandler(w http.ResponseWriter, r *http.Request) error {
	templates, err := utils.Templates.Templates.Clone()
	if err != nil {
		return err
	}
	templates, err = templates.ParseGlob(filepath.Join(utils.Templates.BasePath, "blog", "*.html"))
	if err != nil {
		return err
	}

	if err := templates.ExecuteTemplate(w, "blog", nil); err != nil {
		return err
	}
	return nil
}
