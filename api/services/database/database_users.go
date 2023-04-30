package database

import (
	"database/sql"
	"time"
)

type TUser struct {
	Id               int
	Name             string
	HashedPass       string
	Email            string
	LoginAt          time.Time
	CreateAt         time.Time
	UpdateAt         time.Time
	Deleted          bool
	VCode            string
	TwoVerificatedAt sql.NullTime
	TwoVerificated   bool
	// db get
	dbDeleted        int64
	dbTwoVerificated int64
}

type UserRepository struct{}

func (u *UserRepository) Create(name, email, hashedPass, vcode string) (*TUser, error) {
	// exec insert
	s := `INSERT INTO users (name, email, password, create_at, login_at, update_at, two_step_verification_code) VALUES(?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	result, err := DB.Exec(s, name, email, hashedPass, now, now, now, vcode)
	if err != nil {
		return nil, err
	}

	// get id
	id64, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	tu := &TUser{
		Id:               int(id64),
		Name:             name,
		Email:            email,
		HashedPass:       hashedPass,
		LoginAt:          now,
		CreateAt:         now,
		UpdateAt:         now,
		Deleted:          false,
		VCode:            vcode,
		TwoVerificatedAt: sql.NullTime{},
		TwoVerificated:   false,
	}
	return tu, nil
}

func (u *UserRepository) QueryById(userId int) (*TUser, error) {
	// query user
	tu := &TUser{Id: userId}
	s := `SELECT name, email, password, login_at, create_at, update_at, deleted,
					two_step_verification_code, two_verificated_at, two_verificated
				FROM users WHERE id = ?`
	row := DB.QueryRow(s, userId)
	err := row.Scan(&tu.Name, &tu.Email, &tu.HashedPass, &tu.LoginAt, &tu.CreateAt, &tu.UpdateAt, &tu.dbDeleted,
		&tu.VCode, &tu.TwoVerificatedAt, &tu.dbTwoVerificated,
	)
	if err != nil {
		return nil, err
	}

	// map
	tu.Deleted = tu.dbDeleted == 1
	tu.TwoVerificated = tu.dbTwoVerificated == 1
	return tu, nil
}

func (u *UserRepository) QueryByEMail(email string) (*TUser, error) {
	// query user
	tu := &TUser{Email: email}
	s := `SELECT id, name, email, password, login_at, create_at, update_at, deleted,
					two_step_verification_code, two_verificated_at, two_verificated
	 			FROM users WHERE email = ?`
	row := DB.QueryRow(s, email)
	err := row.Scan(&tu.Id, &tu.Name, &tu.Email, &tu.HashedPass, &tu.LoginAt, &tu.CreateAt, &tu.UpdateAt, &tu.dbDeleted,
		&tu.VCode, &tu.TwoVerificatedAt, &tu.dbTwoVerificated,
	)
	if err != nil {
		return nil, err
	}

	// map
	tu.Deleted = tu.dbDeleted == 1
	tu.TwoVerificated = tu.dbTwoVerificated == 1
	return tu, nil
}

// deleted flag on
func (u *UserRepository) SoftDeleteById(userId int) error {
	s := `UPDATE users SET deleted=1 WHERE id = ?`
	_, err := DB.Exec(s, userId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) ActivateById(userId int) error {
	s := `UPDATE users SET deleted=0 WHERE id = ?`
	_, err := DB.Exec(s, userId)
	if err != nil {
		return err
	}
	return nil
}

// Delete row.
func (u *UserRepository) HardDeleteById(userId int) error {
	// exec delete row
	s := `DELETE FROM users WHERE id = ?`
	_, err := DB.Exec(s, userId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) UpdateLoginTimeAndResetVCode(userId int, vcode string, now time.Time) error {
	// update db
	s := `UPDATE users SET two_step_verification_code = ?, two_verificated = 0, login_at = ? WHERE id = ?`
	_, err := DB.Exec(s, vcode, now, userId)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) UpdateVerifiscatedAt(userId int, now time.Time) error {
	// update db
	s := `UPDATE users SET two_verificated = 1, two_verificated_at = ? WHERE id = ?`
	_, err := DB.Exec(s, now, userId)
	if err != nil {
		return err
	}

	return nil
}
