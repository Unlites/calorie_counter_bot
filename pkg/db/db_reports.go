package db

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func (d *Db) SelectDayCallories(chatID int64) (string, string, error) {
	db, err := d.OpenDB()
	if err != nil {
		return "", "", err
	}
	var sum string
	err = db.QueryRow(fmt.Sprintf("SELECT SUM(callories) FROM lunchs WHERE DATE(`date`) = CURRENT_DATE() AND `user_id` = %d", chatID)).Scan(&sum)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "0", "0", nil
		}
		return "", "", err
	}
	var avg string
	err = db.QueryRow(fmt.Sprintf("SELECT AVG(callories) FROM lunchs WHERE DATE(`date`) = CURRENT_DATE() AND `user_id` = %d", chatID)).Scan(&avg)
	if err != nil {
		return "", "", err
	}
	defer db.Close()
	log.Printf("%d - Succecfully formed day report.", chatID)
	return sum, avg, nil
}

func (d *Db) SelectWeekCallories(chatID int64) (string, string, error) {
	db, err := d.OpenDB()
	if err != nil {
		return "", "", err
	}
	var sum string
	err = db.QueryRow(fmt.Sprintf("SELECT SUM(callories) FROM lunchs WHERE `date` BETWEEN  now() - interval 1 week and now() AND `user_id` = %d", chatID)).Scan(&sum)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "0", "0", nil
		}
		return "", "", err
	}
	var avg string
	err = db.QueryRow(fmt.Sprintf("SELECT AVG(callories) FROM lunchs WHERE `date` BETWEEN  now() - interval 1 week and now() AND `user_id` = %d", chatID)).Scan(&avg)
	if err != nil {
		return "", "", err
	}
	defer db.Close()
	log.Printf("%d - Succecfully formed week report.", chatID)
	return sum, avg, nil
}

func (d *Db) SelectMonthCallories(chatID int64) (string, string, error) {
	db, err := d.OpenDB()
	if err != nil {
		return "", "", err
	}
	var sum string
	err = db.QueryRow(fmt.Sprintf("SELECT SUM(callories) FROM lunchs WHERE `date` BETWEEN  now() - interval 1 month and now() AND `user_id` = %d", chatID)).Scan(&sum)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return "0", "0", nil
		}
		return "", "", err
	}
	var avg string
	err = db.QueryRow(fmt.Sprintf("SELECT AVG(callories) FROM lunchs WHERE `date` BETWEEN  now() - interval 1 week and now() AND `user_id` = %d", chatID)).Scan(&avg)
	if err != nil {
		return "", "", err
	}
	defer db.Close()
	log.Printf("%d - Succecfully formed month report.", chatID)
	return sum, avg, nil
}
