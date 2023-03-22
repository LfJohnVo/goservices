package handlers

import (
	"goservices/models"
	"goservices/pkg"

	"github.com/gofiber/fiber/v2"
)

func PostMail(c *fiber.Ctx) error {
	correot := new(models.PostMail)

	if err := c.BodyParser(correot); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	if err := pkg.SendEmail([]string{correot.Email}, correot.Subject, correot.Body, true); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "success",
	})
}
