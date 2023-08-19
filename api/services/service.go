package services

import (
	"himakiwa/services/sessions"
	"himakiwa/services/users"
)

type RepositoryServices struct {
	UserServices    *users.UserServices
	SessionServices *sessions.SessionServices
}
type UseRepositoryServices func(loginID int) *RepositoryServices

func NewRepositoryServices() UseRepositoryServices {
	useUser := users.NewUserServices()
	useSession := sessions.NewSessionServices()
	return func(loginID int) *RepositoryServices {
		return &RepositoryServices{
			useUser(loginID),
			useSession(loginID),
		}
	}
}
