package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/main/database"
	"github.com/sing3demons/main/models"
	"gorm.io/gorm"
)

type Products struct{}

func (*Products) DB() *gorm.DB {
	return database.GetDB()
}

func (tx *Products) FindProducts(c *fiber.Ctx) error {
	var products []models.Product

	if err := tx.DB().Find(&products).Error; err != nil {
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"products": products,
	})
}
func findByID(c *fiber.Ctx) (uint, error) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (tx *Products) findByProduct(c *fiber.Ctx) (*models.Product, error) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return nil, err
	}

	var product models.Product
	if err := tx.DB().First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (tx *Products) FindProduct(c *fiber.Ctx) error {
	product, err := tx.findByProduct(c)

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"products": product,
	})
}

type productsForm struct {
	Title       string  `json:"title" form:"title"`
	Description string  `json:"description" form:"description"`
	Image       string  `json:"image" form:"image"`
	Price       float64 `json:"price" form:"price"`
}

func (tx *Products) CreateProducts(c *fiber.Ctx) error {
	var product models.Product
	var form productsForm
	if err := c.BodyParser(&form); err != nil {
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	copier.Copy(&product, &form)

	if err := tx.DB().Create(&product).Error; err != nil {
		return c.JSON(err)
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (tx *Products) UpdateProduct(c *fiber.Ctx) error {

	product, err := tx.findByProduct(c)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var form productsForm
	if err := c.BodyParser(&form); err != nil {
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	copier.Copy(&product, &form)
	tx.DB().Model(&product).Updates(&product)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"products": product,
	})
}

func (tx *Products) DeleteProduct(c *fiber.Ctx) error {

	id, err := findByID(c)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := tx.DB().Unscoped().Delete(&models.Product{}, id).Error; err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
