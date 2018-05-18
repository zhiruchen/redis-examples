package db

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
)

// ORM gorm db
var ORM *gorm.DB

// InitMysql init mysql connection pool
func InitMysql() error {
	db, err := gorm.Open("mysql", "root:@tcp(localhost:3306)/testdb?parseTime=true&charset=utf8mb4,utf8&loc=Local")
	if err != nil {
		return err
	}

	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(100)
	ORM = db
	return nil
}
