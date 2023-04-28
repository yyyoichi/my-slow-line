package database

import (
	"database/sql"
	"himakiwa/auth"
	"testing"
)

func TestCreateUser(t *testing.T) {
	test := &SignInUser{
		Email:            "demo@demodemo",
		Password:         "password",
		Name:             "name",
		VerificationCode: auth.GenerateRandomSixNumber(),
	}
	id, err := test.SignIn(nil)
	if err != nil {
		t.Errorf("occur error got='%s'", err)
	}
	defer DeleteUserRow(nil, id)

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
		Email:            "demo@demodemo.com",
		Password:         "password",
		Name:             "name",
		VerificationCode: auth.GenerateRandomSixNumber(),
	}
	id, err := test.SignIn(db)
	t.Logf("\ncreate user id='%d'", id)
	if err != nil {
		t.Errorf("occur error got='%s'", err)
	}
	defer DeleteUserRow(nil, id)

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

	verificationCode := auth.GenerateRandomSixNumber()

	u := &SignInUser{
		Name:             "yaoyao",
		Email:            "email@examle.com",
		Password:         "password",
		VerificationCode: verificationCode,
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

	if du.TwoStepVerificationCode != verificationCode {
		t.Errorf("expected twostep code is %s but got='%s'", verificationCode, du.TwoStepVerificationCode)
	}
	if du.TwoVerificated {
		t.Errorf("expected deleted flag is false but got='true'")
	}
}

func TestLoginUser(t *testing.T) {
	db, err := GetDatabase()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	verificationCode := auth.GenerateRandomSixNumber()

	u := &SignInUser{
		Email:            "test@example.com",
		Name:             "testuser",
		Password:         "password",
		VerificationCode: verificationCode,
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

	verificationCode := auth.GenerateRandomSixNumber()

	u := &SignInUser{
		Email:            "test@example.com",
		Name:             "testuser",
		Password:         "password",
		VerificationCode: verificationCode,
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
			err:   sql.ErrNoRows,
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

	for i, tt := range test {
		lu := &LoginUser{
			Email:    tt.email,
			Password: tt.pass,
		}
		_, err := lu.Login(db)

		if err != tt.err {
			t.Errorf("%d expeced err is '%v' but got='%v'", i, tt.err, err)
		}
	}
}

func TestUpdateCode(t *testing.T) {
	db, err := GetDatabase()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	u := &SignInUser{
		Email:            "test@example.com",
		Name:             "testuser",
		Password:         "password",
		VerificationCode: auth.GenerateRandomSixNumber(),
	}
	userId, del := testmok(t, db, u)
	defer del()

	verificationCode := auth.GenerateRandomSixNumber()
	LogEntryStamp(db, userId, verificationCode)

	du, err := QueryUser(db, userId)

	if err != nil {
		t.Error(err)
	}

	if du.TwoStepVerificationCode != verificationCode {
		t.Errorf("expected code '%s' but got='%s'", verificationCode, du.TwoStepVerificationCode)
	}

}
