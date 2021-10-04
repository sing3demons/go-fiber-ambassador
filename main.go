package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sing3demons/main/database"
	"github.com/sing3demons/main/routes"
)

func main() {
	database.Connect()
	database.SetupRedis()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		resp := fiber.Map{"status": "ok"}
		return c.Status(fiber.StatusOK).JSON(resp)
	})

	routes.Serve(app)
	app.Listen(":8080")
}
