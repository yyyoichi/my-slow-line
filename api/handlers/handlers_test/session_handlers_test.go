package handlers_test

import (
	"encoding/json"
	"himakiwa/handlers"
	"himakiwa/handlers/utils"
	"himakiwa/services"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetSessions(t *testing.T) {
	sessionsHandlers := handlers.NewSessionsHandlers(services.NewRepositoryServicesMock())

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
	rs := services.NewRepositoryServicesMock()
	sessionsHandlers := handlers.NewSessionsHandlers(rs)

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
	}
}

func TestGetSessionAt(t *testing.T) {
	sessionAtHandlers := handlers.NewSessionAtHandlers(services.NewRepositoryServicesMock())

	// Create a new test request with a mock user context
	req, err := http.NewRequest("GET", "/sessions/3", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	// Set a dummy user context for testing
	req = utils.WithUserContext(req, "1")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/sessions/{sessionID}", sessionAtHandlers).Methods(http.MethodGet)
	r.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}
	responseBody := rr.Body.Bytes()
	var result handlers.SessionAtResp

	// Unmarshal the response body into the GetSessionsResp struct
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		t.Errorf("Error parsing response JSON: %s", err.Error())
		return
	}

	if result.Name != "Session3" {
		t.Errorf("Excepted Name is 'Session3', but got='%s'", result.Name)
	}

	if len(result.Participants) != 2 {
		t.Errorf("Expected len(result.Participants) is 2, but got='%d'", len(result.Participants))
	}
}

func TestPostSessionAt(t *testing.T) {
	test := []struct {
		userID        string
		expStatusCode int
	}{
		{"1", http.StatusOK},
		{"2", http.StatusInternalServerError},
	}

	for _, tt := range test {
		sessionAtHandlers := handlers.NewSessionAtHandlers(services.NewRepositoryServicesMock())

		// create Request
		body := `{"sessionName": "New Session"}`
		req, err := http.NewRequest("PUT", "/sessions/1", strings.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req = utils.WithUserContext(req, tt.userID)
		rr := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/sessions/{sessionID}", sessionAtHandlers).Methods(http.MethodPut)
		r.ServeHTTP(rr, req)

		// Check the response status code
		if rr.Code != tt.expStatusCode {
			t.Errorf("Expected status code %d, but got %d", tt.expStatusCode, rr.Code)
		}
	}
}
