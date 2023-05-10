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

func (w *WebpushRepository) Query(userId int) (*[]TWebpush, error) {
	return nil, nil
}

func (w *WebpushRepository) Create(userId int, endpoint, p256dh, auth string, expTime *time.Time) error {
	return nil
}

func (w *WebpushRepository) Delete(userId int) error {
	return nil
}
