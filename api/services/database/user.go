package database

import (
	"database/sql"
	"time"
)

type UserRepositories struct {
	UserRepository                UserRepositoryInterface
	RecruitmentRepository         RecruitmentRepositoryInterface
	WebpushSubscriptionRepository WebpushSubscriptionRepositoryInterface
}

func NewUserRepositories() *UserRepositories {
	return &UserRepositories{
		&UserRepository{},
		&RecruitmentRepository{},
		&WebpushSubscriptionRepository{},
	}
}

type UserRepository struct{}
type UserRepositoryInterface interface {
	QueryByID(tx *sql.Tx, userID int) (*TQueryUser, error)
	QueryByEMail(tx *sql.Tx, email string) (*TQueryUser, error)
	QueryByRecruitUUID(tx *sql.Tx, uuid string) (*TQueryRecruitUser, error)
	Create(tx *sql.Tx, name, email, hashedPass string) (int, error)
	UpdateLoginTime(tx *sql.Tx, userID int) error
	SoftDeleteByID(tx *sql.Tx, userID int) error
	ActivateByID(tx *sql.Tx, userID int) error
	HardDeleteByID(tx *sql.Tx, userID int) error
	UpdateVCode(tx *sql.Tx, userID int, vcode string) error
	UpdateVerifiscatedAt(tx *sql.Tx, userID int) error
}

type TQueryUser struct {
	ID               int
	Name             string
	HashedPass       string
	Email            string
	LoginAt          sql.NullTime
	CreateAt         time.Time
	UpdateAt         time.Time
	Deleted          bool
	VCode            string
	TwoVerificatedAt sql.NullTime
	TwoVerificated   bool
}

type TQueryRecruitUser struct {
	ID               int
	Name             string
	HashedPass       string
	Email            string
	LoginAt          sql.NullTime
	CreateAt         time.Time
	UpdateAt         time.Time
	Deleted          bool
	VCode            string
	TwoVerificatedAt sql.NullTime
	TwoVerificated   bool

	UUID           string
	Message        string
	RecruitDeleted bool
}

func (ur *UserRepository) QueryByID(tx *sql.Tx, userID int) (*TQueryUser, error) {
	// query user
	query := `SELECT id, name, email, password, login_at, create_at, update_at, deleted,
					two_step_verification_code, two_verificated_at, two_verificated
				FROM users WHERE id = ?`
	row := tx.QueryRow(query, userID)

	result := &TQueryUser{}
	err := row.Scan(&result.ID, &result.Name, &result.Email, &result.HashedPass, &result.LoginAt, &result.CreateAt, &result.UpdateAt, &result.Deleted,
		&result.VCode, &result.TwoVerificatedAt, &result.TwoVerificated,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ur *UserRepository) QueryByEMail(tx *sql.Tx, email string) (*TQueryUser, error) {
	// query user
	query := `SELECT id, name, email, password, login_at, create_at, update_at, deleted,
					two_step_verification_code, two_verificated_at, two_verificated
	 			FROM users WHERE email = ?`
	row := tx.QueryRow(query, email)

	result := &TQueryUser{}
	err := row.Scan(&result.ID, &result.Name, &result.Email, &result.HashedPass, &result.LoginAt, &result.CreateAt, &result.UpdateAt, &result.Deleted,
		&result.VCode, &result.TwoVerificatedAt, &result.TwoVerificated,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ur *UserRepository) QueryByRecruitUUID(tx *sql.Tx, uuid string) (*TQueryRecruitUser, error) {
	// query user
	query := `SELECT u.id, u.name, u.email, u.password, u.login_at, u.create_at, u.update_at, u.deleted,
					two_step_verification_code, u.two_verificated_at, u.two_verificated,
					r.uuid, r.message, r.deleted
	 			FROM recruitments AS r
					JOIN users AS u ON r.user_id = u.id	
				WHERE r.uuid = ?`
	row := tx.QueryRow(query, uuid)

	result := &TQueryRecruitUser{}
	err := row.Scan(&result.ID, &result.Name, &result.Email, &result.HashedPass, &result.LoginAt, &result.CreateAt, &result.UpdateAt, &result.Deleted,
		&result.VCode, &result.TwoVerificatedAt, &result.TwoVerificated,
		&result.UUID, &result.Message, &result.RecruitDeleted,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ur *UserRepository) Create(tx *sql.Tx, name, email, hashedPass string) (int, error) {
	// exec insert
	query := `INSERT INTO users (name, email, password) VALUES(?, ?, ?)`
	result, err := tx.Exec(query, name, email, hashedPass)
	if err != nil {
		return 0, err
	}

	// get id
	id, err := result.LastInsertId()
	return int(id), err
}

func (ur *UserRepository) UpdateLoginTime(tx *sql.Tx, userID int) error {
	query := `UPDATE users SET login_at = NOW() WHERE id = ?`
	_, err := tx.Exec(query, userID)
	return err
}

// deleted flag on
func (ur *UserRepository) SoftDeleteByID(tx *sql.Tx, userID int) error {
	query := `UPDATE users SET deleted=1 WHERE id = ?`
	_, err := tx.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

// deleted flag off
func (ur *UserRepository) ActivateByID(tx *sql.Tx, userID int) error {
	query := `UPDATE users SET deleted = 0 WHERE id = ?`
	_, err := tx.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

// Delete row.
func (ur *UserRepository) HardDeleteByID(tx *sql.Tx, userID int) error {
	// exec delete row
	query := `DELETE FROM users WHERE id = ?`
	_, err := tx.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) UpdateVCode(tx *sql.Tx, userID int, vcode string) error {
	// update db
	query := `UPDATE users SET two_step_verification_code = ?, two_verificated = 0 WHERE id = ?`
	_, err := tx.Exec(query, vcode, userID)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UpdateVerifiscatedAt(tx *sql.Tx, userID int) error {
	// update db
	query := `UPDATE users SET two_verificated = 1, two_verificated_at = NOW() WHERE id = ?`
	_, err := tx.Exec(query, userID)
	if err != nil {
		return err
	}

	return nil
}

type RecruitmentRepository struct{}
type RecruitmentRepositoryInterface interface {
	QueryByUserID(tx *sql.Tx, userID int) ([]*TQueryRecruitment, error)
	QueryByUUID(tx *sql.Tx, uuid string) (*TQueryRecruitment, error)
	Update(tx *sql.Tx, uuid string, message string, deleted bool) error
	Create(tx *sql.Tx, userID int, uuid, message string) (int, error)
	Delete(tx *sql.Tx, uuid string) error
}

type TQueryRecruitment struct {
	ID       int
	UserID   int
	UUID     string
	Message  string
	CreateAt time.Time
	UpdateAt time.Time
	Deleted  bool
}

func (rr *RecruitmentRepository) QueryByUserID(tx *sql.Tx, userID int) ([]*TQueryRecruitment, error) {
	// query
	query := `SELECT id, user_id, uuid, message, create_at, update_at, deleted FROM recruitments WHERE user_id = ?`
	rows, err := tx.Query(query, userID)
	if err != nil {
		return nil, err
	}

	// responsed
	var results []*TQueryRecruitment
	for rows.Next() {
		rlt := &TQueryRecruitment{}
		if err := rows.Scan(&rlt.ID, &rlt.UserID, &rlt.UUID, &rlt.Message, &rlt.CreateAt, &rlt.UpdateAt, &rlt.Deleted); err != nil {
			return nil, err
		}
		results = append(results, rlt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (rr *RecruitmentRepository) QueryByUUID(tx *sql.Tx, uuid string) (*TQueryRecruitment, error) {
	query := `SELECT id, user_id, uuid, message, create_at, update_at, deleted FROM recruitments WHERE uuid = ?`
	row := tx.QueryRow(query, uuid)
	// result
	result := &TQueryRecruitment{}
	err := row.Scan(&result.ID, &result.UserID, &result.UUID, &result.Message, &result.CreateAt, &result.UpdateAt, &result.Deleted)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found, return nil without an error.
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (rr *RecruitmentRepository) Update(tx *sql.Tx, uuid string, message string, deleted bool) error {
	// query
	query := `UPDATE recruitments SET message = ?, deleted = ? WHERE uuid = ? `
	_, err := tx.Exec(query, message, deleted, uuid)
	if err != nil {
		return err
	}
	return nil
}

func (rr *RecruitmentRepository) Create(tx *sql.Tx, userID int, uuid, message string) (int, error) {
	// query
	query := `INSERT INTO recruitments (user_id, uuid, message) VALUE(?,?,?)`
	result, err := tx.Exec(query, userID, uuid, message)
	if err != nil {
		return 0, err
	}

	// get id
	id, err := result.LastInsertId()
	return int(id), err
}

func (rr *RecruitmentRepository) Delete(tx *sql.Tx, uuid string) error {
	// exec delete row
	query := `DELETE FROM recruitments WHERE uuid = ?`
	_, err := tx.Exec(query, uuid)
	if err != nil {
		return err
	}
	return nil
}

type WebpushSubscriptionRepository struct{}
type WebpushSubscriptionRepositoryInterface interface {
	QueryByUserID(tx *sql.Tx, userID int) ([]*TQueryWebpushSubscription, error)
	Create(tx *sql.Tx, userID int, endpoint, p256dh, auth, userAgent string, expTime *time.Time) (int, error)
	DeleteAll(tx *sql.Tx, userID int) error
}

type TQueryWebpushSubscription struct {
	ID             int
	UserID         int
	Endpoint       string
	P256dh         string
	Auth           string
	UserAgent      string
	ExpirationTime sql.NullTime
	CreateAt       time.Time
}

func (wsr *WebpushSubscriptionRepository) QueryByUserID(tx *sql.Tx, userID int) ([]*TQueryWebpushSubscription, error) {
	// query
	query := `SELECT id, user_id, endpoint, p256dh, auth, user_agent, expiration_time, create_at FROM webpush WHERE user_id = ?`
	rows, err := tx.Query(query, userID)
	if err != nil {
		return nil, err
	}

	// responsed
	var results []*TQueryWebpushSubscription
	for rows.Next() {
		rlt := &TQueryWebpushSubscription{}
		if err := rows.Scan(&rlt.ID, &rlt.UserID, &rlt.Endpoint, &rlt.P256dh, &rlt.Auth, &rlt.UserAgent, &rlt.ExpirationTime, &rlt.CreateAt); err != nil {
			return nil, err
		}
		results = append(results, rlt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (wsr *WebpushSubscriptionRepository) Create(tx *sql.Tx, userID int, endpoint, p256dh, auth, userAgent string, expTime *time.Time) (int, error) {
	// exec insert
	query := `INSERT INTO webpush (user_id, endpoint, p256dh, auth, user_agent, expiration_time) VALUES(?, ?, ?, ?, ?, ?)`
	result, err := tx.Exec(query, userID, endpoint, p256dh, auth, userAgent, expTime)
	if err != nil {
		return 0, err
	}

	// get id
	id, err := result.LastInsertId()
	return int(id), err
}

func (wsr *WebpushSubscriptionRepository) DeleteAll(tx *sql.Tx, userID int) error {
	// delete
	query := `DELETE FROM webpush WHERE user_id = ?`
	_, err := tx.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
