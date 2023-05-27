package calculator

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"invest/internal/models"
)

type request struct {
	IndustryID        *uint               `json:"industry_id"`
	WorkerCount       uint32              `json:"worker_count"`
	DistrictID        uint                `json:"district_id"`
	LandArea          float32             `json:"land_area"`
	CapBuildingArea   float32             `json:"cap_building_area"`
	CapRebuildingArea float32             `json:"cap_rebuilding_area"`
	RegistrationID    models.Registration `json:"registration_id"`
	TaxID             models.Tax          `json:"tax_id"`
	PatentID          *uint               `json:"patent_id"`
	Equipments        []equipmentPost     `json:"equipments"`
	Buildings         []buildingPost      `json:"buildings"`
	CalculationID     uint                `json:"calculation_id"`
}

type response struct {
	PersonalFrom float32 `json:"personal_from"`
	PersonalTo   float32 `json:"personal_to"`
	EstateFrom   float32 `json:"estate_from"`
	EstateTo     float32 `json:"estate_to"`
	TaxFrom      float32 `json:"tax_from"`
	TaxTo        float32 `json:"tax_to"`
	ServiceFrom  float32 `json:"service_from"`
	ServiceTo    float32 `json:"service_to"`
	TotalFrom    float64 `json:"total_from"`
	TotalTo      float64 `json:"total_to"`
	ReportLink   string  `json:"report_link"`
}

type buildingPost struct {
	Name string  `json:"name"`
	Area float32 `json:"area"`
}

type equipmentPost struct {
	Name     string  `json:"name"`
	PriceRUB float32 `json:"price_rub"`
	Count    uint32  `json:"count"`
}

func HandlerPost(db *gorm.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var user models.User
		if token, ok := c.GetReqHeaders()["Authorization"]; ok {
			db.Where("token =?", token).First(&user)
		}

		req := request{}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		var dist models.District
		var regTax models.RegistrationTax
		var ptn models.Patent
		var ind models.Industry
		if err := db.Transaction(func(tx *gorm.DB) error {
			tx.Select("price").Where("id = ?", req.DistrictID).First(&dist)

			tx.Select("from", "to", "fee").Where("registration = ? AND tax = ?", req.RegistrationID, req.TaxID).First(&regTax)

			if req.PatentID != nil {
				tx.Select("price").Where("id =?", *req.PatentID).First(&ptn)
			}
			if req.IndustryID != nil {
				tx.Where("id =?", *req.IndustryID).First(&ind)
			}
			return nil
		}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		var equipmentPrice float32
		equipments := make([]models.Equipment, 0, len(req.Equipments))
		for _, eq := range req.Equipments {
			equipments = append(equipments, models.Equipment{
				Name:     eq.Name,
				PriceRUB: uint32(eq.PriceRUB * 100),
				Count:    eq.Count,
			})
			equipmentPrice += eq.PriceRUB * float32(eq.Count) * 100
		}

		buildings := make([]models.Building, 0, len(req.Buildings))
		for _, eq := range req.Buildings {
			buildings = append(buildings, models.Building{
				Name: eq.Name,
				Area: eq.Area,
			})
		}

		coef := float32(req.WorkerCount) / float32(ind.Workers)

		// personal
		personalSalaryFrom := float32(ind.Salary*req.WorkerCount*12) * 0.85 / 100
		personalSalaryTo := float32(ind.Salary*req.WorkerCount*12) * 1.15 / 100
		personalSocialFrom := personalSalaryFrom * 0.051
		personalSocialTo := personalSalaryTo * 0.051
		personalPensionFrom := personalSalaryFrom * 0.22
		personalPensionTo := personalSalaryTo * 0.22
		//personalNDFLFrom := personalSalaryFrom * 0.13
		//personalNDFLTo := personalSalaryTo * 0.13

		// estate
		estatePriceFrom := float32(dist.Price) * req.LandArea * 0.85 / 100
		estatePriceTo := float32(dist.Price) * req.LandArea * 1.15 / 100
		estateTaxFrom := estatePriceFrom * 0.015 // float32(ind.EstateTax) * coef * 0.85 / 100
		estateTaxTo := estatePriceFrom * 0.015   // float32(ind.EstateTax) * coef * 1.15 / 100

		// taxes
		moscowTaxFrom := float32(ind.MoscowTax) * coef * 0.85 / 100
		moscowTaxTo := float32(ind.MoscowTax) * coef * 1.15 / 100
		propertyTaxFrom := float32(ind.PropertyTax) * coef * 0.85 / 100
		propertyTaxTo := float32(ind.PropertyTax) * coef * 1.15 / 100
		profitTaxFrom := float32(ind.ProfitTax) * coef * 0.85 / 100
		profitTaxTo := float32(ind.ProfitTax) * coef * 1.15 / 100
		transportTaxFrom := float32(ind.TransportTax) * coef * 0.85 / 100
		transportTaxTo := float32(ind.TransportTax) * coef * 1.15 / 100
		otherTaxFrom := float32(ind.OtherTax) * coef * 0.85 / 100
		otherTaxTo := float32(ind.OtherTax) * coef * 1.15 / 100
		govReg := float32(regTax.Fee)
		patentPrice := float32(ptn.Price) / 100

		// service
		capBuildFrom := req.CapBuildingArea * models.CapBuildingFrom / 100
		capBuildTo := req.CapBuildingArea * models.CapBuildingTo / 100
		capRebuildFrom := req.CapRebuildingArea * models.CapRebuildingFrom / 100
		capRebuildTo := req.CapRebuildingArea * models.CapRebuildingTo / 100
		financialFrom := float32(regTax.From * 12)
		financialTo := float32(regTax.To * 12)

		res := response{
			PersonalFrom: personalSalaryFrom + personalSocialFrom + personalPensionFrom,
			PersonalTo:   personalSalaryTo + personalSocialTo + personalPensionTo,
			EstateFrom:   estatePriceFrom + estateTaxFrom + equipmentPrice,
			EstateTo:     estatePriceTo + estateTaxTo + equipmentPrice,
			TaxFrom:      moscowTaxFrom + propertyTaxFrom + profitTaxFrom + transportTaxFrom + transportTaxFrom + otherTaxFrom + govReg + patentPrice,
			TaxTo:        moscowTaxTo + propertyTaxTo + profitTaxTo + transportTaxTo + transportTaxTo + otherTaxTo + govReg + patentPrice,
			ServiceFrom:  capBuildFrom + capRebuildFrom + financialFrom,
			ServiceTo:    capBuildTo + capRebuildTo + financialTo,
			ReportLink:   fmt.Sprintf("[]byte"),
		}
		res.TotalFrom = float64(res.PersonalFrom + res.EstateFrom + res.TaxFrom)
		res.TotalTo = float64(res.PersonalTo + res.EstateTo + res.TaxTo)

		calc := &models.Calculation{
			IndustryID:        req.IndustryID,
			WorkerCount:       req.WorkerCount,
			DistrictID:        req.DistrictID,
			LandArea:          req.LandArea,
			CapBuildingArea:   req.CapBuildingArea,
			CapRebuildingArea: req.CapRebuildingArea,
			RegistrationTaxID: regTax.ID,
			PatentID:          req.PatentID,
			Equipments:        equipments,
			Buildings:         buildings,
			ResultFrom:        res.TotalFrom,
			ResultTo:          res.TotalTo,
		}
		if req.CalculationID != 0 {
			calc.ID = req.CalculationID
		}
		if user.ID != 0 {
			calc.UserID = &user.ID
		}

		tx := db.Save(&calc)

		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": tx.Error.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"result": res,
		})
	}
}
