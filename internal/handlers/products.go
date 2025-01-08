package handlers

import (
	"net/http"

	"github.com/thinkingojha/go-htmx/internal/utils"
)

func ProductHandler(w http.ResponseWriter, r *http.Request) error {
	templates := utils.Templates.Templates
	if err := templates.ExecuteTemplate(w, "products", nil); err != nil {
		return err
	}
	return nil
}	