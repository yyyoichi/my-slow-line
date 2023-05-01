package services

import (
	"himakiwa/services/users"
)

type RepositoryServices struct {
	users users.UsersService
}

func NewRepositoryServices() *RepositoryServices {
	return &RepositoryServices{
		users: users.UsersService{},
	}
}

func (s *RepositoryServices) GetUser() users.UsersService {
	return s.users
}
