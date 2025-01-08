package handlers

import (
	"net/http"

	"github.com/thinkingojha/go-htmx/internal/utils"
)

func WritingsHandler(w http.ResponseWriter, r *http.Request) error {
	if err := utils.Templates.ExecuteTemplate(w, "blog", nil); err != nil {
		return err
	}
	return nil
}