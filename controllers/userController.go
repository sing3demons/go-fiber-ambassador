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

func (auth *Auth) Rankings(c *fiber.Ctx) error {
	var users []models.User

	auth.DB().Find(&users, models.User{
		IsAmbassador: true,
	})

	var result []interface{}

	for _, user := range users {
		ambassador := models.Ambassador(user)
		ambassador.CalculateRevenue(auth.DB())
		result = append(result, fiber.Map{
			user.Name(): ambassador.Revenue,
		})
	}
	return c.JSON(result)
}
