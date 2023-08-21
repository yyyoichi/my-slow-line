package services

import (
	"himakiwa/services/sessions"
	"himakiwa/services/users"
)

func NewRepositoryServicesMock() UseRepositoryServices {
	useUser := users.NewUserServicesMock()
	useSession := sessions.NewSessionServicesMock()
	return func(loginID int) *RepositoryServices {
		return &RepositoryServices{
			useUser(loginID),
			useSession(loginID),
		}
	}
}
