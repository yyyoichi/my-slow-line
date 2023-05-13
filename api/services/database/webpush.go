package database

import (
	"database/sql"
	"time"
)

type WebpushRepository struct{}

type TWebpush struct {
	Id             int
	UserId         int
	Endpoint       string
	P256dh         string
	Auth           string
	ExpirationTime sql.NullTime
	CreateAt       time.Time
}

func (w *WebpushRepository) QueryByUserId(userId int) ([]TWebpush, error) {
	// query
	s := `SELECT id, endpoint, p256dh, auth, expiration_time, create_at FROM webpush WHERE user_id = ?`
	rows, err := DB.Query(s, userId)
	if err != nil {
		return nil, err
	}

	// responsed
	var results []TWebpush
	for rows.Next() {
		twp := TWebpush{UserId: userId}
		if err := rows.Scan(&twp.Id, &twp.Endpoint, &twp.P256dh, &twp.Auth, &twp.ExpirationTime, &twp.CreateAt); err != nil {
			return nil, err
		}
		results = append(results, twp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (w *WebpushRepository) Create(userId int, endpoint, p256dh, auth string, expTime *time.Time) error {
	// exec insert
	s := `INSERT INTO webpush (user_id, endpoint, p256dh, auth, expiration_time, create_at) VALUES(?, ?, ?, ?, ?, ?)`
	now := time.Now()
	_, err := DB.Exec(s, userId, endpoint, p256dh, auth, expTime, now)
	if err != nil {
		return err
	}
	return nil
}

func (w *WebpushRepository) DeleteAll(userId int) error {
	// delete
	s := `DELETE FROM webpush WHERE user_id = ?`
	_, err := DB.Exec(s, userId)
	if err != nil {
		return err
	}
	return nil
}
