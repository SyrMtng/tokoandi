package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"ilmudata/Project1/database"
	"ilmudata/Project1/models"
)

type LoginForm struct {
	// declare variables
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

type AuthController struct {
	// declare variables
	store *session.Store
	Db *gorm.DB
}
func InitAuthController(s *session.Store) *AuthController {
	db := database.InitDb()
	// gorm
	db.AutoMigrate(&models.User{})
	return &AuthController{Db: db, store: s}
}
// get /login
		func (controller *AuthController) Login(c *fiber.Ctx) error {
			return c.Render("login", fiber.Map{
				"Title": "Login",
			})
		}
		// post /login
		func (controller *AuthController) LoginPosted(c *fiber.Ctx) error {
			sess, err := controller.store.Get(c)
			if err != nil {
				panic(err)
			}

			var user models.User
			var myform LoginForm
			if err := c.BodyParser(&myform); err != nil {
				return c.Redirect("/login")
			}

			er := models.FindByUsername(controller.Db, &user, myform.Username)
			if er != nil {
				return c.Redirect("/login") // http 500 internal server error
			}

			// hardcode auth
			mycompare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(myform.Password))
			if mycompare != nil {
				sess.Set("username", user.Username)
				sess.Set("userID", user.Id)
				sess.Save()

				return c.Redirect("/products")
			}
			return c.Redirect("/login")
		}

func (controller *AuthController) Registrasi(c *fiber.Ctx) error {
	return c.Render("registrasi", fiber.Map{
		"Title": "Registrasi",
	})
}
// post /login
func (controller *AuthController) RegistrasiPosted(c *fiber.Ctx) error {
	var myform models.User
	var convertpass LoginForm

	if err := c.BodyParser(&myform); err != nil {
		return c.Redirect("/login")
	}
	comvertpassword, _ := bcrypt.GenerateFromPassword([]byte(convertpass.Password), 10)
	sHash := string(comvertpassword)

	myform.Password = sHash

	// save product
	err := models.CreateUser(controller.Db, &myform)
	if err != nil {
		return c.Redirect("/Registrasi")
	}
	// if succeed
	return c.Redirect("/Registrasi")
}

// /profile
func (controller *AuthController) Profile(c *fiber.Ctx) error {
	sess,err := controller.store.Get(c)
	if err!=nil {
		panic(err)
	}
	val := sess.Get("username")

	return c.JSON(fiber.Map{
		"username": val,
	})
}
// /logout
func (controller *AuthController) Logout(c *fiber.Ctx) error {
	
	sess,err := controller.store.Get(c)
	if err!=nil {
		panic(err)
	}
	sess.Destroy()
	return c.Render("login", fiber.Map{
		"Title": "Login",
	})
}