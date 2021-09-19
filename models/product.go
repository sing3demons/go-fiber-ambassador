package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
}

type _Product struct {
	ID          uint
	Title       string
	Description string
	Image       string
	Price       float64
}
