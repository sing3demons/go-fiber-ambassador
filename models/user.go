package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string
	LastName     string
	Email        string `gorm:"unique;not null"`
	Password     string `gorm:"not null"`
	IsAmbassador bool
}

func (user *User) HashPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
	user.Password = string(hashedPassword)
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
