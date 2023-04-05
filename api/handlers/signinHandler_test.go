package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"himakiwa/database"
	"himakiwa/middleware"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSignin(t *testing.T) {
	database.Connect()
	db := database.DB
	// create server
	h := middleware.CROSMiddleware(http.HandlerFunc(SigninHandler))
	ts := httptest.NewServer(h)
	defer ts.Close()

	// create req body
	email := os.Getenv("EMAIL_ADDRESS")
	buf := fmt.Sprintf(`{"name":"%s", "email":"%s", "password":"%s"}`, "user1", email, "passw11ordd")

	// delete db row
	defer database.DeleteUserRowWithEmail(db, email)

	// request
	res, err := http.Post(ts.URL, "application/json; charset=UTF-8", bytes.NewBufferString(buf))
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 but got='%d'", res.StatusCode)
	}
	defer res.Body.Close()

	// parse body
	var userId int
	err = json.NewDecoder(res.Body).Decode(&userId)
	if err != nil {
		t.Fatalf("failed to decode HTTP response: %v", err)
	}
	if userId == 0 {
		t.Errorf("userId got='%d'", userId)
	}
}
