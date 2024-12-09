package server

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	err  string
	code int64
}

type HttpHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func makeHttpHandlerFunc(httpFunc HttpHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := httpFunc(w, r); err != nil {
			fmt.Fprint(w, &ApiError{
				err:  err.Error(),
				code: 500,
			})
		}
	}
}
