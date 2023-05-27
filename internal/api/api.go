package api

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"invest/internal/api/auth/authorization"
	"invest/internal/api/auth/registration"
	user_update "invest/internal/api/auth/update"
	"invest/internal/api/calculations/delete"
	"invest/internal/api/calculations/list"
	"invest/internal/api/calculator"
	"invest/internal/api/districts"
	"invest/internal/api/equipments"
	"invest/internal/api/industries"
	"invest/internal/api/patents"
	"invest/internal/api/regtax"
	"invest/internal/api/users"
)

type API struct {
	app    *fiber.App
	db     *gorm.DB
	logger *zap.Logger
}

func New(app *fiber.App, db *gorm.DB, logger *zap.Logger) *API {
	api := &API{
		app:    app,
		db:     db,
		logger: logger,
	}

	return api
}

func (a *API) MakeHandlers() {
	// POST /registration
	a.app.Post("/registration", registration.Handler(a.db))
	// POST /authorization
	a.app.Post("/authorization", authorization.Handler(a.db))
	// POST /user_update
	a.app.Post("/user_update", user_update.Handler(a.db))

	// GET /districts
	a.app.Get("/districts", districts.Handler(a.db))
	// GET /equipments
	a.app.Get("/equipments", equipments.Handler(a.db))
	// GET /industries
	a.app.Get("/industries", industries.Handler(a.db))
	// GET /patents
	a.app.Get("/patents", patents.Handler(a.db))
	// GET /regtax
	a.app.Get("/regtax", regtax.Handler(a.db))
	// POST /users
	a.app.Post("/users", users.Handler(a.db))

	// GET /calculator
	a.app.Get("/calculator", calculator.HandlerGet(a.db))
	// POST /calculator
	a.app.Post("/calculator", calculator.HandlerPost(a.db))
	// GET /calculations
	a.app.Get("/calculations/list", list.Handler(a.db))
	// DELETE /calculations/{id}
	a.app.Delete("/calculations/:id", delete.Handler(a.db))

	routes := a.app.GetRoutes()
	handlers := make(map[string]string)
	for _, r := range routes {
		if _, ok := handlers[r.Path]; !ok {
			handlers[r.Path] = r.Method
		}
	}

	a.app.Get("/handlers", func(c *fiber.Ctx) error {
		return c.JSON(handlers)
	})
}
