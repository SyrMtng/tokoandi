package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"

	"ilmudata/Project1/controllers"
)

func main() {
	// session
	store := session.New()

	// load template engine
	engine := html.New("./views",".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// static
	app.Static("/public","./public")

	// controllers
	prodController := controllers.InitProductController(store)
	authController := controllers.InitAuthController(store)

	prod := app.Group("/products", func(c *fiber.Ctx) error {
		sess,_ := store.Get(c)
		val := sess.Get("username")
		if val != nil {
			return c.Next()
		}

		return c.Redirect("/login")

	})
	prod.Get("/", prodController.IndexProduct)
	prod.Get("/create", prodController.AddProduct)
	prod.Post("/create", prodController.AddPostedProduct)
	prod.Get("/detail/:id", prodController.GetDetailProduct)
	prod.Get("/editproduct/:id", prodController.EditlProduct)
	prod.Post("/editproduct/:id", prodController.EditlPostedProduct)
	prod.Get("/deleteproduct/:id", prodController.DeleteProduct)
	
	app.Get("/login",authController.Login)
	app.Post("/login",authController.LoginPosted)
	app.Get("/Registrasi",authController.Registrasi)
	app.Post("/Registrasi",authController.RegistrasiPosted)
	app.Get("/logout",authController.Logout)
	
	
	app.Get("/Profile", func(c *fiber.Ctx) error {
		sess,_ := store.Get(c)
		val := sess.Get("username")
		if val != nil {
			return c.Next()
		}

		return c.Redirect("/login")

	}, authController.Profile)

	app.Listen(":3000")
}