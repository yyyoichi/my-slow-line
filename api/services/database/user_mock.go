package database

import (
	"database/sql"
	"time"
)

type UserDataMock struct {
	userByID          map[int]*TQueryUser
	recruitmentByUUID map[string]*TQueryRecruitment
}

func NewUserRepositoriesMock() *UserRepositories {
	userByID := make(map[int]*TQueryUser)
	recritmentByUUID := make(map[string]*TQueryRecruitment)
	mock := &UserDataMock{userByID, recritmentByUUID}
	return &UserRepositories{
		&UserRepositoryMock{mock},
		&RecruitmentRepositoryMock{mock},
	}
}

type UserRepositoryMock struct {
	mock *UserDataMock
}

func (ur *UserRepositoryMock) QueryByID(tx *sql.Tx, userID int) (*TQueryUser, error) {
	user, found := ur.mock.userByID[userID]
	if !found {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func (ur *UserRepositoryMock) QueryByEMail(tx *sql.Tx, email string) (*TQueryUser, error) {
	var result *TQueryUser
	for _, user := range ur.mock.userByID {
		if user.Email == email {
			result = user
			break
		}
	}
	if result == nil {
		return nil, sql.ErrNoRows
	}
	return result, nil
}

func (ur *UserRepositoryMock) QueryByRecruitUUID(tx *sql.Tx, uuid string) (*TQeuryRecruitUser, error) {
	recruit, found := ur.mock.recruitmentByUUID[uuid]
	if !found {
		return nil, sql.ErrNoRows
	}
	user, err := ur.QueryByID(tx, recruit.UserID)
	if err != nil {
		return nil, err
	}
	result := &TQeuryRecruitUser{
		user.ID,
		user.Name,
		user.HashedPass,
		user.Email,
		user.LoginAt,
		user.CreateAt,
		user.UpdateAt,
		user.Deleted,
		user.VCode,
		user.TwoVerificatedAt,
		user.TwoVerificated,

		recruit.UUID,
		recruit.Message,
		recruit.Deleted,
	}
	return result, nil
}

func (ur *UserRepositoryMock) Create(tx *sql.Tx, name, email, hashedPass, vcode string) (int, error) {
	id := 0
	// the largest number to get next id
	for ID := range ur.mock.userByID {
		if id < ID {
			id = ID
		}
	}

	id += 1
	user := &TQueryUser{
		id,
		name,
		hashedPass,
		email,
		time.Now(),
		time.Now(),
		time.Now(),
		false,
		vcode,
		sql.NullTime{},
		false,
	}
	ur.mock.userByID[id] = user
	return id, nil
}

func (ur *UserRepositoryMock) SoftDeleteByID(tx *sql.Tx, userID int) error {
	user, found := ur.mock.userByID[userID]
	if !found {
		return sql.ErrNoRows
	}
	user.Deleted = true
	return nil
}

func (ur *UserRepositoryMock) ActivateByID(tx *sql.Tx, userID int) error {
	user, found := ur.mock.userByID[userID]
	if !found {
		return sql.ErrNoRows
	}
	user.Deleted = false
	return nil
}

func (ur *UserRepositoryMock) HardDeleteByID(tx *sql.Tx, userID int) error {
	delete(ur.mock.userByID, userID)
	return nil
}

func (ur *UserRepositoryMock) UpdateVCode(tx *sql.Tx, userID int, vcode string) error {
	user, found := ur.mock.userByID[userID]
	if !found {
		return sql.ErrNoRows
	}
	user.VCode = vcode
	user.TwoVerificated = false
	return nil
}

func (ur *UserRepositoryMock) UpdateVerifiscatedAt(tx *sql.Tx, userID int) error {
	user, found := ur.mock.userByID[userID]
	if !found {
		return sql.ErrNoRows
	}
	user.TwoVerificatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	user.TwoVerificated = true
	return nil
}

type RecruitmentRepositoryMock struct {
	mock *UserDataMock
}

func (rr *RecruitmentRepositoryMock) QueryByUserID(tx *sql.Tx, userID int) ([]*TQueryRecruitment, error) {
	// query
	var results []*TQueryRecruitment
	for _, recruit := range rr.mock.recruitmentByUUID {
		if recruit.UserID == userID {
			results = append(results, recruit)
		}
	}
	return results, nil
}

func (rr *RecruitmentRepositoryMock) QueryByUUID(tx *sql.Tx, uuid string) (*TQueryRecruitment, error) {
	result, found := rr.mock.recruitmentByUUID[uuid]
	if !found {
		return nil, sql.ErrNoRows
	}
	return result, nil
}

func (rr *RecruitmentRepositoryMock) Update(tx *sql.Tx, uuid string, message string, deleted bool) error {
	recruit, found := rr.mock.recruitmentByUUID[uuid]
	if !found {
		return sql.ErrNoRows
	}
	recruit.Message = message
	recruit.Deleted = deleted
	return nil
}

func (rr *RecruitmentRepositoryMock) Create(tx *sql.Tx, userID int, uuid, message string) (int, error) {
	_, found := rr.mock.recruitmentByUUID[uuid]
	if found {
		return 0, sql.ErrNoRows
	}

	id := 0
	// the largest number to get next id
	for _, recruit := range rr.mock.recruitmentByUUID {
		if id < recruit.ID {
			id = recruit.ID
		}
	}
	id += 1

	recruit := &TQueryRecruitment{
		id,
		userID,
		uuid,
		message,
		time.Now(),
		time.Now(),
		false,
	}

	rr.mock.recruitmentByUUID[uuid] = recruit
	return id, nil
}

func (rr *RecruitmentRepositoryMock) Delete(tx *sql.Tx, uuid string) error {
	delete(rr.mock.recruitmentByUUID, uuid)
	return nil
}
