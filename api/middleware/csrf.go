package middleware

import (
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
)

var key = []byte("xox-oxo-xox-oxo-xox-oxo-xox-oxo-")
var CSRFMiddleware = csrf.Protect(key,
	csrf.Secure(false),
	csrf.TrustedOrigins([]string{"http://127.0.0.1:3000"}),
	csrf.ErrorHandler(http.HandlerFunc(csrfError)),
)

func csrfError(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("invalid token from '%s' method '%s' to '%s'", r.Host, r.Method, r.RequestURI)
	http.Error(w, "error", http.StatusBadRequest)
}
