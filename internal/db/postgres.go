package db

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type dataSource struct {
	Host, Port, User, Passwd, Dbname, Sslmode string
}

var PstgCon dataSource

func GetConnection() {
	var res *dataSource = &PstgCon
	envVars := []string{"HOST", "PORT", "USER", "PASSWD", "DBNAME", "SSLMODE"}

	for _, v := range envVars {
		value := os.Getenv(v)
		if value == "" {
			panic(fmt.Sprintf("invalid environment variable %s", v))
		} else {
			field := reflect.ValueOf(res).Elem().FieldByNameFunc(
				func(fieldName string) bool {
					return strings.ToLower(fieldName) == strings.ToLower(v)
				})
			if field.IsValid() {
				field.SetString(value)
			}
		}
	}
}

var dbConStr string = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", PstgCon.Host, PstgCon.Port, PstgCon.User, PstgCon.Passwd, PstgCon.Dbname, PstgCon.Sslmode)

func createTable() error {
	db, err := sql.Open("postgres", dbConStr)
	if err != nil {
		return err
	}

	defer db.Close()

	//Создаем таблицу users
	if _, err = db.Exec(`
		CREATE TABLE users(
			ID SERIAL PRIMARY KEY,
			TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			USERNAME TEXT,
			CHAT_ID INT,
			MESSAGE TEXT,
			ANSWER TEXT);`); err != nil {
		return err
	}

	return nil
}

// Данные по пользователям из бота записываем в БД
func collectData(username string, chatid int64, message string, answer []string) error {
	db, err := sql.Open("postgres", dbConStr)
	if err != nil {
		return err
	}

	defer db.Close()

	answ := strings.Join(answer, ", ")
	data := `INSERT INTO users(username, chatid, message, answer) values ($1, $2, $3, $4)`

	if _, err := db.Exec(data, `@`+username, chatid, message, answ); err != nil {
		return err
	}

	return nil
}

func getNumberOfUsers() (int64, error) {
	var count int64

	db, err := sql.Open("postgres", dbConStr)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	//Отправляем запрос в БД для подсчета числа уникальных пользователей
	row := db.QueryRow("SELECT COUNT(DISTINCT username) FROM users;")
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
