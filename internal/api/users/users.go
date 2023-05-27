package users

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"invest/internal/models"
)

type request struct {
	IDs []uint `json:"ids"`
}

type response struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

func Handler(db *gorm.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var user models.User
		if token, ok := c.GetReqHeaders()["Authorization"]; ok {
			db.Where("token =?", token).First(&user)
		}
		if user.ID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "You don't have access to this page",
			})
		}

		req := request{}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		var users []models.User
		var tx *gorm.DB
		if req.IDs != nil && len(req.IDs) > 0 {
			tx = db.Select("id", "first_name", "middle_name", "last_name").Where("id IN ?", req.IDs).Find(&users)
		} else {
			tx = db.Select("id", "first_name", "middle_name", "last_name").Find(&users)
		}
		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": tx.Error.Error(),
			})
		}

		res := make([]response, 0, len(users))
		for _, u := range users {
			if u.ID != user.ID && user.Role != models.Admin {
				continue
			}
			res = append(res, response{
				ID:         u.ID,
				FirstName:  u.FirstName,
				MiddleName: u.MiddleName,
				LastName:   u.LastName,
			})
		}

		return c.JSON(fiber.Map{
			"users": res,
		})
	}
}
