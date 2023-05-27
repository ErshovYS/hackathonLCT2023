package regtax

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"invest/internal/models"
)

type response struct {
	RegistrationForms []registrationForm `json:"registration_forms"`
	TaxForms          []taxForm          `json:"tax_forms"`
}

type registrationForm struct {
	ID   models.Registration `json:"id"`
	Name string              `json:"name"`
}

type taxForm struct {
	Name string     `json:"name"`
	ID   models.Tax `json:"id"`
}

func Handler(db *gorm.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		reg := []registrationForm{{models.RegOOO, "ООО"}, {models.RegIP, "ИП"}}
		tax := []taxForm{{"ОСН", models.TaxOCH}, {"УСН", models.TaxYCH}, {"Патент", models.TaxPatent}}

		return c.JSON(fiber.Map{
			"data": response{
				RegistrationForms: reg,
				TaxForms:          tax,
			},
		})
	}
}
