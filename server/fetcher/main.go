package main

import (
	"context"
	"fmt"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func main() {

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         []string{"localhost:9092"},
		GroupID:         "consumer-group-id",
		Topic:           "hl-height",
		MinBytes:        10e3,                   // 10KB
		MaxBytes:        10e6,                   // 10MB
		MaxWait:         100 * time.Millisecond, // 减少等待时间
		ReadLagInterval: -1,                     // 禁用滞后检查以提高性能
		CommitInterval:  time.Second,
	})

	for range 5 {
		msg, err := r.FetchMessage(context.Background())
		if err != nil {
			log.Fatal("failed to fetch message:", err)
		}

		fmt.Println(string(msg.Value))

		err = r.CommitMessages(context.Background(), msg)
		if err != nil {
			log.Fatal("failed to commit message:", err)
		}
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}
