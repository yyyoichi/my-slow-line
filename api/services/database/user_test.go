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
	user1 := testingUser1.User
	expUser1 := testingUser1.ExpUser

	close := func() {
		testingUser1.Delete(t)
		user1, err = ur.QueryByID(tx, user1.ID)
		if err != sql.ErrNoRows {
			t.Errorf("Expected err is '%s', but got='%s'", sql.ErrNoRows, err.Error())
		}
		if user1 != nil {
			t.Error("Expected user1 is nil, but is not nil")
		}
	}
	defer close()

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

	// create double email
	_, err = ur.Create(tx, expUser1.Name, expUser1.Email, expUser1.HashedPass, expUser1.VCode)
	if err == nil {
		t.Error("Expecte err but got='nil'")
	}
}

//////////////////////////////////
////// test recruitment //////////
//////////////////////////////////

func TestRecruitment(t *testing.T) {
	testRecruitment(t, NewUserRepositories())
}
func TestRecruitmentMock(t *testing.T) {
	testRecruitment(t, NewUserRepositoriesMock())
}

func testRecruitment(t *testing.T, repos *UserRepositories) {
	tx, err := DB.Begin()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// create
	testingUser := CreateTestingUser(t, tx, repos)
	defer testingUser.Delete(t)
	userID := testingUser.User.ID

	rr := repos.RecruitmentRepository

	// create and query test
	uuid1 := "abc"
	id1, err := rr.Create(tx, userID, uuid1, "Hello")
	if err != nil {
		t.Error(err)
	}
	recruit1, err := rr.QueryByUUID(tx, uuid1)
	if err != nil {
		t.Error(err)
	}
	expRecruit1 := &TQueryRecruitment{
		id1,
		userID,
		uuid1,
		"Hello",
		time.Now(),
		time.Now(),
		false,
	}
	recruitmentIsEqual(t, recruit1, expRecruit1)
	recruitmentIsNotNil(t, recruit1)

	recruits, err := rr.QueryByUserID(tx, userID)
	if err != nil {
		t.Error(err)
	}
	if len(recruits) != 1 {
		t.Errorf("Expected len(recruits) is 1, but got='%d'", len(recruits))
	}
	recruit1 = recruits[0]
	recruitmentIsEqual(t, recruit1, expRecruit1)
	recruitmentIsNotNil(t, recruit1)

	// query user and recruit
	recruitUser, err := repos.UserRepository.QueryByRecruitUUID(tx, uuid1)
	if err != nil {
		t.Error(err)
	}
	user := &TQueryUser{
		recruitUser.ID,
		recruitUser.Name,
		recruitUser.HashedPass,
		recruitUser.Email,
		recruitUser.LoginAt,
		recruitUser.CreateAt,
		recruitUser.UpdateAt,
		recruitUser.Deleted,
		recruitUser.VCode,
		recruitUser.TwoVerificatedAt,
		recruitUser.TwoVerificated,
	}
	userIsEqual(t, user, testingUser.ExpUser)
	userIsNotNil(t, user)
	if recruitUser.UUID != uuid1 {
		t.Errorf("Expected UUID is '%s', but got='%s'", uuid1, recruitUser.UUID)
	}
	if recruitUser.Message != "Hello" {
		t.Errorf("Expected Message is '%s', but got='%s'", "Hello", recruitUser.Message)
	}
	if recruitUser.RecruitDeleted {
		t.Error("Expected RecruitDeleted is 'false', but got='true'")
	}

	// update
	err = rr.Update(tx, uuid1, "Hi", true)
	if err != nil {
		t.Error(err)
	}
	expRecruit1.Message = "Hi"
	expRecruit1.Deleted = true
	recruit1, err = rr.QueryByUUID(tx, uuid1)
	if err != nil {
		t.Error(err)
	}
	recruitmentIsEqual(t, recruit1, expRecruit1)
	recruitmentIsNotNil(t, recruit1)

	// create
	uuid2 := "aaa"
	_, err = rr.Create(tx, userID, uuid2, "Hi")
	if err != nil {
		t.Error(err)
	}
	recruits, err = rr.QueryByUserID(tx, userID)
	if err != nil {
		t.Error(err)
	}
	if len(recruits) != 2 {
		t.Errorf("Expected len(recruits) is 2, but got='%d'", len(recruits))
	}

	// delete
	err = rr.Delete(tx, uuid1)
	if err != nil {
		t.Error(err)
	}
	err = rr.Delete(tx, uuid2)
	if err != nil {
		t.Error(err)
	}
	recruits, err = rr.QueryByUserID(tx, userID)
	if err != nil {
		t.Error(err)
	}
	if len(recruits) != 0 {
		t.Errorf("Expected len(recruits) is 0, but got='%d'", len(recruits))
	}
	fmt.Print("Done")
}

//////////////////////////////////
/// test webpush-subscription ////
//////////////////////////////////

func TestWebpushSubscription(t *testing.T) {
	testWebpushSubscription(t, NewUserRepositories())
}
func TestWebpushSubscriptionMock(t *testing.T) {
	testWebpushSubscription(t, NewUserRepositoriesMock())
}

func testWebpushSubscription(t *testing.T, repos *UserRepositories) {
	tx, err := DB.Begin()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	testingUser := CreateTestingUser(t, tx, repos)
	defer testingUser.Delete(t)
	userID := testingUser.User.ID

	wsr := repos.WebpushSubscriptionRepository
	// create
	id, err := wsr.Create(tx, userID, "endpoint", "p256dh", "auth", "userAgent", nil)
	if err != nil {
		t.Error(err)
	}
	subscriptions, err := wsr.QueryByUserID(tx, userID)
	if err != nil {
		t.Error(err)
	}
	if len(subscriptions) != 1 {
		t.Errorf("Expected len(subscriptions) is 1, but got='%d'", len(subscriptions))
	}

	subscription := subscriptions[0]
	if subscription.ID != id {
		t.Errorf("Expected ID is '%d', but got='%d'", id, subscription.ID)
	}
	if subscription.UserID != userID {
		t.Errorf("Expected UserID is '%d', but got='%d'", userID, subscription.UserID)
	}
	if subscription.Auth != "auth" {
		t.Errorf("Expected ID is 'auth', but got='%s'", subscription.Auth)
	}
	if subscription.P256dh != "p256dh" {
		t.Errorf("Expected P256dh is 'p256dh', but got='%s'", subscription.Auth)
	}
	if subscription.UserAgent != "userAgent" {
		t.Errorf("Expected UserAgent is 'userAgent', but got='%s'", subscription.Auth)
	}
	if subscription.ExpirationTime.Valid {
		t.Error("Expected ExpirationTime.Valid is 'false', but got'true'")
	}
	if subscription.CreateAt.IsZero() {
		t.Error("Expected CreateAt is not nil, but got='nil'")
	}

	// create
	_, err = wsr.Create(tx, userID, "endpoint", "p256dh", "auth", "userAgent", nil)
	if err != nil {
		t.Error(err)
	}
	subscriptions, err = wsr.QueryByUserID(tx, userID)
	if err != nil {
		t.Error(err)
	}
	if len(subscriptions) != 2 {
		t.Errorf("Expected len(subscriptions) is 2, but got='%d'", len(subscriptions))
	}

	// delete
	err = wsr.DeleteAll(tx, userID)
	if err != nil {
		t.Error(err)
	}
	subscriptions, err = wsr.QueryByUserID(tx, userID)
	if err != nil {
		t.Error(err)
	}
	if len(subscriptions) != 0 {
		t.Errorf("Expected len(subscriptions) is 0, but got='%d'", len(subscriptions))
	}
	fmt.Print("Done")
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

func recruitmentIsEqual(t *testing.T, act, exp *TQueryRecruitment) {
	if act.UserID != exp.UserID {
		t.Errorf("Expected UserID is '%d', but got='%d'", exp.UserID, act.UserID)
	}
	if act.UUID != exp.UUID {
		t.Errorf("Expected UUID is '%s', but got='%s'", exp.UUID, act.UUID)
	}
	if act.Message != exp.Message {
		t.Errorf("Expected Message is '%s', but got='%s'", exp.Message, act.Message)
	}
	if act.Deleted != exp.Deleted {
		t.Errorf("Expected Deleted is '%v', but got='%v'", exp.Deleted, act.Deleted)
	}
}

func recruitmentIsNotNil(t *testing.T, r *TQueryRecruitment) {
	if r.ID == 0 {
		t.Error("Expected ID is not nil, but got='nil'")
	}
	if r.UserID == 0 {
		t.Error("Expected UserID is not nil, but got='nil'")
	}
	if r.UUID == "" {
		t.Error("Expected UUID is not nil, but got='nil'")
	}
	if r.Message == "" {
		t.Error("Expected Message is not nil, but got='nil'")
	}
	if r.CreateAt.IsZero() {
		t.Error("Expected CreateAt is not nil, but got='nil'")
	}
	if r.UpdateAt.IsZero() {
		t.Error("Expected UpdateAt is not nil, but got='nil'")
	}
}
