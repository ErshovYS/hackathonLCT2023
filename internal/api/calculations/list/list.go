package list

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

		limit := 20
		var offset int
		var err error
		qLimit := c.Query("limit")
		if qLimit != "" {
			limit, err = strconv.Atoi(qLimit)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "limit is not a number",
				})
			}
		}
		qOffset := c.Query("offset")
		if qOffset != "" {
			offset, err = strconv.Atoi(qOffset)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "offset is not a number",
				})
			}
		}
		userId := c.Query("user_id")

		var calculations []models.Calculation
		if user.Role == models.Investor {
			iTx := db.Where("user_id =?", user.ID).Order("updated_at desc").Limit(limit).Offset(offset).Find(&calculations)
			if iTx.Error != nil {
				return iTx.Error
			}
		} else if user.Role != models.Admin {
			if userId != "" {
				userID, err := strconv.ParseUint(userId, 10, 64)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"message": "user_id is not a number",
					})
				}
				iTx := db.Where("user_id =?", userID).Order("updated_at desc").Limit(limit).Offset(offset).Find(&calculations)
				if iTx.Error != nil {
					return iTx.Error
				}
			} else {
				iTx := db.Order("updated_at desc").Limit(limit).Offset(offset).Find(&calculations)
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
