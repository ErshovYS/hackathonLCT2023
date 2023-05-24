package equipments

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"invest/internal/models"
)

type response struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	PriceUSD uint32 `json:"price_usd"`
	PriceRUB uint32 `json:"price_rub"`
}

func Handler(db *gorm.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var equips []models.EquipmentList
		tx := db.Find(&equips)
		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": tx.Error.Error(),
			})
		}

		res := make([]response, 0, len(equips))
		for _, d := range equips {
			res = append(res, response{
				ID:       d.ID,
				Name:     d.Name,
				PriceUSD: d.PriceUSD,
				PriceRUB: d.PriceRUB,
			})
		}

		return c.JSON(fiber.Map{
			"equipments": res,
		})
	}
}
