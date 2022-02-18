package db

import (
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func (d *Db) AddNewUser(chatID int64) error {
	db, err := d.OpenDB()
	if err != nil {
		return err
	}
	insert, err := db.Prepare("INSERT INTO `users_status` (`user_id`, `waiting`, `product_name`, `callories`, `current_callories`) VALUES (?, 'no_waiting', '', '', '0')")
	if err != nil {
		return err
	}
	insert.Exec(chatID)
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func (d *Db) UserExists(chatID int64) (bool, error) {
	db, err := d.OpenDB()
	if err != nil {
		return false, err
	}
	var result string
	err = db.QueryRow(fmt.Sprintf("SELECT `user_id` FROM `users_status` WHERE `user_id` = (%d)", chatID)).Scan(&result)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return false, nil
		}
		return false, err
	}
	defer db.Close()
	return true, nil
}
