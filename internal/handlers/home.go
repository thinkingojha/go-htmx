package handlers

import (
	"net/http"
	"github.com/thinkingojha/go-htmx/internal/utils"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) error {
	if err := utils.Templates.ExecuteTemplate(w, "home", nil); err != nil {
		return err
	}
	return nil
}
