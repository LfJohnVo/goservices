package handlers

import (
	"fiber-boilerplate/databases"
	"fiber-boilerplate/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.StandardClaims
}

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	if data["password"] != data["password_confirm"] {
		return c.Status(400).JSON(fiber.Map{
			"message": "Password do not match!",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		Password:  password,
	}

	databases.Database.Create(&user)

	return c.Status(200).JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	var user models.User
	databases.Database.Where("email = ?", data["email"]).First(&user)

	if user.ID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Correo no registrado dentro del sistema",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Incorrect password!",
		})
	}

	claims := jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	//add token to cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	//send to frontend
	c.Cookie(&cookie)

	return c.Status(200).JSON(fiber.Map{
		"jwt": t,
	})
}

func AuthenticatedUser(c *fiber.Ctx) error {
	//get autenticated user using cookies
	cookie := c.Cookies("jwt")

	// decode secret token
	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	//this validate token
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated!",
		})
	}

	//get user id from claims
	claimsData := token.Claims.(*Claims)

	var user models.User
	//find user register user claims cookie
	databases.Database.Where("id = ?", claimsData.Issuer).First(&user)

	return c.Status(200).JSON(user)

}

func Logout(c *fiber.Ctx) error {
	//remove cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	//delete cookie
	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logout success",
	})
}
