package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
	"log"
)

// App - application structure
type App struct {
	fiber       *fiber.App
	kafkaWriter *kafka.Writer
}

// New - create new application
func New() *App {
	// Kafka
	w := &kafka.Writer{
		Addr:         kafka.TCP("127.0.0.1:9095", "27.0.0.1:9096", "127.0.0.1:9097"),
		Topic:        "test_topic",
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    1048576,
		BatchTimeout: 1000,
		Compression:  kafka.Gzip,
		RequiredAcks: kafka.RequireOne,
	}

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
		kafkaWriter: w,
	}
}

// Run launches the application
func (app *App) Run(ctx context.Context) {
	g, ctx := errgroup.WithContext(ctx)

	// Start HTTP server
	g.Go(func() error {
		return app.StartHTTPServer(ctx)
	})

	// Close Kafka writer when done
	g.Go(func() error {
		defer func() {
			log.Println("Kafka close")
			if err := app.kafkaWriter.Close(); err != nil {
				log.Println("Error closing kafkaWriter:", err)
			}
		}()
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Printf("Application error: %v", err)
	} else {
		log.Println("Microservice producer finished successfully")
	}
}
