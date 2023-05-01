package middleware

import (
	jwttoken "himakiwa/handlers/jwt"
	"himakiwa/handlers/utils"
	"net/http"
	"os"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := utils.ReadJWTCookie(r)
		secret := os.Getenv("JWT_SECRET")
		jt := jwttoken.NewJwt(secret)
		rc, err := jt.ParseToken(reqToken)
		if err != nil {
			http.Error(w, "auth error", http.StatusBadRequest)
			return
		}
		// contextにログイン中のuserIdを格納
		r = utils.WithUserContext(r, rc.ID)
		next.ServeHTTP(w, r)
	})
}
