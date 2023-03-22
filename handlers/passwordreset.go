package handlers

import (
	"encoding/hex"
	"fiber-boilerplate/config"
	"fiber-boilerplate/databases"
	"fiber-boilerplate/models"
	"fiber-boilerplate/pkg"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Forgot(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	//token := RandStringRunes(12)
	token := tokenGenerator()

	PasswordReset := models.PasswordReset{
		Email: data["email"],
		Token: token,
	}

	to := []string{"mail@mail.com"}
	subject := "Recuperación de contraseña"
	url := config.PWD() + "/api/reset/" + token
	html := true
	message := "Click <a href=\"" + url + "\">here</a> to reset your password!"

	pkg.SendEmail(to, subject, message, html)

	databases.Database.Create(&PasswordReset)

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
	})
}

func ResetPassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		return c.Status(400).JSON(fiber.Map{
			"message": "Password do not match!",
		})
	}

	var PasswordReset = models.PasswordReset{}

	if err := databases.Database.Where("token = ?", data["token"]).Last(&PasswordReset); err.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid token!",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	databases.Database.Model(&models.User{}).Where("email = ?", PasswordReset.Email).Update("password", password)
	return c.Status(200).JSON(fiber.Map{
		"message": "success",
	})
}

func tokenGenerator() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func RandStringRunes(n int) string {
	//int 32 characters
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	//make slice with caracters
	b := make([]rune, n)
	//loop slice and grant integer
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
