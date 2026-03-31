package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/thinkingojha/go-htmx/internal/utils"
)

type ProductsPageData struct {
	PageName string
	Title    string
}

func ProductHandler(w http.ResponseWriter, r *http.Request) error {
	templates, err := utils.Templates.Templates.Clone()
	if err != nil {
		return err
	}
	templates, err = templates.ParseGlob(filepath.Join(utils.Templates.BasePath, "products", "*.html"))
	if err != nil {
		return err
	}
	data := ProductsPageData{
		PageName: "products",
		Title:    "products",
	}
	if err := templates.ExecuteTemplate(w, "products", data); err != nil {
		return err
	}
	return nil
}