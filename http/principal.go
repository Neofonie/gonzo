package http

import "net/http"

/** Handle principal calls */
type Principal struct {
	Next func(http.ResponseWriter, *http.Request, string)
}

func (p Principal) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if pr := r.Header.Get("X-Principal"); pr != "" {
		p.Next(w, r, pr)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
