package database

import "testing"

func TestFriendRecruitment(t *testing.T) {
	usersR := &UserRepository{}

	m := createMockUser()
	tu, close := userMock(t, usersR, m)
	defer close()

	frRepository := &FRecruitmentRepository{}

	uuid := "test_uuid"
	message := "please"

	// create
	if err := frRepository.Create(tu.Id, uuid, message); err != nil {
		t.Errorf("cannot create friend-recruitment: %s", err.Error())
	}

	// get
	frs, err := frRepository.QueryByUserId(tu.Id)
	if err != nil {
		t.Error(err)
	}

	if len(frs) != 1 {
		t.Error("friend-recruitment doesnot exist")
	}

	fr := frs[0]
	if fr.Message != message {
		t.Errorf("expected message is '%s' but got='%s'", message, fr.Message)
	}
	if fr.Uuid != uuid {
		t.Errorf("expected uuid is '%s' but got='%s'", uuid, fr.Uuid)
	}

	updateMessage := "Hello"
	// update
	if err = frRepository.UpdateMessage(uuid, updateMessage); err != nil {
		t.Error(err)
	}
	// get
	frs, err = frRepository.QueryByUserId(tu.Id)
	if err != nil {
		t.Error(err)
	}

	if len(frs) != 1 {
		t.Error("friend-recruitment doesnot exist")
	}

	fr = frs[0]
	if fr.Message != updateMessage {
		t.Errorf("expected message is '%s' but got='%s'", updateMessage, fr.Message)
	}

	// delete
	if err = frRepository.DeleteAll(tu.Id); err != nil {
		t.Error(err)
	}

	frs, err = frRepository.QueryByUserId(tu.Id)
	if err != nil {
		t.Error(err)
	}

	if len(frs) != 0 {
		t.Error("cannot delete recruitment")
	}
}
