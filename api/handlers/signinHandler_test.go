package handlers

import (
	"bytes"
	"fmt"
	"himakiwa/database"
	"himakiwa/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignin(t *testing.T) {
	h := middleware.CROSMiddleware(http.HandlerFunc(SigninHandler))
	ts := httptest.NewServer(h)
	defer ts.Close()
	email := "example@example.com"
	defer database.DeleteUserRowWithEmail(nil, email)
	buf := fmt.Sprintf(`{"name":"%s", "email":"%s", "password":"%s"}`, "user1", email, "passw11ordd")
	res, err := http.Post(ts.URL, "application/json; charset=UTF-8", bytes.NewBufferString(buf))
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 but got='%d'", res.StatusCode)
	}
}
