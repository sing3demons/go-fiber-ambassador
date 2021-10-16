package controllers

import (
	"github.com/bxcodec/faker/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/main/database"
	"github.com/sing3demons/main/middlewares"
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

type CreateLinkRequest struct {
	Products []int
}

func (tx *Link) CreateLink(c *fiber.Ctx) error {
	var request CreateLinkRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	id, _ := middlewares.GetUserID(c)

	link := models.Link{
		UserId: id,
		Code:   faker.Username(),
	}

	for _, productId := range request.Products {
		product := models.Product{}
		product.ID = uint(productId)
		link.Products = append(link.Products, product)
	}

	tx.DB().Create(&link)

	return c.JSON(link)
}

func (tx *Link) Stats(c *fiber.Ctx) error {
	id, _ := middlewares.GetUserID(c)
	var links []models.Link

	tx.DB().Find(&links, models.Link{
		UserId: id,
	})

	var result []interface{}
	var orders []models.Order

	for _, link := range links {
		tx.DB().Preload("OrderItems").Find(&orders, models.Order{
			Code:     link.Code,
			Complete: true,
		})

		revenue := 0.0

		for _, order := range orders {
			revenue += order.GetTotal()
		}

		result = append(result, fiber.Map{
			"code":    link.Code,
			"count":   len(orders),
			"revenue": revenue,
		})
	}
	return c.JSON(result)
}

func (tx *Link) GetLink(c *fiber.Ctx) error {
	code := c.Params("code")

	link := models.Link{
		Code: code,
	}

	tx.DB().Preload("User").Preload("Products").First(&link)

	return c.JSON(link)
}
