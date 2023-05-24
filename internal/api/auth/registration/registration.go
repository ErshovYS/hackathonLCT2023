package registration

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"invest/internal/models"
)

type request struct {
	Fullname     string `json:"fullname"`
	Email        string `json:"email"`
	Organization string `json:"organization"`
	INN          string `json:"INN"`
	Site         string `json:"site"`
	IndustryID   uint   `json:"industry_id"`
	Country      string `json:"country"`
	City         string `json:"city"`
	Job          string `json:"job"`
	Password     string `json:"password"`
}

func Handler(db *gorm.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		body := request{}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		user := models.User{
			Fullname:     body.Fullname,
			Email:        body.Email,
			Organization: body.Organization,
			INN:          body.INN,
			Site:         body.Site,
			Country:      body.Country,
			City:         body.City,
			Job:          body.Job,
			Password:     body.Password,
			IndustryID:   body.IndustryID,
			Role:         models.Investor,
		}
		data := []byte(fmt.Sprintf("%s:%s:%d", user.Fullname, user.Password, user.Role))
		user.Token = base64.StdEncoding.EncodeToString(data)

		tx := db.Create(&user)
		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": tx.Error.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"user_id": user.ID,
			"token":   user.Token,
		})
	}
}
