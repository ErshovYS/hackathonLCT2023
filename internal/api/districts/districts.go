package districts

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"invest/internal/models"
)

type response struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func Handler(db *gorm.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var dists []models.District
		tx := db.Find(&dists)
		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": tx.Error.Error(),
			})
		}

		res := make([]response, 0, len(dists))
		for _, d := range dists {
			res = append(res, response{
				ID:    d.ID,
				Name:  d.Name,
				Price: d.Price,
			})
		}

		return c.JSON(fiber.Map{
			"districts": res,
		})
	}
}
