package calculator

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"invest/internal/models"
)

type district struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type equipment struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	PriceUSD uint32 `json:"price_usd"`
	PriceRUB uint32 `json:"price_rub"`
}

type industry struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type patent struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func HandlerGet(db *gorm.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var dists []models.District
		var equips []models.EquipmentList
		var industries []models.Industry
		var patents []models.Patent

		if err := db.Transaction(func(tx *gorm.DB) error {
			iTx := db.Find(&dists)
			if iTx.Error != nil {
				return iTx.Error
			}

			iTx = db.Find(&equips)
			if iTx.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": iTx.Error.Error(),
				})
			}

			iTx = db.Select("id, name").Find(&industries)
			if iTx.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": iTx.Error.Error(),
				})
			}

			iTx = db.Select("id, name, price").Find(&patents)
			if iTx.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": iTx.Error.Error(),
				})
			}

			return nil
		}); err != nil {

		}

		respDist := make([]district, 0, len(dists))
		for _, d := range dists {
			respDist = append(respDist, district{
				ID:    d.ID,
				Name:  d.Name,
				Price: d.Price,
			})
		}

		respEq := make([]equipment, 0, len(equips))
		for _, d := range equips {
			respEq = append(respEq, equipment{
				ID:       d.ID,
				Name:     d.Name,
				PriceUSD: d.PriceUSD,
				PriceRUB: d.PriceRUB,
			})
		}

		respInd := make([]industry, 0, len(industries))
		for _, d := range industries {
			respInd = append(respInd, industry{
				ID:   d.ID,
				Name: d.Name,
			})
		}

		respPt := make([]patent, 0, len(patents))
		for _, d := range patents {
			respPt = append(respPt, patent{
				ID:    d.ID,
				Name:  d.Name,
				Price: d.Price,
			})
		}

		reg := map[string]models.Registration{"ООО": models.RegOOO, "ИП": models.RegIP}
		tax := map[string]models.Tax{"ОСН": models.TaxOCH, "УСН": models.TaxYCH, "Патент": models.TaxPatent}

		return c.JSON(fiber.Map{
			"districts":     respDist,
			"equipments":    respEq,
			"industries":    respInd,
			"patents":       respPt,
			"registrations": reg,
			"taxes":         tax,
		})
	}
}
