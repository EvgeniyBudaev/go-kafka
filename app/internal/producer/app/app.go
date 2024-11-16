package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
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
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err := app.StartHTTPServer(ctx); err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	wg.Wait()
}
