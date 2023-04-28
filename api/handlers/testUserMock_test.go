package handlers

import (
	"database/sql"
	"himakiwa/auth"
	"himakiwa/services/database"
	"os"
	"testing"
)

type userMock struct {
	userId int
	code   string
	email  string
	pass   string
	name   string
}

func (u *userMock) close(db *sql.DB) {
	database.DeleteUserRow(db, u.userId)
}

func newUserMock(t *testing.T, db *sql.DB) *userMock {
	verificationCode := auth.GenerateRandomSixNumber()

	u := &database.SignInUser{
		Email:            os.Getenv("EMAIL_ADDRESS"),
		Name:             "testuser",
		Password:         "passw0rd",
		VerificationCode: verificationCode,
	}

	result, err := database.ExistEmail(db, u.Email)
	if err != nil {
		t.Error(err)
	}
	if result {
		if err = database.DeleteUserRowWithEmail(db, u.Email); err != nil {
			t.Error(err)
		}
	}

	userId, err := u.SignIn(db)
	if err != nil {
		t.Error(err)
	}
	if userId == 0 {
		t.Errorf("userId is 0")
	}

	return &userMock{userId, verificationCode, u.Email, u.Password, u.Name}
}
