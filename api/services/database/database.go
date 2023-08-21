package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	var err error
	DB, err = GetDatabase()
	if err != nil || DB.Ping() != nil {
		log.Panic(err, DB.Ping())
	}
}

func getConf() string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWOR")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_DATABASE")
	conf := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, password, host, port, dbname,
	)
	return conf
}

func GetDatabase() (*sql.DB, error) {
	var db *sql.DB
	db, err := sql.Open("mysql", getConf())
	if err != nil {
		return nil, err
	}
	return db, err
}

type TUseTransaction func(func(tx *sql.Tx) error) error

func UseTransaction(fn func(tx *sql.Tx) error) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
