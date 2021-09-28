package models

import "gorm.io/gorm"

type _Product struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
}

type Product struct {
	gorm.Model
	Title       string
	Description string
	Image       string
	Price       float64
}
