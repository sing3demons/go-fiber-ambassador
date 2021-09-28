package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/main/controllers"
	"github.com/sing3demons/main/middlewares"
)

func Serve(app *fiber.App) {
	v1 := app.Group("/api/v1/")

	adminGroup := v1.Group("admin/")
	adminController := controllers.Auth{}
	productController := controllers.Products{}
	{
		adminGroup.Post("register", adminController.Register)
		adminGroup.Post("login", adminController.Login)
		adminAuthenticated := adminGroup.Use(middlewares.IsAuthenticated)
		adminAuthenticated.Get("user", adminController.User)
		adminAuthenticated.Get("logout", adminController.Logout)
		adminAuthenticated.Put("users/info", adminController.UpdateInfo)
		adminAuthenticated.Patch("users/password", adminController.UpdatePassword)
		adminAuthenticated.Get("ambassadors", adminController.Ambassadors)
		adminAuthenticated.Get("ambassadors", adminController.Ambassadors)
		adminAuthenticated.Get("products", productController.FindProducts)
		adminAuthenticated.Get("products/:id", productController.FindProduct)
		adminAuthenticated.Post("products", productController.CreateProducts)
		adminAuthenticated.Put("products/:id", productController.UpdateProduct)
		adminAuthenticated.Delete("products/:id", productController.DeleteProduct)
	}
	ambassador := v1.Group("ambassador/")
	{
		ambassador.Get("products", func(c *fiber.Ctx) error {
			return c.JSON("products")
		})
	}
}
