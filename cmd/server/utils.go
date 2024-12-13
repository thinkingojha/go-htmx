package server

import (
	"fmt"
	"log"
	"net/http"
)

type ApiError struct {
	err  string
	code int64
}

type HttpHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func makeHttpHandlerFunc(httpFunc HttpHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, "--->" , r.RequestURI,)
		if err := httpFunc(w, r); err != nil {
			fmt.Fprint(w, &ApiError{
				err:  err.Error(),
				code: 500,
			})
		}
	}
}
