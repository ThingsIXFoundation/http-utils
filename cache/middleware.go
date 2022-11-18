package cache

import (
	"net/http"
)

// DisableCacheOnGetRequests sets Cache-Control: no-store on GET responses by
// default. Handlers can overwrite this when required.
func DisableCacheOnGetRequests(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// set it before callback is executed allowing handlers to
			// overwrite it
			w.Header().Set("Cache-Control", "no-store")
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
