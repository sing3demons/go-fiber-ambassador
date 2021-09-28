package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Title       string
	Description string
	Image       string
	Price       float64
}
