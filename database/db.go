package database

import (
	"github.com/sing3demons/main/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	var err error

	db, err = gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/ambassador"), &gorm.Config{})

	if err != nil {
		panic("Could not connect with the database!")
	}
	db.AutoMigrate(models.Product{})
}

func GetDB() *gorm.DB {
	return db
}

func AutoMigrate() {
	db.AutoMigrate(models.Product{})
}
