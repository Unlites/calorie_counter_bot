package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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
	log.Print(chatID) // DEBUG
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

func (d *Db) SelectFood(productName string) (string, error) {
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
	log.Printf("Succecfully readed product %s.", productName)
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

func (d *Db) IncreaseCurrentCallories(callories string) error {
	db, err := d.OpenDB()
	if err != nil {
		return err
	}
	result, err := db.Prepare("UPDATE `current_callories` SET `callories` = `callories` + (?)")
	if err != nil {
		return err
	}
	_, err = result.Exec(callories)
	if err != nil {
		return err
	}

	defer db.Close()
	log.Printf("Succecfully increased current_callories on %s", callories)
	return nil
}

func (d *Db) SelectCurrentCallories() (string, error) {
	db, err := d.OpenDB()
	if err != nil {
		return "", err
	}
	var result string
	err = db.QueryRow("SELECT `callories` FROM `current_callories`").Scan(&result)
	if err != nil {
		return "", err
	}
	insert, err := db.Prepare("INSERT INTO `lunchs` (`date`, `callories`) VALUES (now(), ?)")
	if err != nil {
		return "", err
	}
	insert.Exec(result)
	if err != nil {
		return "", err
	}
	delete, err := db.Prepare("UPDATE `current_callories` SET `callories` = '0'")
	if err != nil {
		return "", err
	}
	_, err = delete.Exec()
	if err != nil {
		return "", err
	}
	defer db.Close()
	log.Print("Succecfully readed current_callories.")
	return result, nil
}

func (d *Db) SelectProductCallories(productName string) (string, error) {
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
	log.Printf("Succecfully readed %s's callories.", productName)
	return result, nil
}

func (d *Db) SelectDayCallories() (string, string, error) {
	db, err := d.OpenDB()
	if err != nil {
		return "", "", err
	}
	var sum string
	err = db.QueryRow("SELECT SUM(callories) FROM lunchs WHERE DATE(`date`) = CURRENT_DATE()").Scan(&sum)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "0", "0", nil
		}
		return "", "", err
	}
	var avg string
	err = db.QueryRow("SELECT AVG(callories) FROM lunchs WHERE DATE(`date`) = CURRENT_DATE()").Scan(&avg)
	if err != nil {
		return "", "", err
	}
	defer db.Close()
	log.Println("Succecfully formed day report.")
	return sum, avg, nil
}

func (d *Db) SelectWeekCallories() (string, string, error) {
	db, err := d.OpenDB()
	if err != nil {
		return "", "", err
	}
	var sum string
	err = db.QueryRow("SELECT SUM(callories) FROM lunchs WHERE `date` BETWEEN  now() - interval 1 week and now()").Scan(&sum)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "0", "0", nil
		}
		return "", "", err
	}
	var avg string
	err = db.QueryRow("SELECT AVG(callories) FROM lunchs WHERE `date` BETWEEN  now() - interval 1 week and now()").Scan(&avg)
	if err != nil {
		return "", "", err
	}
	defer db.Close()
	log.Println("Succecfully formed week report.")
	return sum, avg, nil
}

func (d *Db) SelectMonthCallories() (string, string, error) {
	db, err := d.OpenDB()
	if err != nil {
		return "", "", err
	}
	var sum string
	err = db.QueryRow("SELECT SUM(callories) FROM lunchs WHERE `date` BETWEEN  now() - interval 1 month and now()").Scan(&sum)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "0", "0", nil
		}
		return "", "", err
	}
	var avg string
	err = db.QueryRow("SELECT AVG(callories) FROM lunchs WHERE `date` BETWEEN  now() - interval 1 month and now()").Scan(&avg)
	if err != nil {
		return "", "", err
	}
	defer db.Close()
	log.Println("Succecfully formed month report.")
	return sum, avg, nil
}

func (d *Db) ResetCurrentCallories() error {
	db, err := d.OpenDB()
	if err != nil {
		return err
	}
	delete, err := db.Prepare("UPDATE `current_callories` SET `callories` = '0'")
	if err != nil {
		return err
	}
	_, err = delete.Exec()
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}
