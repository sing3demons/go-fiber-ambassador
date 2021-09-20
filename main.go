package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sing3demons/main/database"
	"github.com/sing3demons/main/routes"
)

func main() {
	database.Connect()
	// database.AutoMigrate()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	app.Get("/healcheck", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	routes.Serve(app)
	app.Listen(":8080")
}
