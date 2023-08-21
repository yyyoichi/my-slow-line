package users

import (
	"database/sql"
	"himakiwa/services/database"
)

func NewUserServicesMock() UseUserServicesFunc {
	tx := &sql.Tx{}

	useTransaction := func(fn func(tx *sql.Tx) error) error {
		return fn(tx)
	}
	mock := database.NewUserRepositoriesMock()
	us := &UserServices{useTransaction, mock, 0}

	// create user 1
	database.CreateTestingUser(tx, mock)
	// create user 2
	database.CreateTestingUser(tx, mock)

	// loginUser is user 1 //
	us.loginUserID = 1
	us.CreateRecruitment("Hello")
	us.CreateRecruitment("Hi!")
	us.AddWebpushSubscription("endpoint", "p256dh", "auth", "userAgent", nil)

	// loginUser is user 2 //
	us.loginUserID = 2
	us.CreateRecruitment("Hey!")

	// expected [loginID] is 0, 1 or 2
	return func(loginID int) *UserServices {
		us.loginUserID = loginID
		return us
	}
}
