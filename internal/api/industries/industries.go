package industries

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"invest/internal/models"
)

type response struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func Handler(db *gorm.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var industries []models.Industry
		tx := db.Select("id, name").Find(&industries)
		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": tx.Error.Error(),
			})
		}

		res := make([]response, 0, len(industries))
		for _, d := range industries {
			res = append(res, response{
				ID:   d.ID,
				Name: d.Name,
			})
		}

		return c.JSON(fiber.Map{
			"industries": res,
		})
	}
}
