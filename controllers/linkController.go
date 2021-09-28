package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/main/database"
	"github.com/sing3demons/main/models"
	"gorm.io/gorm"
)

type Link struct{}

func (*Link) DB() *gorm.DB {
	return database.GetDB()
}

func (tx *Link) FindLink(c *fiber.Ctx) error {
	id, _ := findByID(c)

	var links []models.Link

	tx.DB().Where("user_id = ?", id).Find(&links)

	for i, link := range links {
		var orders []models.Order

		tx.DB().Where("code = ? and complete = true", link.Code).Find(&orders)

		links[i].Orders = orders
	}

	return c.JSON(links)
}
