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
	done := make(chan struct{})
	profileController := controller.NewProfileController(app.kafkaWriter)
	router.Post("/profiles/likes", profileController.AddLike())
	go func() {
		port := ":9000"
		if err := app.fiber.Listen(port); err != nil {
			log.Println(err)
		}
		close(done)
	}()
	select {
	case <-ctx.Done():
		if err := app.fiber.Shutdown(); err != nil {
			log.Println(err)
		}
	case <-done:
		log.Println("producer server finished successfully")
	}
	return nil
}
