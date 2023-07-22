package handlers_test

import (
	"himakiwa/handlers"
	"himakiwa/handlers/utils"
	"himakiwa/services/database"
	"himakiwa/services/users"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockType = struct {
	userID  int
	uuid    string
	message string
}

func NewTestRecruitHandlers(init []mockType) handlers.RecruitHandlers {

	// add init data
	repo := database.NewMockFRecruitmentRepository()
	for _, d := range init {
		repo.Create(d.userID, d.uuid, d.message)
	}

	return handlers.RecruitHandlers{GetFriendRecruitService: func(userID int) users.FriendRecruitService {
		return users.FriendRecruitService{UserId: userID, FRecruitmentRepo: repo}
	}}
}

func TestPostRecruitmentsHandler(t *testing.T) {
	// Create the test instance of RecruitHandlers
	recruitHandlers := NewTestRecruitHandlers([]mockType{})

	// Prepare a test request with a valid PostRecruitmentBody
	body := `{"message": "Test message"}`
	req, err := http.NewRequest("POST", "/recruitments", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	// Set a dummy user context for testing
	req = utils.WithUserContext(req, "123") // Assuming user ID is "123"

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the postRecruitmentsHandler
	http.HandlerFunc(recruitHandlers.PostRecruitmentsHandler).ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check if the recruitment was created in the mock repository (assuming Create method in the service returns nil on success)
	// Assuming you have implemented the MockFRecruitmentRepository correctly.
	// You might need to update this part based on your actual implementation.
	if mockRepo, ok := recruitHandlers.GetFriendRecruitService(123).FRecruitmentRepo.(*database.MockFRecruitmentRepository); ok {
		if len(mockRepo.RecruitmentByID[123]) != 1 {
			t.Errorf("Expected 1 recruitment for user ID 123, but got %d", len(mockRepo.RecruitmentByID[123]))
		}
	} else {
		t.Error("Expected a MockFRecruitmentRepository, but got something else")
	}
}

func TestPutRecruitmentsHandler(t *testing.T) {
	// Create the test instance of RecruitHandlers with a mock service
	recruitHandlers := NewTestRecruitHandlers([]mockType{
		{123, "12345", ""},
	})

	// Prepare a test request with a valid PutRecruitmentBody JSON in the request body
	body := `{"uuid": "12345", "message": "Test message", "deleted": true}`
	req, err := http.NewRequest("PUT", "/recruitments", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	// Set a dummy user context for testing
	req = utils.WithUserContext(req, "123") // Assuming user ID is "123"

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the PutRecruitmentsHandler
	http.HandlerFunc(recruitHandlers.PutRecruitmentsHandler).ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check if the UpdateAt method was called on the mock service
	if mockRepo, ok := recruitHandlers.GetFriendRecruitService(123).FRecruitmentRepo.(*database.MockFRecruitmentRepository); ok {
		if len(mockRepo.RecruitmentByID[123]) != 1 {
			t.Errorf("Expected 1 recruitment for user ID 123, but got %d", len(mockRepo.RecruitmentByID[123]))
		}
		r := mockRepo.RecruitmentByID[123][0]
		if r.Message != "Test message" {
			t.Errorf("Expected '%s' message, but got %s", "Test message", r.Message)
		}
	} else {
		t.Error("Expected a MockFRecruitmentRepository, but got something else")
	}

}
