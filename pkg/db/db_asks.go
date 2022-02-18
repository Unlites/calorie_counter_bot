package db

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func (d *Db) ResetCurrentCallories(chatID int64) error {
	db, err := d.OpenDB()
	if err != nil {
		return err
	}
	delete, err := db.Prepare("UPDATE `users_status` SET `current_callories` = '0' WHERE user_id = (?)")
	if err != nil {
		return err
	}
	_, err = delete.Exec(chatID)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Printf("%d - Reset all parameters.", chatID)
	return nil

}

func (d *Db) SelectFood(chatID int64, productName string) (string, error) {
	db, err := d.OpenDB()
	if err != nil {
		return "", err
	}
	var result string
	err = db.QueryRow(fmt.Sprintf("SELECT `product_name` FROM `food` WHERE product_name = '%s'", productName)).Scan(&result)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "", nil
		}
		return "", err
	}
	defer db.Close()
	log.Printf("%d - Succecfully readed product %s.", chatID, productName)
	return result, nil
}

func (d *Db) InsertFood(productName string, callories string) error {
	db, err := d.OpenDB()
	if err != nil {
		return err
	}
	result, err := db.Prepare("INSERT INTO `food` (`product_name`, `callories`) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = result.Exec(productName, callories)
	if err != nil {
		return err
	}

	defer db.Close()
	log.Printf("Succecfully added product %s with %s callories", productName, callories)
	return nil
}

func (d *Db) IncreaseCurrentCallories(chatID int64, callories string) error {
	db, err := d.OpenDB()
	if err != nil {
		return err
	}
	result, err := db.Prepare("UPDATE `users_status` SET `current_callories` = `current_callories` + (?) WHERE user_id = (?)")
	if err != nil {
		return err
	}
	_, err = result.Exec(callories, chatID)
	if err != nil {
		return err
	}

	defer db.Close()
	log.Printf("%d - Succecfully increased current_callories on %s", chatID, callories)
	return nil
}

func (d *Db) SelectCurrentCallories(chatID int64) (string, error) {
	db, err := d.OpenDB()
	if err != nil {
		return "", err
	}
	var result string
	err = db.QueryRow(fmt.Sprintf("SELECT `current_callories` FROM `users_status` WHERE `user_id` = %d", chatID)).Scan(&result)
	if err != nil {
		return "", err
	}
	insert, err := db.Prepare("INSERT INTO `lunchs` (`user_id`,`date`, `callories`) VALUES (?, now(), ?)")
	if err != nil {
		return "", err
	}
	insert.Exec(chatID, result)
	if err != nil {
		return "", err
	}
	delete, err := db.Prepare("UPDATE `users_status` SET `current_callories` = '0' WHERE `user_id` = (?)")
	if err != nil {
		return "", err
	}
	_, err = delete.Exec(chatID)
	if err != nil {
		return "", err
	}
	defer db.Close()
	log.Printf("%d - Succecfully readed current_callories.", chatID)
	return result, nil
}

func (d *Db) SelectProductCallories(chatID int64, productName string) (string, error) {
	db, err := d.OpenDB()
	if err != nil {
		return "", err
	}
	var result string
	err = db.QueryRow(fmt.Sprintf("SELECT `callories` FROM `food` WHERE `product_name` = '%s'", productName)).Scan(&result)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "", nil
		}
		return "", err
	}
	defer db.Close()
	log.Printf("%d - Succecfully readed %s's callories.", chatID, productName)
	return result, nil
}
