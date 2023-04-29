package users

import (
	"himakiwa/services/database"
	"testing"
)

type mock struct {
	name  string
	email string
	pass  string
}

func testMock(t *testing.T, u *UsersService, m *mock) (*database.TUser, func()) {
	tu, err := u.Signin(m.email, m.pass, m.name)
	if err != nil {
		t.Error(err)
	}
	return tu, func() {
		if err = u.HardDelete(tu.Id); err != nil {
			t.Error(err)
		}
	}
}

func TestQuery(t *testing.T) {
	m := &mock{"demo", "test@sample.com", "pa55word"}

	db, err := database.GetDatabase()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	u := NewUsersService(db)
	tu, close := testMock(t, u, m)
	defer close()

	qu, err := u.Query(tu.Id)
	if err != nil {
		t.Error(err)
	}

	if tu.Id != qu.Id {
		t.Errorf("expected Id is %d but got='%d'", tu.Id, qu.Id)
	}
}

func TestLogin(t *testing.T) {
	m := &mock{"demo", "test@sample.com", "pa55word"}

	db, err := database.GetDatabase()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	u := NewUsersService(db)
	tu, close := testMock(t, u, m)
	defer close()

	test := []struct {
		email string
		pass  string
		err   error
	}{
		{
			email: m.email,
			pass:  m.pass,
			err:   nil,
		},
		{
			email: "",
			pass:  m.pass,
			err:   ErrInValidParams,
		},
		{
			email: m.email,
			pass:  "",
			err:   ErrInValidParams,
		},
		{
			email: "demo@demo",
			pass:  m.pass,
			err:   ErrInvalidEmail,
		},
		{
			email: m.email,
			pass:  "buzzword",
			err:   ErrInvalidPassword,
		},
		{
			email: m.email,
			pass:  "ppp",
			err:   ErrInvalidPassword,
		},
	}

	for i, tt := range test {
		qu, err := u.Login(tt.email, tt.pass)
		if err != tt.err {
			t.Errorf("%d: expected err is '%s' but got='%s'", i, tt.err, err)
		}
		if err != nil {
			continue
		}
		// regular
		if tu.VCode == qu.VCode {
			t.Errorf("%d: expected difference vcode but equal", i)
		}

		if tu.LoginAt == qu.LoginAt {
			t.Errorf("%d: expected difference time of LoginAt but equal", i)
		}
	}
}

func TestVerificate(t *testing.T) {
	m := &mock{"demo", "test@sample.com", "pa55word"}

	db, err := database.GetDatabase()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	u := NewUsersService(db)
	tu, close := testMock(t, u, m)
	defer close()

	test := []struct {
		vcode string
		err   error
	}{
		{
			vcode: tu.VCode,
			err:   nil,
		},
		{
			vcode: "",
			err:   ErrInValidParams,
		},
		{
			vcode: "12345",
			err:   ErrInValidParams,
		},
		{
			vcode: "123456",
			err:   ErrInvalidVCode,
		},
	}
	for i, tt := range test {
		qu, err := u.Verificate(tu.Id, tt.vcode)
		if err != tt.err {
			t.Errorf("%d: expected err is '%s' but got='%s'", i, tt.err, err)
		}
		if err != nil {
			continue
		}
		// regular
		if !qu.TwoVerificated {
			t.Errorf("%d: expected TwoVerificated is true but got='false'", i)
		}

		if !qu.TwoVerificatedAt.Valid {
			t.Errorf("%d: expected TwoVerificatedAt is valid but invalid", i)
		}
	}
}
