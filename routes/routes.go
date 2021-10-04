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
	linkController := controllers.Link{}
	orderController := controllers.Order{}
	{
		adminGroup.Post("register", adminController.Register)
		adminGroup.Post("login", adminController.Login)
		adminAuthenticated := adminGroup.Use(middlewares.IsAuthenticated)
		{
			adminAuthenticated.Get("user", adminController.User)
			adminAuthenticated.Get("logout", adminController.Logout)
			adminAuthenticated.Put("users/info", adminController.UpdateInfo)
			adminAuthenticated.Patch("users/password", adminController.UpdatePassword)
			adminAuthenticated.Get("ambassadors", adminController.Ambassadors)
			adminAuthenticated.Get("ambassadors", adminController.Ambassadors)
		}
		//-->products
		{
			adminAuthenticated.Get("products", productController.FindProducts)
			adminAuthenticated.Get("products/:id", productController.FindProduct)
			adminAuthenticated.Post("products", productController.CreateProducts)
			adminAuthenticated.Put("products/:id", productController.UpdateProduct)
			adminAuthenticated.Delete("products/:id", productController.DeleteProduct)
		}
		{
			adminAuthenticated.Get("users/:id/link", linkController.FindLink)
			adminAuthenticated.Get("orders", orderController.Orders)
		}

	}
	ambassador := v1.Group("ambassadors")
	{
		ambassador.Post("register", adminController.Register)
		ambassador.Post("login", adminController.Login)
		ambassador.Get("products/frontend", productController.ProductsFrontend)
	}
	ambassadorAuthenticated := ambassador.Use(middlewares.IsAuthenticated)
	ambassadorController := controllers.Auth{}
	{
		ambassadorAuthenticated.Get("user", ambassadorController.User)
		ambassadorAuthenticated.Post("logout", ambassadorController.Logout)
		ambassadorAuthenticated.Put("users/info", ambassadorController.UpdateInfo)
		ambassadorAuthenticated.Put("users/password", ambassadorController.UpdatePassword)
		// ambassadorAuthenticated.Post("links", ambassadorController.CreateLink)
		// ambassadorAuthenticated.Get("stats", ambassadorController.Stats)
		// ambassadorAuthenticated.Get("rankings", ambassadorController.Rankings)
	}
}
