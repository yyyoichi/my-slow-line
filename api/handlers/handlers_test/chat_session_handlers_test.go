package handlers_test

import (
	"encoding/json"
	"himakiwa/handlers"
	"himakiwa/handlers/utils"
	chat_services "himakiwa/services/chats"
	"himakiwa/services/database"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestChatSessionsHandlers() handlers.ChatSessionsHandlers {
	csRepo := database.NewMockChatSessionRepository()
	cspRepo := database.NewMockChatSessionParticipantRepository()
	cRepo := database.NewMockChatRepository()
	return handlers.ChatSessionsHandlers{GetChatService: func(loginUserID int) *chat_services.ChatService {
		return &chat_services.ChatService{
			ChatSessionRepo:            csRepo,
			ChatSessionParticipantRepo: cspRepo,
			ChatRepo:                   cRepo,
			UserID:                     loginUserID,
		}
	}}
}

func TestGetChatSessionsHandler(t *testing.T) {
	// Create the test instance
	chatSessionsHandlers := newTestChatSessionsHandlers()

	// add data in mock
	services := chatSessionsHandlers.GetChatService(123)
	// add session
	err := services.CreateChatSession("", "Test Session", []int{})
	if err != nil {
		t.Error(err)
	}
	// add chat
	services.SendMessage(1, "Test Hello")

	// Create a new test request with a mock user context
	req, err := http.NewRequest("GET", "/chatsessions", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set a dummy user context for testing
	req = utils.WithUserContext(req, "123") // Assuming user ID is "123"

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the postRecruitmentsHandler
	http.HandlerFunc(chatSessionsHandlers.GetChatSessionsHandler).ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Get the response body as a byte slice
	responseBody := rr.Body.Bytes()

	// Create a new instance of GetSessionsHandlerResp to store the parsed data
	var results []handlers.GetSessionsHandlerResp

	// Unmarshal the response body into the GetSessionsHandlerResp struct
	err = json.Unmarshal(responseBody, &results)
	if err != nil {
		t.Errorf("Error parsing response JSON: %s", err.Error())
		return
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 Sessions, but got %d Sessions", len(results))
	}

	result := results[0]

	// Now you can test the values of the parsed data
	if result.ID != 1 {
		t.Errorf("Expected ID: 1, but got: %d", result.ID)
	}

	if result.Name != "Test Session" {
		t.Errorf("Expected Name: 'Test Session', but got: %s", result.Name)
	}

	if result.Status != database.Joined {
		t.Errorf("Expected Status: %s, but got: %s", database.Joined, result.Status)
	}

	if len(result.Chats) != 1 {
		t.Errorf("Expected 1 Chats, but got %d Chats", len(result.Chats))
	}

	if result.Chats[0].Content != "Test Hello" {
		t.Errorf("Expected Status: Test Hello, but got: %s", result.Chats[0].Content)
	}
}
