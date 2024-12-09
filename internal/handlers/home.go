package handlers

import (
	"fmt"
	"net/http"
)

func Home (w http.ResponseWriter, r *http.Request) error {
	if _, err := fmt.Fprintf(w, "Another attempt at making basic stuff in golang"); err != nil {
		return err
	}
	return nil
}