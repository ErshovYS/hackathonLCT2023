package list

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"invest/internal/models"
)

type request struct {
	UserId uint `json:"user_id"`
	Limit  int  `json:"limit"`
	Offset int  `json:"offset"`
}

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

		var req request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		var calculations []models.Calculation
		if user.Role == models.Investor {
			iTx := db.Where("user_id =?", user.ID).Order("updated_at desc").Limit(req.Limit).Offset(req.Offset).Find(&calculations)
			if iTx.Error != nil {
				return iTx.Error
			}
		} else if user.Role != models.Admin {
			if req.UserId != 0 {
				iTx := db.Where("user_id =?", req.UserId).Order("updated_at desc").Limit(req.Limit).Offset(req.Offset).Find(&calculations)
				if iTx.Error != nil {
					return iTx.Error
				}
			} else {
				iTx := db.Order("updated_at desc").Limit(req.Limit).Offset(req.Offset).Find(&calculations)
				if iTx.Error != nil {
					return iTx.Error
				}
			}
		}

		return c.JSON(fiber.Map{
			"calculations": calculations,
		})
	}
}
