package controllers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/main/database"
	"github.com/sing3demons/main/middlewares"
	"github.com/sing3demons/main/models"
	"gorm.io/gorm"
)

type Auth struct{}

func (*Auth) DB() *gorm.DB {
	return database.GetDB()
}

func (auth *Auth) Register(c *fiber.Ctx) error {
	var data registerForm
	if err := c.BodyParser(&data); err != nil {
		return c.JSON(err)
	}

	if data.Password != data.PasswordConfirm {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	user := models.User{
		FirstName:    data.FirstName,
		LastName:     data.LastName,
		Email:        data.Email,
		IsAmbassador: strings.Contains(c.Path(), "/api/v1/ambassadors"),
	}

	user.HashPassword(data.Password)

	if err := auth.DB().Create(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": user,
	})
}

func (auth *Auth) Login(c *fiber.Ctx) error {
	var data loginForm

	if err := c.BodyParser(&data); err != nil {
		return c.JSON(err)
	}

	var user models.User

	if err := auth.DB().Where("email = ?", data.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if user.ID == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Credentials",
		})
	}

	if err := user.ComparePassword(data.Password); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	isAmbassador := strings.Contains(c.Path(), "/api/v1/ambassadors")

	var scope string

	if isAmbassador {
		scope = "ambassador"
	} else {
		scope = "admin"
	}

	if !isAmbassador && user.IsAmbassador {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	token, err := middlewares.GenerateJWT(user.ID, scope)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 1),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func (auth *Auth) User(c *fiber.Ctx) error {
	id, err := middlewares.GetUserID(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"error": "unauthenticated",
		})
	}
	var user models.User

	if err := auth.DB().First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	serializeUser := userResponse{}
	copier.Copy(&serializeUser, &user)
	return c.JSON(serializeUser)
}

func (auth *Auth) Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logout success",
	})
}

func (auth *Auth) UpdateInfo(c *fiber.Ctx) error {
	var form updateprofile

	if err := c.BodyParser(&form); err != nil {
		return err
	}

	id, _ := middlewares.GetUserID(c)
	var user models.User
	copier.Copy(&user, &form)
	user.ID = id

	if err := auth.DB().Model(&user).Updates(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	serializeUser := userResponse{}
	copier.Copy(&serializeUser, &user)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": serializeUser,
	})
}

func (auth *Auth) UpdatePassword(c *fiber.Ctx) error {
	var form updatePassword
	if err := c.BodyParser(&form); err != nil {
		return err
	}
	id, _ := middlewares.GetUserID(c)
	var user models.User
	user.ID = id

	if form.Password != form.PasswordConfirm {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	user.HashPassword(form.Password)

	if err := auth.DB().Model(&user).Update("password", user.Password).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
