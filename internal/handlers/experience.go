package handlers

import (
	"net/http"

	"github.com/thinkingojha/go-htmx/internal/utils"
)


func ExpHandler( w http.ResponseWriter, r *http.Request) error {
	if err := utils.Templates.ExecuteTemplate(w, "info", nil); err != nil {
		return err
	}
	return nil
}