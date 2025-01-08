package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/thinkingojha/go-htmx/internal/utils"
)

func ExpHandler(w http.ResponseWriter, r *http.Request) error {
	templates, err := utils.Templates.Templates.Clone()
	if err != nil {
		return err
	}
	templates, err = templates.ParseGlob(filepath.Join(utils.Templates.BasePath, "info", "*.html"))
	if err != nil {
		return err
	}
	if err = templates.ExecuteTemplate(w, "info", nil); err != nil {
		return err
	}
	return nil
}
