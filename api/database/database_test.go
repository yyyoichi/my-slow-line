package database

import "testing"

func TestDatabaseConnection(t *testing.T) {
	db := GetDB()
	if db == nil {
		t.Errorf("cannot connect db")
		return
	}
	defer db.DB.Close()

	err := db.DB.Ping()
	if err != nil {
		t.Errorf("db ping error got='%s'", err)
	}
}
