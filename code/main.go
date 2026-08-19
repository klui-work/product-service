package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"primo/service/internal/di"
	"primo/service/internal/infra/conf"
	"primo/service/internal/infra/database/postgresql"
	productv1 "primo/service/internal/product/handler/v1"

	_ "primo/service/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title          Primo Product API
// @version        1.0
// @description    Product service API for managing products
// @host           localhost:3333
// @BasePath       /
func main() {
	app := fiber.New()

	conf.LoadConfig()

	db := postgresql.SetupPostgresDB()
	postgresql.RunMigration(db)
	postgresql.RunSeed(db)

	productHandler := di.InitializeProductAdapter(db)
	productv1.RegisterRoutes(app, productHandler)

	app.Get("/api-docs/*", swagger.HandlerDefault)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	time.Sleep(3 * time.Second)

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
