package database

import (
	"reflect"
	"testing"
)

type mock struct {
	name  string
	email string
	pass  string
	vcode string
}

func createMockUser() *mock {
	return &mock{"sample", "test@example.com", "pa55word", "098123"}
}

func userMock(t *testing.T, usersR *UserRepository, m *mock) (*TUser, func()) {
	tu, err := usersR.Create(m.name, m.email, m.pass, m.vcode)
	if err != nil {
		t.Error(err)
	}
	return tu, func() {
		err = usersR.HardDeleteById(tu.Id)
		if err != nil {
			t.Error(err)
		}
	}
}

func testNotNil(t *testing.T, u *TUser) {
	if u.Id == 0 {
		t.Error("Id got='0'")
	}
	if u.Name == "" {
		t.Error("Name got=''")
	}
	if u.Email == "" {
		t.Error("Email got=''")
	}
	if u.HashedPass == "" {
		t.Error("Name got=''")
	}
	if u.VCode == "" {
		t.Error("VCode got=''")
	}

	if u.LoginAt.IsZero() {
		t.Error("LoginAt is zero")
	}
	if u.UpdateAt.IsZero() {
		t.Error("LoginAt is zero")
	}
	if u.CreateAt.IsZero() {
		t.Error("LoginAt is zero")
	}
}
func TestQuery(t *testing.T) {
	db, err := GetDatabase()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	usersR := GetUsers(db)

	m := createMockUser()
	tu, close := userMock(t, usersR, m)
	defer close()

	testNotNil(t, tu)

	qu, err := usersR.QueryById(tu.Id)
	if err != nil {
		t.Error(err)
	}
	testNotNil(t, qu)

	if reflect.DeepEqual(tu, qu) {
		t.Errorf("expected is %v but got='%v' ", tu, qu)
	}

	qu, err = usersR.QueryByEMail(tu.Email)
	if err != nil {
		t.Error(err)
	}
	testNotNil(t, qu)

	if reflect.DeepEqual(tu, qu) {
		t.Errorf("expected is %v but got='%v' ", tu, qu)
	}
}

func TestSoftDelete(t *testing.T) {
	db, err := GetDatabase()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	usersR := GetUsers(db)

	m := createMockUser()
	tu, close := userMock(t, usersR, m)
	defer close()
	testNotNil(t, tu)

	if err = usersR.SoftDeleteById(tu.Id); err != nil {
		t.Error(err)
	}

	qu, err := usersR.QueryById(tu.Id)
	if err != nil {
		t.Error(err)
	}
	testNotNil(t, qu)

	if !qu.Deleted {
		t.Error("expected deleted flag on but 'off'")
	}

	if err = usersR.ActivateById(tu.Id); err != nil {
		t.Error(err)
	}

	qu, err = usersR.QueryById(tu.Id)
	if err != nil {
		t.Error(err)
	}
	testNotNil(t, qu)

	if qu.Deleted {
		t.Error("expected deleted flag off but 'on'")
	}
}

func TestUpdateCode(t *testing.T) {
	db, err := GetDatabase()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	usersR := GetUsers(db)

	m := createMockUser()
	tu, close := userMock(t, usersR, m)
	defer close()

	newCode := "123456"

	if err := usersR.UpdateLoginTimeAndResetVCode(tu.Id, newCode); err != nil {
		t.Error(err)
	}

	qu, err := usersR.QueryById(tu.Id)
	if err != nil {
		t.Error(err)
	}

	if qu.VCode != newCode {
		t.Errorf("expected is %s but got='%s'", newCode, qu.VCode)
	}

	if qu.LoginAt == tu.LoginAt {
		t.Errorf("expected is %v but got='%v'", tu.LoginAt, qu.LoginAt)
	}

	if qu.TwoVerificated {
		t.Errorf("expected is false but got='true'")
	}

}
