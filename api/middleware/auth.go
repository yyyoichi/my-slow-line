package middleware

import (
	"himakiwa/auth"
	"himakiwa/utils"
	"net/http"
	"os"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := utils.ReadJWTCookie(r)
		secret := os.Getenv("JWT_SECRET")
		jt := auth.NewJwt(secret)
		rc, err := jt.ParseToken(token)
		if err != nil {
			http.Error(w, "auth error", http.StatusBadRequest)
			return
		}
		// contextにログイン中のuserIdを格納
		r = utils.WithUserContext(r, rc.ID)
		next.ServeHTTP(w, r)
	})
}
