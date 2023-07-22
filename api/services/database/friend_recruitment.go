package database

import "time"

type FRecruitmentRepositoryInterface interface {
	QueryByUserId(userId int) ([]TFRecruitment, error)
	Update(uuid string, message string, deleted bool) error
	Create(userId int, uuid, message string) error
	DeleteAll(userId int) error
}

type FRecruitmentRepository struct{}

type TFRecruitment struct {
	Id       int
	UserId   int
	Uuid     string
	Message  string
	CreateAt time.Time
	UpdateAt time.Time
	Deleted  bool

	//
	deleted int
}

func (fr *FRecruitmentRepository) QueryByUserId(userId int) ([]TFRecruitment, error) {
	// query
	s := `SELECT id, user_id, uuid, message, create_at, update_at, deleted FROM friend_recruitment WHERE user_id = ?`
	rows, err := DB.Query(s, userId)
	if err != nil {
		return nil, err
	}

	// responsed
	var results []TFRecruitment
	for rows.Next() {
		tfr := TFRecruitment{UserId: userId}
		if err := rows.Scan(&tfr.Id, &tfr.UserId, &tfr.Uuid, &tfr.Message, &tfr.CreateAt, &tfr.UpdateAt, &tfr.deleted); err != nil {
			return nil, err
		}
		// parse deleted
		tfr.Deleted = tfr.deleted == 1
		results = append(results, tfr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (fr *FRecruitmentRepository) Update(uuid string, message string, deleted bool) error {
	deletedFlg := 0
	if deleted {
		deletedFlg = 1
	}
	// query
	s := `UPDATE friend_recruitment SET message = ?, deleted = ? WHERE uuid = ? `
	_, err := DB.Exec(s, message, deletedFlg, uuid)
	if err != nil {
		return err
	}
	return nil
}

func (fr *FRecruitmentRepository) Create(userId int, uuid, message string) error {
	// query
	s := `INSERT INTO friend_recruitment (user_id, uuid, message) VALUE(?,?,?)`
	_, err := DB.Exec(s, userId, uuid, message)
	if err != nil {
		return err
	}
	return nil
}

func (fr *FRecruitmentRepository) DeleteAll(userId int) error {
	// delete
	s := `DELETE FROM friend_recruitment WHERE user_id = ?`
	_, err := DB.Exec(s, userId)
	if err != nil {
		return err
	}
	return nil
}
