package users

import (
	"database/sql"
	"errors"
	"himakiwa/services/database"
	"time"

	guuid "github.com/google/uuid"
)

var (
	ErrInValidParams   = errors.New("invalid params check your params")
	ErrInvalidEmail    = errors.New("invalid email does not exist")
	ErrInvalidPassword = errors.New("invalid password does not match password")
	ErrInvalidVCode    = errors.New("invalid password does not match vcode")
	ErrUnexpected      = errors.New("unexpected errors occuered")
	ErrInvalidUuid     = errors.New("invalid uuid")
)

type UseUserServicesFunc func(loginID int) *UserServices
type UserServices struct {
	repositories *database.UserRepositories
	loginUserID  int
}

func NewUserServices() UseUserServicesFunc {
	us := &UserServices{repositories: database.NewUserRepositories()}
	return func(loginID int) *UserServices {
		us.loginUserID = loginID
		return us
	}
}

/////////////////////////////////////////////////////
//////////// user repository service ////////////////
/////////////////////////////////////////////////////

func (us *UserServices) Signin(email, pass, name string) (int, error) {
	// validation
	if email == "" || pass == "" {
		return 0, ErrInValidParams
	}
	// hashed password
	hashedPass, err := hashPassword(pass)
	if err != nil {
		return 0, err
	}

	var userID int
	err = database.WithTransaction(func(tx *sql.Tx) error {
		var err error
		userID, err = us.repositories.UserRepository.Create(tx, name, email, hashedPass)
		if err == nil {
			return nil
		}

		// cannnot create user
		// query user
		_, err = us.repositories.UserRepository.QueryByEMail(tx, email)
		if err != nil {
			// does not exist.
			return ErrUnexpected
		}

		// already exist
		return ErrInvalidEmail
	})
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (us *UserServices) Login(email, pass string) (int, error) {
	// validate
	if email == "" || pass == "" {
		return 0, ErrInValidParams
	}
	var user *database.TQueryUser
	err := database.WithTransaction(func(tx *sql.Tx) error {
		var err error
		user, err = us.repositories.UserRepository.QueryByEMail(tx, email)
		if err != nil {
			return err
		}
		// check password
		if result, err := comparePasswordAndHash(pass, user.HashedPass); err != nil {
			return err
		} else if !result {
			return ErrInvalidPassword
		}

		// send logn timestamp
		return us.repositories.UserRepository.UpdateLoginTime(tx, user.ID)
	})
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (us *UserServices) RefreshVCode(userID int) (string, error) {
	if userID == 0 {
		return "", ErrInValidParams
	}
	vcode := generateRandomSixNumber()
	err := database.WithTransaction(func(tx *sql.Tx) error {
		return us.repositories.UserRepository.UpdateVCode(tx, userID, vcode)
	})
	if err != nil {
		return "", nil
	}
	return vcode, nil
}

func (us *UserServices) Verificate(userID int, vcode string) error {
	if userID == 0 || vcode == "" || len(vcode) != 6 {
		return ErrInValidParams
	}
	return database.WithTransaction(func(tx *sql.Tx) error {
		user, err := us.repositories.UserRepository.QueryByID(tx, userID)
		if err != nil {
			return err
		}
		if result := verificateSixNumber(vcode, user.VCode, user.LoginAt.Time); !result {
			return ErrInvalidVCode
		}
		return us.repositories.UserRepository.UpdateVerifiscatedAt(tx, userID)
	})
}

func (us *UserServices) GetUser(userID int) (*database.TQueryUser, error) {
	if userID == 0 {
		return nil, ErrInValidParams
	}
	var user *database.TQueryUser
	err := database.WithTransaction(func(tx *sql.Tx) error {
		var err error
		user, err = us.repositories.UserRepository.QueryByID(tx, userID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserServices) GetUserByRecruitUUID(uuid string) (*database.TQueryRecruitUser, error) {
	if uuid == "" {
		return nil, ErrInValidParams
	}
	var user *database.TQueryRecruitUser
	err := database.WithTransaction(func(tx *sql.Tx) error {
		var err error
		user, err = us.repositories.UserRepository.QueryByRecruitUUID(tx, uuid)
		return err
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

/////////////////////////////////////////////////////
///////// recruitment repository service ////////////
/////////////////////////////////////////////////////

func (us *UserServices) GetRecruitments() ([]*database.TQueryRecruitment, error) {
	var recruitments []*database.TQueryRecruitment
	err := database.WithTransaction(func(tx *sql.Tx) error {
		var err error
		recruitments, err = us.repositories.RecruitmentRepository.QueryByUserID(tx, us.loginUserID)
		return err
	})
	return recruitments, err
}

func (us *UserServices) UpdateRecruitment(uuid, message string, deleted bool) error {
	return database.WithTransaction(func(tx *sql.Tx) error {
		// is owned by loginUser
		if recruitment, err := us.repositories.RecruitmentRepository.QueryByUUID(tx, uuid); err != nil {
			return ErrInvalidUuid
		} else if recruitment.UserID != us.loginUserID {
			return ErrInvalidUuid
		}

		// exec
		return us.repositories.RecruitmentRepository.Update(tx, uuid, message, deleted)
	})
}

func (us *UserServices) CreateRecruitment(message string) error {
	uuid := guuid.NewString()
	return database.WithTransaction(func(tx *sql.Tx) error {
		_, err := us.repositories.RecruitmentRepository.Create(tx, us.loginUserID, uuid, message)
		return err
	})
}

func (us *UserServices) DeleteRecruitment(uuid string) error {
	return database.WithTransaction(func(tx *sql.Tx) error {
		return us.repositories.RecruitmentRepository.Delete(tx, uuid)
	})
}

/////////////////////////////////////////////////////
////// webpush subscription repository service //////
/////////////////////////////////////////////////////

func (us *UserServices) GetWebpushSubscriptions() ([]*database.TQueryWebpushSubscription, error) {
	var subscriptions []*database.TQueryWebpushSubscription
	err := database.WithTransaction(func(tx *sql.Tx) error {
		var err error
		subscriptions, err = us.repositories.WebpushSubscriptionRepository.QueryByUserID(tx, us.loginUserID)
		return err
	})
	return subscriptions, err
}

func (us *UserServices) AddWebpushSubscription(endpoint, p256dh, auth, userAgent string, expTime *time.Time) error {
	// users can have a subscription
	return database.WithTransaction(func(tx *sql.Tx) error {
		err := us.repositories.WebpushSubscriptionRepository.DeleteAll(tx, us.loginUserID)
		if err != nil {
			return err
		}
		_, err = us.repositories.WebpushSubscriptionRepository.Create(tx, us.loginUserID, endpoint, p256dh, auth, userAgent, expTime)
		return err
	})
}
