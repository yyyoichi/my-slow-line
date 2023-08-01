package handlers_test

import (
	"himakiwa/handlers"
	"himakiwa/handlers/utils"
	"himakiwa/services/sessions"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestPostParticipantsAt(t *testing.T) {
	UseSessionServicesFunc := sessions.NewSessionServicesMock()
	participantsAtHandlers := handlers.NewParticipantsAtHandlers(UseSessionServicesFunc)

	test := []struct {
		endPoint      string
		expStatusCode int
	}{
		{"/participants/3", http.StatusOK},
		{"/participants/2", http.StatusInternalServerError},
		{"/participants/1", http.StatusInternalServerError},
	}
	for _, tt := range test {
		body := `{"userID":1, "status":"rejected"}`
		req, err := http.NewRequest("PUT", tt.endPoint, strings.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req = utils.WithUserContext(req, "2")
		rr := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/participants/{sessionID}", participantsAtHandlers).Methods(http.MethodPut)
		r.ServeHTTP(rr, req)

		// Check the response status code
		if rr.Code != tt.expStatusCode {
			t.Errorf("Expected status code %d, but got %d", tt.expStatusCode, rr.Code)
		}
	}
}
