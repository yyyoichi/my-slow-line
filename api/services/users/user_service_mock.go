package users

import (
	"database/sql"
	"himakiwa/services/database"
)

func NewUserServicesMock() UseUserServicesFunc {
	urs := database.NewUserRepositoriesMock()
	tx := &sql.Tx{}
	// create user 1
	database.CreateTestingUser(tx, urs)
	// create user 2
	database.CreateTestingUser(tx, urs)

	us := &UserServices{repositories: urs}
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
