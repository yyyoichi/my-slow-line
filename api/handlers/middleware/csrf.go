package middleware

import (
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
)

var key = []byte("xox-oxo-xox-oxo-xox-oxo-xox-oxo-")
var CSRFMiddleware = csrf.Protect(key,
	csrf.SameSite(csrf.SameSiteStrictMode),
	csrf.ErrorHandler(http.HandlerFunc(csrfError)),
)

func csrfError(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("invalid token from '%s' method '%s' to '%s'", r.Host, r.Method, r.RequestURI)
	http.Error(w, "error", http.StatusBadRequest)
}
