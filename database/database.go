package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
func Connection() gorm.DB {
	db, err := gorm.Open("mysql", "root:secret@/postapp?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database\n")
	}
	return db
}