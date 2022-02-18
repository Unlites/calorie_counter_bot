package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	user     string
	password string
}

func InitDB(user string, password string) *Db {
	return &Db{user: user, password: password}
}

func (d *Db) OpenDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/callorie_counter", d.user, d.password))
	if err != nil {
		log.Print(err)
		return db, err
	}
	return db, nil
}
