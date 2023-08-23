package handlers_test

import (
	"database/sql"
	"fmt"
	"himakiwa/handlers"
	"himakiwa/handlers/utils"
	"himakiwa/services"
	"himakiwa/services/database"
	"himakiwa/services/webpush"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestSessionKey(t *testing.T) {
	// create dataset
	tx := &sql.Tx{}
	useReposiotryServices := services.NewRepositoryServicesMock()
	userRepositories := useReposiotryServices(0).UserServices.GetUserRepositories()

	// create users 0~2
	users := make([]*database.TestingUser, 3)
	for i := range users {
		user, err := database.CreateTestingUser(tx, userRepositories)
		if err != nil {
			t.Error(err)
		}
		users[i] = user
	}
	user0ID := users[0].User.ID
	user1ID := users[1].User.ID
	user2ID := users[2].User.ID
	// user0 and user1 register webpush-subscription but not user2 and user3
	useReposiotryServices(user0ID).UserServices.AddWebpushSubscription("ep", "p", "at", "ua", nil)
	useReposiotryServices(user1ID).UserServices.AddWebpushSubscription("ep", "p", "at", "ua", nil)

	// session 1 to which user1 is invited is created by user0
	session1ID, _ := useReposiotryServices(user0ID).SessionServices.CreateSession("", "HOGE", user1ID)
	// session 2 to which user2 is invited is created by user0 and user2 doesnot register webpush-subscription
	session2ID, _ := useReposiotryServices(user0ID).SessionServices.CreateSession("", "HOGE", user2ID)

	test := []struct {
		expStatusCode int
		sessionID     int
		loginID       int
		inviteeID     int
	}{
		{200, session1ID, user0ID, user1ID},                   // regular
		{http.StatusBadRequest, session2ID, user0ID, user2ID}, // no subscription
		{http.StatusBadRequest, session1ID, user0ID, user2ID}, // not invited
		{http.StatusBadRequest, session1ID, user1ID, user1ID}, // not join
		{http.StatusBadRequest, 9999, user0ID, user1ID},       // not exist session
	}

	// create server
	sessionKeyHandlers := handlers.NewSessionKeyHandlers(useReposiotryServices, webpush.NewWebpushServicesMock())
	r := mux.NewRouter()
	r.HandleFunc("/sessionkey", sessionKeyHandlers).Methods(http.MethodPost)

	for i, tt := range test {
		// create request
		body := fmt.Sprintf(`{"sessionID":%d, "inviteeID":%d, "key": "key"}`, tt.sessionID, tt.inviteeID)
		req, err := http.NewRequest("POST", "/sessionkey", strings.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req = utils.WithUserContext(req, strconv.Itoa(tt.loginID))
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		// test
		if rr.Code != tt.expStatusCode {
			t.Errorf("%d: Expected Code is %d, but got='%d'", i, tt.expStatusCode, rr.Code)
		}
	}
}
