package handlers_test

import (
	"encoding/json"
	"himakiwa/handlers"
	"himakiwa/handlers/utils"
	"himakiwa/services/database"
	"himakiwa/services/sessions"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetSessions(t *testing.T) {
	UseSessionServicesFunc := sessions.NewSessionServicesMock()
	sessionsHandlers := handlers.NewSessionsHandlers(UseSessionServicesFunc)

	// Create a new test request with a mock user context
	req, err := http.NewRequest("GET", "/sessions", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	test := []struct {
		userID       string
		expResultLen int
	}{
		{"1", 4},
		{"2", 3},
	}

	for _, tt := range test {

		// Set a dummy user context for testing
		req = utils.WithUserContext(req, tt.userID)

		// Create a response recorder to record the response
		rr := httptest.NewRecorder()

		// Call the postRecruitmentsHandler
		http.HandlerFunc(sessionsHandlers).ServeHTTP(rr, req)

		// Check the response status code
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
		}

		// Get the response body as a byte slice
		responseBody := rr.Body.Bytes()

		// Create a new instance of GetSessionsResp to store the parsed data
		var results []handlers.GetSessionsResp

		// Unmarshal the response body into the GetSessionsResp struct
		err = json.Unmarshal(responseBody, &results)
		if err != nil {
			t.Errorf("Error parsing response JSON: %s", err.Error())
			return
		}

		if len(results) != tt.expResultLen {
			t.Errorf("Expected len(results) is '%d', but got='%d'", tt.expResultLen, len(results))
		}
	}
}

func TestPostSessions(t *testing.T) {
	UseSessionServicesFunc := sessions.NewSessionServicesMock()
	sessionsHandlers := handlers.NewSessionsHandlers(UseSessionServicesFunc)

	test := []struct {
		body          string
		expStatusCode int
		expName       string
	}{
		{`{"recruitUUID": "Test UUID of userID2", "sessionName": "Test Session", "publicKey": "AA"}`, http.StatusOK, "Test Session"},
		{`{"recruitUUID": "Test", "sessionName": "Test Session", "publicKey": "AA"}`, http.StatusInternalServerError, ""},
	}

	for _, tt := range test {

		// Create a new test request
		req, err := http.NewRequest("POST", "/sessions", strings.NewReader(tt.body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		// Set a dummy user context for testing
		req = utils.WithUserContext(req, "1")
		rr := httptest.NewRecorder()
		http.HandlerFunc(sessionsHandlers).ServeHTTP(rr, req)

		// Check the response status code
		if rr.Code != tt.expStatusCode {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
		}
		if tt.expStatusCode != http.StatusOK {
			continue
		}

		sessionServices := UseSessionServicesFunc(1)
		querySessions, err := sessionServices.GetActiveOrArchivedSessions()
		if err != nil {
			t.Error(err)
		}
		var session *database.TQuerySessions
		for _, ss := range querySessions {
			if ss.PublicKey == "AA" {
				session = ss
			}
		}

		if session.Name != tt.expName {
			t.Errorf("Expected session.Name is '%s', but got='%s'", tt.expName, session.Name)
		}
	}
}
