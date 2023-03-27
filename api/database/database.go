package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type MariaDB struct {
	DB *sql.DB
}

func GetDB() *MariaDB {
	db, err := getDatabase()
	if err != nil {
		fmt.Println(err)
	}
	return &MariaDB{DB: db}
}

func getConf() string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWOR")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_DATABASE")
	conf := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		user, password, host, port, dbname,
	)
	return conf
}

func getDatabase() (*sql.DB, error) {
	var db *sql.DB
	db, err := sql.Open("mysql", getConf())
	if err != nil {
		return nil, err
	}
	return db, err
}
