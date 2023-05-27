package delete

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"invest/internal/models"
	"strconv"
)

func Handler(db *gorm.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var user models.User
		if token, ok := c.GetReqHeaders()["Authorization"]; ok {
			db.Where("token =?", token).First(&user)
		}
		if user.ID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "You don't have access to this page",
			})
		}

		buf := c.Params("id")
		calcID, err := strconv.ParseUint(buf, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		var calculation models.Calculation
		iTx := db.Where("id =?", calcID).First(&calculation)
		if iTx.Error != nil {
			return iTx.Error
		}
		if user.Role != models.Admin && user.ID != *calculation.UserID {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "It's not your data",
			})
		}

		iTx = db.Delete(&calculation)
		return iTx.Error
	}
}
