package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

func main() {
	// make a writer that produces to test_topic, using the least-bytes distribution
	w := &kafka.Writer{
		Addr:         kafka.TCP("127.0.0.1:9095", "27.0.0.1:9096", "127.0.0.1:9097"),
		Topic:        "test_topic",
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    1048576,
		BatchTimeout: 1000,
		Compression:  kafka.Gzip,
		RequiredAcks: kafka.RequireOne,
	}

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!2"),
		},
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   []byte("Key-C"),
			Value: []byte("Two!2"),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
