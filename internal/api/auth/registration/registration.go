package registration

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"invest/internal/models"
	"time"
)

type request struct {
	FirstName    string `json:"first_name"`
	MiddleName   string `json:"middle_name"`
	LastName     string `json:"last_name"`
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
			FirstName:    body.FirstName,
			MiddleName:   body.MiddleName,
			LastName:     body.LastName,
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
		data := []byte(fmt.Sprintf("%s:%s:%s:%d", user.FirstName, user.LastName, user.Password, time.Now().Unix()))
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
