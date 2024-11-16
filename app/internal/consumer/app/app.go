package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/EvgeniyBudaev/kafka-go/app/internal/consumer/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
)

var brokers = []string{"127.0.0.1:9095", "127.0.0.1:9096", "127.0.0.1:9097"}

// App - application structure
type App struct {
	fiber *fiber.App
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
	var hc *entity.HubContent

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("err: ", err)
			break
		}
		// Deserialize the JSON data
		err = json.Unmarshal(m.Value, &hc)
		if err != nil {
			log.Fatal("failed from json unmarshal:", err)
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
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
		fiber: f,
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
