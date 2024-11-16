package app

import (
	"context"
	"log"
)

func (app *App) StartHTTPServer(ctx context.Context) error {
	app.fiber.Static("/static", "./static")

	done := make(chan struct{})
	go func() {
		port := ":9001"
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
		log.Println("consumer server finished successfully")
	}
	return nil
}
