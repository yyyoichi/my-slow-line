package database

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

//////////////////////////////////
///////// test user //////////////
//////////////////////////////////

type TestingUser struct {
	User         *TQueryUser
	ExpUser      *TQueryUser
	tx           *sql.Tx
	repositories *UserRepositories
}

func (tu *TestingUser) Delete(t *testing.T) error {
	err := tu.repositories.UserRepository.HardDeleteByID(tu.tx, tu.User.ID)
	if err != nil {
		t.Error(err)
	}
	return err
}
func (tu *TestingUser) GetUserRipositories() *UserRepositories {
	return tu.repositories
}

// counter of user to escape double email and name in db
var TestUserCount = 0

// create user in database and return created user, expected userdata and delete function.
func CreateTestingUser(t *testing.T, tx *sql.Tx, urs *UserRepositories) *TestingUser {
	if urs == nil {
		urs = NewUserRepositories()
	}

	TestUserCount += 1
	name := fmt.Sprintf("Test user %d", TestUserCount)
	email := fmt.Sprintf("test%d@example.com", TestUserCount)

	// create
	userID, err := urs.UserRepository.Create(tx, name, email, "pa55word", "123456")
	if err != nil {
		t.Error(err)
	}
	expUser := &TQueryUser{
		userID,
		name,
		"pa55word",
		email,
		sql.NullTime{},
		time.Now(),
		time.Now(),
		false,
		"123456",
		sql.NullTime{},
		false,
	}
	user, err := urs.UserRepository.QueryByID(tx, userID)
	if err != nil {
		t.Error(err)
	}
	return &TestingUser{
		user,
		expUser,
		tx,
		urs,
	}
}

func TestUser(t *testing.T) {
	testUser(t, NewUserRepositories())
}
func TestUserMock(t *testing.T) {
	testUser(t, NewUserRepositoriesMock())
}

func testUser(t *testing.T, repos *UserRepositories) {
	tx, err := DB.Begin()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	ur := repos.UserRepository

	// create and query test
	testingUser1 := CreateTestingUser(t, tx, repos)
	defer testingUser1.Delete(t)

	user1 := testingUser1.User
	expUser1 := testingUser1.ExpUser
	userIsNotNil(t, user1)
	userIsEqual(t, user1, expUser1)

	user1, err = ur.QueryByEMail(tx, expUser1.Email)
	if err != nil {
		t.Error(err)
	}
	userIsNotNil(t, user1)
	userIsEqual(t, user1, expUser1)

	// update
	err = ur.UpdateLoginTime(tx, user1.ID)
	if err != nil {
		t.Error(err)
	}
	user1, err = ur.QueryByID(tx, user1.ID)
	if err != nil {
		t.Error(err)
	}
	if !user1.LoginAt.Valid {
		t.Error("Expected LoginAt.Valid is true, but got='false'")
	}

	err = ur.SoftDeleteByID(tx, user1.ID)
	if err != nil {
		t.Error(err)
	}
	user1, err = ur.QueryByID(tx, user1.ID)
	if err != nil {
		t.Error(err)
	}
	expUser1.Deleted = true
	userIsEqual(t, user1, expUser1)

	err = ur.ActivateByID(tx, user1.ID)
	if err != nil {
		t.Error(err)
	}
	user1, err = ur.QueryByID(tx, user1.ID)
	if err != nil {
		t.Error(err)
	}
	expUser1.Deleted = false
	userIsEqual(t, user1, expUser1)

	err = ur.UpdateVCode(tx, user1.ID, "654321")
	if err != nil {
		t.Error(err)
	}
	user1, err = ur.QueryByID(tx, user1.ID)
	if err != nil {
		t.Error(err)
	}
	expUser1.VCode = "654321"
	userIsEqual(t, user1, expUser1)

	err = ur.UpdateVerifiscatedAt(tx, user1.ID)
	if err != nil {
		t.Error(err)
	}
	user1, err = ur.QueryByID(tx, user1.ID)
	if err != nil {
		t.Error(err)
	}
	expUser1.TwoVerificated = true
	userIsEqual(t, user1, expUser1)
	if !user1.TwoVerificatedAt.Valid {
		t.Error("Expected TwoVerificatedAt.Valid is true, but got='false'")
	}
}

func userIsEqual(t *testing.T, act, exp *TQueryUser) {
	if act.ID != exp.ID {
		t.Errorf("Expected Id is %d but got='%d'", exp.ID, act.ID)
	}
	if act.Name != exp.Name {
		t.Errorf("Expected Name is %s but got='%s'", exp.Name, act.Name)
	}
	if act.Email != exp.Email {
		t.Errorf("Expected Email is %s but got='%s'", exp.Email, act.Email)
	}
	if act.HashedPass != exp.HashedPass {
		t.Errorf("Expected HashedPass is %s but got='%s'", exp.HashedPass, act.HashedPass)
	}
	if act.VCode != exp.VCode {
		t.Errorf("Expected VCode is %s but got='%s'", exp.VCode, act.VCode)
	}
	if act.Deleted != exp.Deleted {
		t.Errorf("Expected Deleted is %v but got='%v'", exp.Deleted, act.Deleted)
	}
	if act.TwoVerificated != exp.TwoVerificated {
		t.Errorf("Expected TwoVerificated is %v but got='%v'", exp.TwoVerificated, act.TwoVerificated)
	}
}

func userIsNotNil(t *testing.T, u *TQueryUser) {
	if u.ID == 0 {
		t.Error("ID got='0'")
	}
	if u.Name == "" {
		t.Error("Name got=''")
	}
	if u.Email == "" {
		t.Error("Email got=''")
	}
	if u.HashedPass == "" {
		t.Error("HashedPass got=''")
	}
	if u.VCode == "" {
		t.Error("VCode got=''")
	}

	if u.UpdateAt.IsZero() {
		t.Error("UpdateAt is zero")
	}
	if u.CreateAt.IsZero() {
		t.Error("CreateAt is zero")
	}
}
