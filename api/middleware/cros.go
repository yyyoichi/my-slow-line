package middleware

import (
	"fmt"
	"net/http"
)

func CROSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// comment outs when cros access
		// w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		// w.Header().Set("Access-Control-Allow-Credentials", "true")
		// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Csrf-Token")
		// w.Header().Set("Access-Control-Expose-Headers", "X-Csrf-Token")
		fmt.Printf("got from '%s' method '%s' to '%s'\n", r.Host, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
