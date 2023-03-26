package middleware

import (
	"fmt"
	"himakiwa/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func useRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(AuthMiddleware)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// read context
		id := utils.ReadUserContext(r)
		fmt.Println(id)
		w.Write([]byte(id))
	})
	http.Handle("/", r)
	return r
}

var router = useRouter()

func TestMiddleAuth(t *testing.T) {

	secret := "hogehoge"
	userId := "userABC"
	token, err := utils.NewJwt(secret).Generate(userId)
	if err != nil {
		t.Errorf("got error='%s'", err)
	}

	got := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{Name: "token", Value: token}
	req.AddCookie(cookie)

	router.ServeHTTP(got, req)

	t.Log(got.Body)
	if got.Body.String() != userId {
		t.Errorf("expected '%s' but got='%s'", userId, got.Body.String())
	}
}
