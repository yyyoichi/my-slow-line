package handlers_test

import (
	"himakiwa/handlers"
	"net/http"
	"testing"
)

func TestGetMe(t *testing.T) {

	// create authed user mock
	mock := NewAuthedMock(t)
	defer mock.Delete(t)

	// create server
	server := NewReqAuthHttpMock(handlers.MeHandler, mock.Jwt)
	defer server.Close()

	// request
	resp := server.Get(t)
	defer resp.BodyClose()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status-code is 200 but got='%d'", resp.StatusCode)
	}

	meResp := &handlers.GetMeResp{}
	if err := resp.BodyJson(meResp); err != nil {
		t.Error(err)
	}

	if meResp.Email != mock.TUser.Email {
		t.Errorf("expected email %s but got='%s'", mock.TUser.Email, meResp.Email)
	}
}
