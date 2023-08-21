package handlers_test

import (
	"fmt"
	"himakiwa/handlers"
	"himakiwa/handlers/utils"
	"himakiwa/services"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostRecruitmentsHandler(t *testing.T) {
	// Create the test instance of RecruitHandlers
	recruitHandlers := handlers.NewRecruitHandlers(services.NewRepositoryServicesMock())

	// Prepare a test request with a valid PostRecruitmentBody
	body := `{"message": "Test message"}`
	req, err := http.NewRequest("POST", "/recruitments", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	// Set a dummy user context for testing
	req = utils.WithUserContext(req, "1") // Assuming user ID is "1"
	rr := httptest.NewRecorder()

	http.HandlerFunc(recruitHandlers).ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestPutRecruitmentsHandler(t *testing.T) {
	// Create the test instance of RecruitHandlers
	rs := services.NewRepositoryServicesMock()
	recruits, err := rs(1).UserServices.GetRecruitments()
	if err != nil {
		t.Error(err)
	}
	uuid := recruits[0].UUID
	recruitHandlers := handlers.NewRecruitHandlers(rs)

	// Prepare a test request with a valid PutRecruitmentBody JSON in the request body
	body := fmt.Sprintf(`{"uuid": "%s", "message": "Test message", "deleted": true}`, uuid)
	req, err := http.NewRequest("PUT", "/recruitments", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	// Set a dummy user context for testing
	req = utils.WithUserContext(req, "1") // Assuming user ID is "123"

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the PutRecruitmentsHandler
	http.HandlerFunc(recruitHandlers).ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

}
