package main

import (
	"net/http"
)

type HandlerList []http.Handler

func (hl HandlerList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, ok := w.(*response); !ok {
		w = &response{w, 0}
	}
	for _, h := range hl {
		h.ServeHTTP(w, r)
		if w.(*response).code >= 400 {
			break
		}
	}
}

// A wrapper for http.ResponseWriter to record
// if a header has been written
type response struct {
	http.ResponseWriter
	code int
}

func (r *response) WriteHeader(n int) {
	r.code = n
	r.ResponseWriter.WriteHeader(n)
}
