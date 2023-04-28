package database

import "testing"

func TestDatabaseConnection(t *testing.T) {
	db, err := GetDatabase()
	if err != nil {
		t.Errorf("got db error '%s'", err)
	}
	if db == nil {
		t.Errorf("cannot connect db")
		return
	}
	defer db.Close()

	if db.Ping() != nil {
		t.Errorf("db ping error got='%s'", err)
	}
}
