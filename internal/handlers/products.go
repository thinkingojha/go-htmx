package handlers

import (
	"net/http"

	"github.com/thinkingojha/go-htmx/internal/utils"
)

func ProductHandler(w http.ResponseWriter, r *http.Request) error {
	if err := utils.Templates.ExecuteTemplate(w, "products", nil); err != nil {
		return err
	}
	return nil
}	