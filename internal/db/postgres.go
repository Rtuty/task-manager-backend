package db

import (
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

// func createTable() error {
// 	db, err := sql.Open("postgres", dbConStr)
// 	if err != nil {
// 		return err
// 	}
// }
