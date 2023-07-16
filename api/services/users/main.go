package users

import (
	"database/sql"
	"errors"
	"himakiwa/services/database"
	"time"

	googleUuid "github.com/google/uuid"
)

var (
	ErrNotExistUser    = errors.New("does not exist users")
	ErrInValidParams   = errors.New("invalid params check your params")
	ErrInvalidEmail    = errors.New("invalid email does not exist")
	ErrInvalidPassword = errors.New("invalid password does not match password")
	ErrInvalidVCode    = errors.New("invalid password does not match vcode")
	ErrUnexpected      = errors.New("unexpected errors occuered")
	ErrInvalidUuid     = errors.New("invalid uuid")
)

type UsersService struct{}

func (u *UsersService) Signin(email, pass, name string) (*database.TUser, error) {
	// validation
	if email == "" || pass == "" {
		return nil, ErrInValidParams
	}
	// hashed password
	hashedPass, err := hashPassword(pass)
	if err != nil {
		return nil, err
	}

	// create verification code
	vcode := generateRandomSixNumber()

	users := &database.UserRepository{}

	tu, err := users.Create(name, email, hashedPass, vcode)
	if err == nil {
		return tu, nil
	}

	// cannnot create user

	// query user
	_, err = users.QueryByEMail(email)
	if err != nil {
		// does not exist.
		return nil, ErrUnexpected
	}

	// already exist
	return nil, ErrInvalidEmail
}

func (u *UsersService) Login(email, pass string) (*database.TUser, error) {
	// validate
	if email == "" || pass == "" {
		return nil, ErrInValidParams
	}

	users := &database.UserRepository{}

	tu, err := users.QueryByEMail(email)
	if err != nil {
		return nil, ErrInvalidEmail
	}

	// check password
	result, err := comparePasswordAndHash(pass, tu.HashedPass)
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, ErrInvalidPassword
	}

	// create verification code
	vcode := generateRandomSixNumber()
	now := time.Now()

	// update vcode
	if err = users.UpdateLoginTimeAndResetVCode(tu.Id, vcode, now); err != nil {
		return nil, err
	}

	tu.VCode = vcode
	tu.LoginAt = now

	return tu, nil
}

func (u *UsersService) Verificate(userId int, vcode string) (*database.TUser, error) {
	if userId == 0 || vcode == "" || len(vcode) != 6 {
		return nil, ErrInValidParams
	}
	users := &database.UserRepository{}

	tu, err := users.QueryById(userId)
	if err != nil {
		return nil, ErrNotExistUser
	}

	// compare code
	result := verificateSixNumber(vcode, tu.VCode, tu.LoginAt)
	if !result {
		return nil, ErrInvalidVCode
	}

	now := time.Now()
	users.UpdateVerifiscatedAt(tu.Id, now)

	tu.TwoVerificatedAt = sql.NullTime{Time: now, Valid: true}
	tu.TwoVerificated = true

	return tu, nil
}

func (u *UsersService) Query(userId int) (*database.TUser, error) {
	if userId == 0 {
		return nil, ErrInValidParams
	}

	users := &database.UserRepository{}
	tu, err := users.QueryById(userId)

	if err != nil {
		return nil, ErrNotExistUser
	}

	return tu, nil
}

func (u *UsersService) QueryByRecruitUuid(uuid string) (*database.TRecruiteUser, error) {
	if uuid == "" {
		return nil, ErrInValidParams
	}
	users := database.UserRepository{}
	tu, err := users.QueryByRecruitUuid(uuid)
	if err != nil {
		return nil, err
	}
	return tu, nil
}

func (u *UsersService) HardDelete(userId int) error {
	if userId == 0 {
		return ErrInValidParams
	}

	users := &database.UserRepository{}
	return users.HardDeleteById(userId)
}

func (u *UsersService) SoftDelete(userId int) error {
	if userId == 0 {
		return ErrInValidParams
	}

	users := &database.UserRepository{}
	return users.SoftDeleteById(userId)
}

type FriendRecruitService struct {
	UserId int
}

func (u *UsersService) GetFriendRecruitService(userId int) FriendRecruitService {
	return FriendRecruitService{UserId: userId}
}

func (f *FriendRecruitService) Query() ([]database.TFRecruitment, error) {
	repository := database.FRecruitmentRepository{}
	recruits, err := repository.QueryByUserId(f.UserId)
	if err != nil {
		return nil, err
	}
	return recruits, nil
}

func (f *FriendRecruitService) UpdateMessageAt(uuid, message string) error {
	// has uuid in user
	recruits, err := f.Query()
	if err != nil {
		return err
	}

	ownUuid := false
	for _, rc := range recruits {
		if rc.Uuid == uuid {
			ownUuid = true
		}
	}
	if !ownUuid {
		return ErrInvalidUuid
	}

	// update
	repository := database.FRecruitmentRepository{}
	err = repository.UpdateMessage(uuid, message)
	if err != nil {
		return err
	}

	return nil
}

func (f *FriendRecruitService) Create(message string) error {
	uuid := googleUuid.NewString()
	repository := database.FRecruitmentRepository{}
	if err := repository.Create(f.UserId, uuid, message); err != nil {
		return err
	}
	return nil
}

func (f *FriendRecruitService) DeleteHard() error {
	repository := database.FRecruitmentRepository{}
	if err := repository.DeleteAll(f.UserId); err != nil {
		return err
	}
	return nil
}
