package app

import (
	"context"
	"github.com/EvgeniyBudaev/kafka-go/app/internal/producer/controller"
	"log"
)

var prefix = "/gateway/api/v1"

func (app *App) StartHTTPServer(ctx context.Context) error {
	app.fiber.Static("/static", "./static")
	router := app.fiber.Group(prefix)
	profileController := controller.NewProfileController(app.kafkaWriter)
	router.Post("/profiles/likes", profileController.AddLike())
	go func() {
		port := ":9000"
		if err := app.fiber.Listen(port); err != nil {
			log.Println(err)
		}
	}()
	select {
	case <-ctx.Done():
		log.Println("Context cancelled, shutting down...")
		if err := app.fiber.Shutdown(); err != nil {
			log.Println(err)
		}
	}
	return nil
}
