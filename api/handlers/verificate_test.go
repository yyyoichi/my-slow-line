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

func TestVerfication(t *testing.T) {
	database.Connect()
	db := database.DB

	mock := newUserMock(t, db)
	defer mock.close(db)

	h := middleware.CROSMiddleware(http.HandlerFunc(VerificateTwoStepCodeHandler))
	ts := httptest.NewServer(h)
	defer ts.Close()

	test := []struct {
		code   string
		status int
	}{
		{
			code:   mock.code,
			status: http.StatusOK,
		},
		{
			code:   "000000",
			status: http.StatusBadRequest,
		},
	}
	for _, tt := range test {
		buf := fmt.Sprintf(`{"id":%d, "code":"%s"}`, mock.userId, tt.code)
		res, err := http.Post(ts.URL, "application/json; charset=UTF-8", bytes.NewBufferString(buf))
		if err != nil {
			t.Error(err)
		}
		if res.StatusCode != tt.status {
			t.Errorf("expected status '%d' but got='%d'", tt.status, res.StatusCode)
		}
	}
}
