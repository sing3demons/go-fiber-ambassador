package routes

import "github.com/gofiber/fiber/v2"

func Serve(app *fiber.App) {
	v1 := app.Group("/api/v1/")
	ambassador := v1.Group("ambassador/")
	{
		ambassador.Get("products", func(c *fiber.Ctx) error {
			return c.JSON("products")
		})
	}
}
