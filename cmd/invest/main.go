package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"invest/internal/api"
	"invest/internal/config"
	"log"
	"fmt"
)

func main() {
	fmt.Println("Invest started")
	var cfg config.Config
	err := envconfig.Process("invest", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Config: %v\n", cfg)

	zapLogger := zap.NewNop()

	var dbDial gorm.Dialector
	switch cfg.DB.Type {
	case "sqlite":
		dbDial = sqlite.Open(cfg.DB.DSN)
	case "mysql":
		dbDial = mysql.Open(cfg.DB.DSN)
	case "postgres":
		dbDial = postgres.Open(cfg.DB.DSN)
	default:
		log.Fatal("wrong type of database")
	}
	db, err := gorm.Open(dbDial)
	if err != nil {
		log.Fatal(err.Error())
	}

	app := fiber.New()

	restAPI := api.New(app, db, zapLogger)
	restAPI.MakeHandlers()
	fmt.Println("restAPI ready to work")

	err = app.Listen(":80")
	if err != nil {
		zapLogger.Error("failed start server", zap.Error(err))
	}
}
