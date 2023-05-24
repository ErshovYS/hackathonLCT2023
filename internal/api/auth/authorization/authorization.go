package authorization

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"invest/internal/models"
)

type request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Handler(db *gorm.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		body := request{}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		user := models.User{}
		tx := db.Select("id", "token").Where("email = ? AND age = ?", body.Email, body.Password).First(&user)
		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": tx.Error.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"user_id": user.ID,
			"token":   user.Token,
		})
	}
}
