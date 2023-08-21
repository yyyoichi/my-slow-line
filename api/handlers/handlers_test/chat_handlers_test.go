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

func TestGetChats(t *testing.T) {
	rs := services.NewRepositoryServicesMock()
	chatsHandlers := handlers.NewChatsHandlers(rs)

	req, err := http.NewRequest("GET", "/chats", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req = utils.WithUserContext(req, "1")
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/chats", chatsHandlers).Methods(http.MethodGet)
	r.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	exp := make(map[int]string)
	exp[3] = "Hello, I am 2 in the session 3"

	responseBody := rr.Body.Bytes()
	var results []handlers.GetChatsResp

	err = json.Unmarshal(responseBody, &results)
	if err != nil {
		t.Errorf("Error parsing response JSON: %s", err.Error())
		return
	}

	if len(results) != 1 {
		t.Errorf("Expected len(results) is 1, but got='%d'", len(results))
	}

	for _, lastChat := range results {
		if lastChat.Content != exp[lastChat.SessionID] {
			t.Errorf("Expected lastChat Content is '%s', but got='%s'", exp[lastChat.SessionID], lastChat.Content)
		}
	}
}

func TestGetChatAt(t *testing.T) {
	chatsAtHandlers := handlers.NewChatsAtHandlers(services.NewRepositoryServicesMock())

	test := []struct {
		endPoint  string
		expLength int
	}{
		{"/chats/2", 0},
		{"/chats/3", 2},
		{"/chats/4", 2},
	}
	for _, tt := range test {

		req, err := http.NewRequest("GET", tt.endPoint, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req = utils.WithUserContext(req, "1")
		rr := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/chats/{sessionID}", chatsAtHandlers).Methods(http.MethodGet)
		r.ServeHTTP(rr, req)

		// Check the response status code
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
		}

		responseBody := rr.Body.Bytes()
		var results []handlers.GetChatsAtResp

		err = json.Unmarshal(responseBody, &results)
		if err != nil {
			t.Errorf("Error parsing response JSON: %s", err.Error())
			return
		}

		if len(results) != tt.expLength {
			t.Errorf("Expected len(results) is '%d', but got='%d'", tt.expLength, len(results))
		}
	}
}

func TestPostChatsAt(t *testing.T) {
	rs := services.NewRepositoryServicesMock()
	chatsAtHandlers := handlers.NewChatsAtHandlers(rs)

	test := []struct {
		endPoint      string
		expStatusCode int
	}{
		{"/chats/3", http.StatusOK},
		{"/chats/2", http.StatusInternalServerError},
		{"/chats/1", http.StatusInternalServerError},
	}
	for _, tt := range test {
		body := `{"content": "Hello"}`
		req, err := http.NewRequest("POST", tt.endPoint, strings.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req = utils.WithUserContext(req, "2")
		rr := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/chats/{sessionID}", chatsAtHandlers).Methods(http.MethodPost)
		r.ServeHTTP(rr, req)

		// Check the response status code
		if rr.Code != tt.expStatusCode {
			t.Errorf("Expected status code %d, but got %d", tt.expStatusCode, rr.Code)
		}
	}
	sessionServices := rs(2).SessionServices
	chats, err := sessionServices.GetChatsAtIn48Hours(3)
	if err != nil {
		t.Error(err)
	}

	if len(chats) != 3 {
		t.Errorf("Expected len(chats) is '3', but got='%d'", len(chats))
	}
}
