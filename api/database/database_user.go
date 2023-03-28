package database

import (
	"database/sql"
	"errors"
	"himakiwa/auth"
	"time"
)

var (
	ErrInValidParams = errors.New("invalid params check your params")
)

type SignInUser struct {
	Name     string
	Email    string
	Password string
}

// ユーザ登録
func (u *SignInUser) SignIn(db *sql.DB) (int, error) {
	// validation
	if u.Email == "" || u.Password == "" {
		return 0, ErrInValidParams
	}

	// hashed password
	hashedPassword, err := auth.HashPassword(u.Password)
	if err != nil {
		return 0, err
	}

	// connect db
	if db == nil {
		var err error
		db, err = GetDatabase()
		if err != nil {
			return 0, err
		}
	}
	defer db.Close()

	// bigin transaction and insert db
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// exec insert
	s := `INSERT INTO users (name, email, password, create_at, login_at, update_at) VALUES(?, ?, ?, ?, ?, ?)`
	now := time.Now()
	result, err := tx.Exec(s, u.Name, u.Email, hashedPassword, now, now, now)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// get id
	id64, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	return int(id64), nil
}
