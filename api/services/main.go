package services

import (
	"database/sql"
	"himakiwa/services/users"
)

type RepositoryServices struct {
	Users *users.UsersService
}

func NewServices(db *sql.DB) *RepositoryServices {
	return &RepositoryServices{
		Users: users.NewUsersService(db),
	}
}
