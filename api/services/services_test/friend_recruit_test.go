package services_test

import (
	"himakiwa/services"
	"himakiwa/services/users"
	"testing"
)

func TestFriendRecruit(t *testing.T) {
	m := &mock{"demo", "test@sample.com", "pa55word"}
	u := services.NewRepositoryServices().GetUser()
	tu, close := testMock(t, u, m)
	defer close()

	fr := u.GetFriendRecruitService(tu.Id)
	err := fr.Create("")
	if err != nil {
		t.Error(err)
	}
	defer fr.DeleteHard()

	recruits, err := fr.Query()
	if err != nil {
		t.Error(err)
	}
	uuid := recruits[0].Uuid

	test := []struct {
		userId int
		err    error
	}{
		{
			userId: tu.Id,
			err:    nil,
		},
		{
			userId: tu.Id + 1,
			err:    users.ErrInvalidUuid,
		},
	}
	for i, tt := range test {
		service := u.GetFriendRecruitService(tt.userId)
		if err = service.UpdateMessageAt(uuid, ""); err != tt.err {
			t.Errorf("%d: expcted error is '%s' but got='%s'", i, tt.err, err)
		}
	}

	user, err := u.QueryByRecruitUuid(uuid)
	if err != nil {
		t.Error(err)
	}

	if user.Id != tu.Id {
		t.Errorf("expected userId is '%d' but got='%d'", tu.Id, user.Id)
	}
}
