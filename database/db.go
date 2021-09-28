package database

import (
	"github.com/sing3demons/main/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	var err error
	dsn := "root:root@tcp(localhost:3306)/ambassador?charset=utf8&parseTime=True&loc=Local"

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Could not connect with the database!")
	}
	autoMigrate()
}

func GetDB() *gorm.DB {
	return db
}

func autoMigrate() {
	// db.Migrator().DropTable(models.User{})
	db.AutoMigrate(&models.User{})
	// db.AutoMigrate(models.Product{})
}
