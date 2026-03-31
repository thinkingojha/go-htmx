package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/thinkingojha/go-htmx/internal/utils"
)

// HomePageData wraps nil data with PageName for the base template
type HomePageData struct {
	PageName string
	Title    string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) error {
	templates, err := utils.Templates.Templates.Clone()
	if err != nil {
		return err
	}
	templates, err = templates.ParseGlob(filepath.Join(utils.Templates.BasePath, "home", "*.html"))
	if err != nil {
		return err
	}
	data := HomePageData{
		PageName: "home",
		Title:    "ankush.fyi",
	}
	if err = templates.ExecuteTemplate(w, "home", data); err != nil {
		return err
	}
	return nil
}
