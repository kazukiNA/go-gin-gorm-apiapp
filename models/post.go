package models

import (
	"github.com/jinzhu/gorm"
)


type Post struct{
	gorm.Model
	Text string
	JPTime string
	UserRefer uint
	UserEmail string
}