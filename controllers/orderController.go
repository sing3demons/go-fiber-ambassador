package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/main/database"
	"github.com/sing3demons/main/models"
	"gorm.io/gorm"
)

type Order struct{}

func (*Order) DB() *gorm.DB {
	return database.GetDB()
}

func (tx *Order) Orders(c *fiber.Ctx) error {
	var orders []models.Order

	tx.DB().Preload("OrderItems").Find(&orders)

	for i, order := range orders {
		orders[i].Name = order.FullName()
		orders[i].Total = order.GetTotal()
	}

	return c.JSON(orders)
}
