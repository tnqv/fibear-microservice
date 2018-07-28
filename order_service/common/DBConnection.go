package common

import (
	"log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Open a postgres connection
func MysqlConnection() (*gorm.DB, *AppError) {

	dbCon, err := gorm.Open(MAIN_DB_DRIVER, MAIN_DB_CONSTRING)
	if err != nil {
			log.Println("Error : ", err.Error())
			return nil, &AppError{Err: err, Code: ERROR_DATABASE_CONNECTION, Msg: "Error : Cannot connect to db!!"}
	}
	return dbCon, nil
}