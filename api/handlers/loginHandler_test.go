package handlers

import (
	"bytes"
	"fmt"
	"himakiwa/middleware"
	"himakiwa/services/database"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	database.Connect()
	db := database.DB

	mock := newUserMock(t, db)
	defer mock.close(db)

	h := middleware.CROSMiddleware(http.HandlerFunc(LoginHandler))
	ts := httptest.NewServer(h)
	defer ts.Close()

	test := []struct {
		password string
		status   int
	}{
		{
			password: mock.pass,
			status:   http.StatusOK,
		},
		{
			password: "example666",
			status:   http.StatusBadRequest,
		},
	}

	for _, tt := range test {
		buf := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, mock.email, tt.password)
		res, err := http.Post(ts.URL, "application/json; charset=UTF-8", bytes.NewBufferString(buf))
		if err != nil {
			t.Error(err)
		}
		if res.StatusCode != tt.status {
			t.Errorf("expected status '%d' but got='%d'", tt.status, res.StatusCode)
		}
	}
}
