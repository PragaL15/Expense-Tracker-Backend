package main

import (
	"fmt"
	"log"
	"github.com/PragaL15/Expense-Tracker/internal/config"
	"github.com/PragaL15/Expense-Tracker/internal/database"
	"github.com/PragaL15/Expense-Tracker/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.Load()
	database.Connect(cfg.DatabaseURL)
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://expense-tracker-frontend-gamma-lilac.vercel.app", 
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	routes.Register(app)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("listening on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
