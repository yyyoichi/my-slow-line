package services

import (
	"himakiwa/services/users"
	webpush_services "himakiwa/services/webpush"
)

type RepositoryServices struct {
	users   users.UsersService
	webpush webpush_services.WebpushService
}

func NewRepositoryServices() *RepositoryServices {
	return &RepositoryServices{
		users: users.UsersService{},
	}
}

func (s *RepositoryServices) GetUser() users.UsersService {
	return s.users
}

func (s *RepositoryServices) GetWebpush() webpush_services.WebpushService {
	return s.webpush
}
