package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/EvgeniyBudaev/kafka-go/app/internal/consumer/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
	"log"
)

var brokers = []string{"127.0.0.1:10095", "127.0.0.1:10096", "127.0.0.1:10097"}

// App - application structure
type App struct {
	fiber       *fiber.App
	kafkaReader *kafka.Reader
}

// New - create new application
func New() *App {
	// Kafka
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  "consumer-group-id",
		Topic:    "test_topic",
		MaxBytes: 10e6, // 10MB
	})

	// Fiber
	f := fiber.New(fiber.Config{
		ReadBufferSize: 256 << 8,
		BodyLimit:      50 * 1024 * 1024, // 50 MB
	})

	// CORS
	f.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Content-Type, X-Requested-With, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	return &App{
		fiber:       f,
		kafkaReader: r,
	}
}

// Run launches the application
func (app *App) Run(ctx context.Context) {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return app.StartHTTPServer(ctx)
	})

	g.Go(func() error {
		var hc *entity.HubContent

		for {
			select {
			case <-ctx.Done():
				log.Println("Context cancelled, stopping Kafka reader...")
				return app.kafkaReader.Close() // Close the Kafka reader on context cancellation
			default:
				m, err := app.kafkaReader.ReadMessage(ctx)
				if err != nil {
					log.Println("Error reading message:", err)
					return nil
				}

				err = json.Unmarshal(m.Value, &hc)
				if err != nil {
					log.Println("Failed to unmarshal JSON:", err)
					return nil
				}
				fmt.Printf("Message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
			}
		}
	})

	if err := g.Wait(); err != nil {
		log.Printf("Application error: %v", err)
	} else {
		log.Println("Microservice consumer finished successfully")
	}
}
