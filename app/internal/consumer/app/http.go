package app

import (
	"context"
	"log"
)

func (app *App) StartHTTPServer(ctx context.Context) error {
	app.fiber.Static("/static", "./static")

	go func() {
		port := ":9001"
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
