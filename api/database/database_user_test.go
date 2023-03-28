package database

import "testing"

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
