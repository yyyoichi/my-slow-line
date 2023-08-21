package handlers_test

import (
	"encoding/json"
	"himakiwa/handlers"
	"himakiwa/handlers/utils"
	"himakiwa/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetMe(t *testing.T) {

	rs := services.NewRepositoryServicesMock()
	meHandlers := handlers.NewMeHandlers(rs)

	req, err := http.NewRequest("GET", "/me", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req = utils.WithUserContext(req, "1")
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/me", meHandlers.MeHandler).Methods(http.MethodGet)
	r.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	responseBody := rr.Body.Bytes()
	meResp := &handlers.GetMeResp{}
	err = json.Unmarshal(responseBody, &meResp)
	if err != nil {
		t.Errorf("Error parsing response JSON: %s", err.Error())
		return
	}
}
