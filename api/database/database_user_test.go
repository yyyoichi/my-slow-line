package database

import (
	"database/sql"
	"testing"
)

func TestCreateUser(t *testing.T) {
	test := &SignInUser{
		Email:    "demo@demodemo",
		Password: "password",
		Name:     "name",
	}
	id, err := test.SignIn(nil)
	if err != nil {
		t.Errorf("occur error got='%s'", err)
	}

	if id == 0 {
		t.Error("id expected int64 but got nil")
	} else {
		t.Logf("\ninserted id='%d'", id)
	}

}

func TestExistUser(t *testing.T) {

	db, err := GetDatabase()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	test := &SignInUser{
		Email:    "demo@demodemo.com",
		Password: "password",
		Name:     "name",
	}
	id, err := test.SignIn(db)
	t.Logf("\ncreate user id='%d'", id)
	if err != nil {
		t.Errorf("occur error got='%s'", err)
	}

	result, err := ExistEmail(db, test.Email)
	if err != nil {
		t.Errorf("occur error got='%s'", err)
	}

	if !result {
		t.Errorf("expect false bc exist user but got true")
	}

	err = DeleteUserRow(db, id)
	if err != nil {
		t.Errorf("delete error '%s'", err)
	}
	t.Logf("\ndelete user id='%d'", id)
}

func testmok(t *testing.T, db *sql.DB, u *SignInUser) (int, func()) {
	id, err := u.SignIn(db)
	if err != nil {
		t.Error(err)
	}
	return id, func() {
		DeleteUserRow(db, id)
	}
}

func TestQueryUser(t *testing.T) {
	db, err := GetDatabase()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	u := &SignInUser{
		Name:     "yaoyao",
		Email:    "email@examle.com",
		Password: "password",
	}
	userId, del := testmok(t, db, u)
	defer del()

	du, err := QueryUser(db, userId)

	if err != nil {
		t.Error(err)
	}

	if du.Name != u.Name {
		t.Errorf("expected name is %s but got='%s'", u.Name, du.Name)
	}
	if du.Email != u.Email {
		t.Errorf("expected email is %s but got='%s'", u.Email, du.Email)
	}
	if du.Deleted {
		t.Errorf("expected deleted flag is false but got='true'")
	}
	if du.CreateAt.IsZero() {
		t.Errorf("expected create_at but it is zero")
	}
	if du.UpdateAt.IsZero() {
		t.Errorf("expected update_at but it is zero")
	}
	if du.LoginAt.IsZero() {
		t.Errorf("expected login_at but it is zero")
	}
}

func TestLoginUser(t *testing.T) {
	db, err := GetDatabase()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	u := &SignInUser{
		Email:    "test@example.com",
		Name:     "testuser",
		Password: "password",
	}
	userId, del := testmok(t, db, u)
	defer del()

	lu := &LoginUser{
		Email:    u.Email,
		Password: u.Password,
	}

	du, err := lu.Login(db)
	if err != nil {
		t.Errorf("cannot login got='%s'", err)
	}

	if du.Id != userId {
		t.Errorf("expected userid '%d' but got='%d'", userId, du.Id)
	}
}

func TestInvalidLogin(t *testing.T) {
	db, err := GetDatabase()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	u := &SignInUser{
		Email:    "test@example.com",
		Name:     "testuser",
		Password: "password",
	}
	_, del := testmok(t, db, u)
	defer del()

	test := []struct {
		email string
		pass  string
		err   error
	}{
		{
			email: "test@example.com",
			pass:  "passww",
			err:   ErrInvalidPassword,
		},
		{
			email: "testtest@example.com",
			pass:  "password",
			err:   ErrInvalidEmail,
		},
		{
			email: "",
			pass:  "password",
			err:   ErrInValidParams,
		},
		{
			email: "test@example.com",
			pass:  "",
			err:   ErrInValidParams,
		},
	}

	for _, tt := range test {
		lu := &LoginUser{
			Email:    tt.email,
			Password: tt.pass,
		}
		_, err := lu.Login(db)

		if err != tt.err {
			t.Error("expeced err is '%w' but got='%w'", tt.err, err)
		}
	}
}
