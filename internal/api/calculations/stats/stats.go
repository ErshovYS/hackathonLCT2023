package stats

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"invest/internal/models"
	"sort"
	"time"
)

type stat struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

func Handler(db *gorm.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var user models.User
		if token, ok := c.GetReqHeaders()["Authorization"]; ok {
			db.Where("token =?", token).First(&user)
		}
		if user.ID == 0 || user.Role != models.Admin {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "You don't have access to this page",
			})
		}

		var startDate time.Time
		endDate := time.Now()
		var err error
		qStart := c.Query("start_date")
		if qStart != "" {
			startDate, err = time.Parse("2006-01-02", qStart)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "invalid start date",
				})
			}
		}
		qEnd := c.Query("end_date")
		if qEnd != "" {
			endDate, err = time.Parse("2006-01-02", qEnd)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "invalid end date",
				})
			}
		}

		var calculations []models.Calculation
		iTx := db.Select("id", "created_at").Where("created_at BETWEEN ? AND ?", startDate, endDate).Order("updated_at desc").Find(&calculations)
		if iTx.Error != nil {
			return iTx.Error
		}

		statsMap := make(map[time.Time]int, len(calculations))
		for _, calc := range calculations {
			statsMap[calc.CreatedAt.Local().Round(24*time.Hour)]++
		}
		stats := make([]stat, 0, len(statsMap))
		for k, v := range statsMap {
			stats = append(stats, stat{Date: k, Count: v})
		}
		sort.Slice(stats, func(i, j int) bool {
			if stats[i].Date.After(stats[j].Date) {
				return true
			}
			return false
		})

		return c.JSON(fiber.Map{
			"stats": stats,
		})
	}
}
