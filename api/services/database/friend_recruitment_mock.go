package database

import (
	"database/sql"
	"errors"
	"time"
)

type MockFRecruitmentRepository struct {
	RecruitmentByID map[int][]TFRecruitment
}

func NewMockFRecruitmentRepository() *MockFRecruitmentRepository {
	return &MockFRecruitmentRepository{make(map[int][]TFRecruitment)}
}

// QueryByUserId is a mock method to retrieve the Recruitment information corresponding to a specified user ID.
func (m *MockFRecruitmentRepository) QueryByUserId(userId int) ([]TFRecruitment, error) {
	if recruitments, ok := m.RecruitmentByID[userId]; ok {
		return recruitments, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockFRecruitmentRepository) QueryByUUID(uuid string) (*TFRecruitment, error) {
	for _, recruits := range m.RecruitmentByID {
		for _, r := range recruits {
			if r.Uuid == uuid {
				return &r, nil
			}
		}
	}
	return nil, sql.ErrNoRows
}

// Update is a mock method that updates the Recruitment information corresponding to the specified UUID
func (m *MockFRecruitmentRepository) Update(uuid string, message string, deleted bool) error {
	for _, recruitments := range m.RecruitmentByID {
		for i, r := range recruitments {
			if r.Uuid == uuid {
				// UUIDが一致するRecruitmentを見つけたら更新する
				recruitments[i].Message = message
				recruitments[i].Deleted = deleted
				recruitments[i].UpdateAt = time.Now()
				return nil
			}
		}
	}
	return errors.New("recruitment not found")
}

// Create is a mock method that creates new Recruitment information for a given user ID
func (m *MockFRecruitmentRepository) Create(userId int, uuid, message string) error {
	recruitments := m.RecruitmentByID[userId]
	id := m.getLen() + 1
	newRecruitment := TFRecruitment{
		id,
		userId,
		uuid,
		message,
		time.Now(),
		time.Now(),
		false,
		0,
	}
	m.RecruitmentByID[userId] = append(recruitments, newRecruitment)
	return nil
}

// DeleteAll is a mock method that deletes all Recruitment information corresponding to a specified user ID
func (m *MockFRecruitmentRepository) DeleteAll(userId int) error {
	delete(m.RecruitmentByID, userId)
	return nil
}

func (m *MockFRecruitmentRepository) getLen() int {
	length := 0
	for _, recruits := range m.RecruitmentByID {
		length += len(recruits)
	}
	return length
}
