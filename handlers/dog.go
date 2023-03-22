package handlers

import (
	"fiber-boilerplate/databases"
	"fiber-boilerplate/models"

	"github.com/gofiber/fiber/v2"
)

// GetAlldogs is a function to get all dogs data from database
// @Summary Get all dogs
// @Description Get all dogs
// @Tags dogs
// @Accept json
// @Produce json
// @Success 200 {object} ResponseHTTP{data=[]models.Dog}
// @Failure 503 {object} ResponseHTTP{}
// @Router /api/dogs [get]
func GetDogs(c *fiber.Ctx) error {
	var dogs []models.Dog

	databases.Database.Find(&dogs)
	return c.Status(fiber.StatusOK).JSON(dogs)
}

func GetDog(c *fiber.Ctx) error {
	id := c.Params("id")
	var dog models.Dog

	result := databases.Database.Find(&dog, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.Status(200).JSON(&dog)
}

func AddDog(c *fiber.Ctx) error {
	dog := new(models.Dog)

	if err := c.BodyParser(dog); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	databases.Database.Create(&dog)
	return c.Status(201).JSON(dog)
}

func UpdateDog(c *fiber.Ctx) error {
	dog := new(models.Dog)
	id := c.Params("id")

	if err := c.BodyParser(dog); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	databases.Database.Where("id = ?", id).Updates(&dog)
	return c.Status(200).JSON(dog)
}

func RemoveDog(c *fiber.Ctx) error {
	id := c.Params("id")
	var dog models.Dog

	result := databases.Database.Delete(&dog, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}
