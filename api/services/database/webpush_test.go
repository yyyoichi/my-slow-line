package database

import (
	"testing"
	"time"
)

func init() {
	Connect()
}

func TestWebpush(t *testing.T) {
	usersR := &UserRepository{}

	m := createMockUser()
	tu, close := userMock(t, usersR, m)
	defer close()

	endpoint := "endpoint"
	p256dh := "testp256"
	auth := "testauth"
	expTime := time.Now()

	webpushRepository := &WebpushRepository{}
	if err := webpushRepository.Create(tu.Id, endpoint, p256dh, auth, &expTime); err != nil {
		t.Error(err)
	}

	results, err := webpushRepository.QueryByUserId(tu.Id)
	if err != nil {
		t.Error(err)
	}

	if len(results) != 1 {
		t.Errorf("expected webpush subscription length is '1' but got='%d'", len(results))
	}

	result := results[0]
	if result.Endpoint != endpoint {
		t.Errorf("expected endpoint is '%s' but got='%s'", endpoint, result.Endpoint)
	}
	if result.P256dh != p256dh {
		t.Errorf("expected p256dh is '%s' but got='%s'", p256dh, result.P256dh)
	}
	if result.Auth != auth {
		t.Errorf("expected auth is '%s' but got='%s'", auth, result.Auth)
	}
	if result.UserId != tu.Id {
		t.Errorf("expected userId is '%d' but got='%d'", tu.Id, result.UserId)
	}
	_, err = result.ExpirationTime.Value()
	if err != nil {
		t.Error(err)
	}

	if err = webpushRepository.DeleteAll(tu.Id); err != nil {
		t.Error(err)
	}

	results, err = webpushRepository.QueryByUserId(tu.Id)
	if err != nil {
		t.Error(err)
	}
	if len(results) != 0 {
		t.Errorf("expected results length after deletes is 0 but got='%d'", len(results))
	}
}
