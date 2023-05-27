package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"invest/internal/api"
	"invest/internal/config"
	"invest/internal/storage"
	"log"
)

func main() {
	var cfg config.Config
	err := envconfig.Process("invest", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

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

	store := storage.New(cfg.Storage)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Authorization, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	restAPI := api.New(app, db, store, zapLogger)
	restAPI.MakeHandlers()

	err = app.Listen("0.0.0.0:8080")
	if err != nil {
		zapLogger.Error("failed start server", zap.Error(err))
	}
}
