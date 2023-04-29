package database

import (
	"testing"
	"time"
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

func deepEqual(t *testing.T, act, exp *TUser) {
	if act.Id != exp.Id {
		t.Errorf("expected Id is %d but got='%d'", exp.Id, act.Id)
	}
	if act.Name != exp.Name {
		t.Errorf("expected Name is %s but got='%s'", exp.Name, act.Name)
	}
	if act.Email != exp.Email {
		t.Errorf("expected Email is %s but got='%s'", exp.Email, act.Email)
	}
	if act.HashedPass != exp.HashedPass {
		t.Errorf("expected HashedPass is %s but got='%s'", exp.HashedPass, act.HashedPass)
	}
	if act.VCode != exp.VCode {
		t.Errorf("expected VCode is %s but got='%s'", exp.VCode, act.VCode)
	}
	if act.Deleted != exp.Deleted {
		t.Errorf("expected Deleted is %v but got='%v'", exp.Deleted, act.Deleted)
	}
	if act.TwoVerificated != exp.TwoVerificated {
		t.Errorf("expected TwoVerificated is %v but got='%v'", exp.TwoVerificated, act.TwoVerificated)
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

	deepEqual(t, qu, tu)

	qu, err = usersR.QueryByEMail(tu.Email)
	if err != nil {
		t.Error(err)
	}
	testNotNil(t, qu)

	deepEqual(t, qu, tu)
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

	if err := usersR.UpdateLoginTimeAndResetVCode(tu.Id, newCode, time.Now()); err != nil {
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

	if err := usersR.UpdateVerifiscatedAt(tu.Id, time.Now()); err != nil {
		t.Error(err)
	}

	qu, err = usersR.QueryById(tu.Id)
	if err != nil {
		t.Error(err)
	}

	if !qu.TwoVerificated {
		t.Errorf("expected is true but got='false'")
	}

}
