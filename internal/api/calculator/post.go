package calculator

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"invest/internal/models"
	"invest/internal/string_generator"
	"strconv"
	"time"
)

type Storage interface {
	UploadFile(filename string, file []byte) (string, error)
}

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
	PersonalFrom float64 `json:"personal_from"`
	PersonalTo   float64 `json:"personal_to"`
	EstateFrom   float64 `json:"estate_from"`
	EstateTo     float64 `json:"estate_to"`
	TaxFrom      float64 `json:"tax_from"`
	TaxTo        float64 `json:"tax_to"`
	ServiceFrom  float64 `json:"service_from"`
	ServiceTo    float64 `json:"service_to"`
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
	PriceRUB float64 `json:"price_rub"`
	Count    uint32  `json:"count"`
}

func HandlerPost(db *gorm.DB, storage Storage) func(c *fiber.Ctx) error {
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
			tx.Select("name", "price").Where("id = ?", req.DistrictID).First(&dist)

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

		var equipmentPrice float64
		equipments := make([]models.Equipment, 0, len(req.Equipments))
		for _, eq := range req.Equipments {
			equipments = append(equipments, models.Equipment{
				Name:     eq.Name,
				PriceRUB: uint32(eq.PriceRUB * 100),
				Count:    eq.Count,
			})
			equipmentPrice += eq.PriceRUB * float64(eq.Count)
		}

		buildings := make([]models.Building, 0, len(req.Buildings))
		for _, eq := range req.Buildings {
			buildings = append(buildings, models.Building{
				Name: eq.Name,
				Area: eq.Area,
			})
		}

		coef := float64(req.WorkerCount) / float64(ind.Workers)

		// personal
		personalSalaryFrom := float64(ind.Salary*req.WorkerCount*12) * 0.85 / 100
		personalSalaryTo := float64(ind.Salary*req.WorkerCount*12) * 1.15 / 100
		personalSocialFrom := personalSalaryFrom * 0.051
		personalSocialTo := personalSalaryTo * 0.051
		personalPensionFrom := personalSalaryFrom * 0.22
		personalPensionTo := personalSalaryTo * 0.22
		personalNDFLFrom := personalSalaryFrom * 0.13
		personalNDFLTo := personalSalaryTo * 0.13

		// estate
		estatePriceFrom := float64(dist.Price) * float64(req.LandArea) * 0.85 / 100
		estatePriceTo := float64(dist.Price) * float64(req.LandArea) * 1.15 / 100
		estateTaxFrom := estatePriceFrom * 0.015 // float64(ind.EstateTax) * coef * 0.85 / 100
		estateTaxTo := estatePriceFrom * 0.015   // float64(ind.EstateTax) * coef * 1.15 / 100

		// taxes
		moscowTaxFrom := float64(ind.MoscowTax) * coef * 0.85 / 100
		moscowTaxTo := float64(ind.MoscowTax) * coef * 1.15 / 100
		propertyTaxFrom := float64(ind.PropertyTax) * coef * 0.85 / 100
		propertyTaxTo := float64(ind.PropertyTax) * coef * 1.15 / 100
		profitTaxFrom := float64(ind.ProfitTax) * coef * 0.85 / 100
		profitTaxTo := float64(ind.ProfitTax) * coef * 1.15 / 100
		transportTaxFrom := float64(ind.TransportTax) * coef * 0.85 / 100
		transportTaxTo := float64(ind.TransportTax) * coef * 1.15 / 100
		otherTaxFrom := float64(ind.OtherTax) * coef * 0.85 / 100
		otherTaxTo := float64(ind.OtherTax) * coef * 1.15 / 100
		govReg := float64(regTax.Fee)
		patentPrice := float64(ptn.Price) / 100

		// service
		capBuildFrom := req.CapBuildingArea * models.CapBuildingFrom / 100
		capBuildTo := req.CapBuildingArea * models.CapBuildingTo / 100
		capRebuildFrom := req.CapRebuildingArea * models.CapRebuildingFrom / 100
		capRebuildTo := req.CapRebuildingArea * models.CapRebuildingTo / 100
		financialFrom := float64(regTax.From * 12)
		financialTo := float64(regTax.To * 12)

		res := response{
			PersonalFrom: personalSalaryFrom + personalSocialFrom + personalPensionFrom,
			PersonalTo:   personalSalaryTo + personalSocialTo + personalPensionTo,
			EstateFrom:   estatePriceFrom + estateTaxFrom + equipmentPrice,
			EstateTo:     estatePriceTo + estateTaxTo + equipmentPrice,
			TaxFrom:      moscowTaxFrom + propertyTaxFrom + profitTaxFrom + transportTaxFrom + transportTaxFrom + otherTaxFrom + govReg + patentPrice,
			TaxTo:        moscowTaxTo + propertyTaxTo + profitTaxTo + transportTaxTo + transportTaxTo + otherTaxTo + govReg + patentPrice,
			ServiceFrom:  float64(capBuildFrom+capRebuildFrom) + financialFrom,
			ServiceTo:    float64(capBuildTo+capRebuildTo) + financialTo,
		}
		res.TotalFrom = res.PersonalFrom + res.EstateFrom + res.TaxFrom
		res.TotalTo = res.PersonalTo + res.EstateTo + res.TaxTo

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

			resReport := map[string]float64{
				"PersonalCount":       float64(req.WorkerCount),
				"PersonalSalaryFrom":  personalSalaryFrom,
				"PersonalSalaryTo":    personalSalaryTo,
				"PersonalSocialFrom":  personalSocialFrom,
				"PersonalSocialTo":    personalSocialTo,
				"PersonalPensionFrom": personalPensionFrom,
				"PersonalPensionTo":   personalPensionTo,
				"PersonalNDFLFrom":    personalNDFLFrom,
				"PersonalNDFLTo":      personalNDFLTo,
				"EstatePriceFrom":     estatePriceFrom,
				"EstatePriceTo":       estatePriceTo,
				"EstateTaxFrom":       estateTaxFrom,
				"EstateTaxTo":         estateTaxTo,
				"EquipmentPrice":      equipmentPrice,
				"MoscowTaxFrom":       moscowTaxFrom,
				"MoscowTaxTo":         moscowTaxTo,
				"PropertyTaxFrom":     propertyTaxFrom,
				"PropertyTaxTo":       propertyTaxTo,
				"ProfitTaxFrom":       profitTaxFrom,
				"ProfitTaxTo":         profitTaxTo,
				"TransportTaxFrom":    transportTaxFrom,
				"TransportTaxTo":      transportTaxTo,
				"OtherTaxFrom":        otherTaxFrom,
				"OtherTaxTo":          otherTaxTo,
				"PatentPrice":         patentPrice,
				"GovReg":              govReg,
				"CapBuildFrom":        float64(capBuildFrom),
				"CapBuildTo":          float64(capBuildTo),
				"CapRebuildFrom":      float64(capRebuildFrom),
				"CapRebuildTo":        float64(capRebuildTo),
				"FinancialFrom":       financialFrom,
				"FinancialTo":         financialTo,
				"TotalFrom":           res.TotalFrom,
				"TotalTo":             res.TotalTo,
			}
			reportString := map[string]string{
				"Industry":     ind.Name,
				"District":     dist.Name,
				"WorkersCount": strconv.FormatUint(uint64(req.WorkerCount), 10),
			}
			if req.RegistrationID == models.RegOOO {
				reportString["Organization"] = "ООО"
			} else {
				reportString["Organization"] = "ИП"
			}
			b, err := string_generator.GenerateReport(reportString, resReport)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": err.Error(),
				})
			}

			filename := fmt.Sprintf("%s%s%s%d.pdf", user.LastName, user.FirstName, user.MiddleName, time.Now().Unix())
			link, err := storage.UploadFile(filename, b)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": err.Error(),
				})
			}

			calc.ReportLink = link
			res.ReportLink = link
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
