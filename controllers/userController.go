package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/main/models"
)

type User struct{}

func (auth *Auth) Ambassadors(c *fiber.Ctx) error {
	var users []models.User

	auth.DB().Find(&users)
	return c.JSON(users)
}
