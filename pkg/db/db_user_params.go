package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (d *Db) SetWaiting(chatID int64, waiting string) error {
	db, err := d.OpenDB()
	if err != nil {
		return err
	}
	result, err := db.Prepare("UPDATE `users_status` SET `waiting` = (?) WHERE user_id = (?)")
	if err != nil {
		return err
	}
	_, err = result.Exec(waiting, chatID)
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func (d *Db) SetProductName(chatID int64, productName string) error {
	db, err := d.OpenDB()
	if err != nil {
		return err
	}
	result, err := db.Prepare("UPDATE `users_status` SET `product_name` = (?) WHERE user_id = (?)")
	if err != nil {
		return err
	}
	_, err = result.Exec(productName, chatID)
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func (d *Db) SetCallories(chatID int64, callories string) error {
	db, err := d.OpenDB()
	if err != nil {
		return err
	}
	result, err := db.Prepare("UPDATE `users_status` SET `callories` = (?) WHERE user_id = (?)")
	if err != nil {
		return err
	}
	_, err = result.Exec(callories, chatID)
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func (d *Db) Waiting(chatID int64) (string, error) {
	db, err := d.OpenDB()
	if err != nil {
		return "", err
	}
	var result string
	err = db.QueryRow(fmt.Sprintf("SELECT `waiting` FROM `users_status` WHERE `user_id` = %d", chatID)).Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (d *Db) ProductName(chatID int64) (string, error) {
	db, err := d.OpenDB()
	if err != nil {
		return "", err
	}
	var result string
	err = db.QueryRow(fmt.Sprintf("SELECT `product_name` FROM `users_status` WHERE `user_id` = %d", chatID)).Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (d *Db) Callories(chatID int64) (string, error) {
	db, err := d.OpenDB()
	if err != nil {
		return "", err
	}
	var result string
	err = db.QueryRow(fmt.Sprintf("SELECT `callories` FROM `users_status` WHERE `user_id` = %d", chatID)).Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}
