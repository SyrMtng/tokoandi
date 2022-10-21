package controllers

import (
	"strconv"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"ilmudata/Project1/database"
	"ilmudata/Project1/models"
)

// type ProductForm struct {
// 	Email string `form:"email" validate:"required"`
// 	Address string `form:"address" validate:"required"`
// }

type ProductController struct {
	// declare variables
	store *session.Store
	Db *gorm.DB
}
func InitProductController(s *session.Store) *ProductController {
	db := database.InitDb()
	// gorm
	db.AutoMigrate(&models.Product{})
	return &ProductController{Db: db, store: s}
}

// routing
// GET /products
func (controller *ProductController) IndexProduct(c *fiber.Ctx) error {
	// load all products
	var products []models.Product
	err := models.ReadProducts(controller.Db, &products)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.Render("products", fiber.Map{
		"Title": "Daftar Produk",
		"Products": products,
	})
}
// GET /products/create
func (controller *ProductController) AddProduct(c *fiber.Ctx) error {
	return c.Render("addproduct", fiber.Map{
		"Title": "Tambah Produk",
	})
}
// POST /products/create
func (controller *ProductController) AddPostedProduct(c *fiber.Ctx) error {
	var myform models.Product
	
	file, errFile := c.FormFile("gambar")
	if errFile != nil {
		fmt.Println("Error File =", errFile)
	}
	var filename string = file.Filename
	if file != nil {

		errSaveFile := c.SaveFile(file, fmt.Sprintf("./public/upload/%s", filename))
			if errSaveFile != nil {
				fmt.Println("404")
			}
	} else {
		fmt.Println("404")
	}

	if err := c.BodyParser(&myform); err != nil {
		return c.Redirect("/products")
	}

	myform.Gambar = filename

	// save product
	err := models.CreateProduct(controller.Db, &myform)
	if err!=nil {
		return c.Redirect("/products")
	}
	// if succeed
	return c.Redirect("/products")	
}

// GET /products/detail/xxx
func (controller *ProductController) GetDetailProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)


	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.Render("productdetail", fiber.Map{
		"Title": "Detail Produk",
		"Product": product,
	})
}
/// GET products/editproduct/xx
func (controller *ProductController) EditlProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)


	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.Render("editproduct", fiber.Map{
		"Title": "Edit Produk",
		"Product": product,
	})
}
/// POST products/editproduct/xx
func (controller *ProductController) EditlPostedProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	var myform models.Product

	if err := c.BodyParser(&myform); err != nil {
		return c.Redirect("/products")
	}

	file, errFile := c.FormFile("gambar")
	if errFile != nil {
		fmt.Println("Error File =", errFile)
	}
	var filename string = file.Filename
	if file != nil {

		errSaveFile := c.SaveFile(file, fmt.Sprintf("./public/upload/%s", filename))
		if errSaveFile != nil {
			fmt.Println("folder tidak ada")
		}
	} else {
		fmt.Println("Error !!!!!!!!")
	}
	myform.Gambar = filename
	product.Name = myform.Name
	product.Gambar = myform.Gambar
	product.Deskripsi = myform.Deskripsi
	product.Price = myform.Price
	// save product
	models.UpdateProduct(controller.Db, &product)

	return c.Redirect("/products")

}

/// GET /products/deleteproduct/xx
func (controller *ProductController) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)

	var product models.Product
	models.DeleteProductById(controller.Db, &product, idn)
	return c.Redirect("/products")	
}