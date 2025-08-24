package main

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {

	w := &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		Topic:        "hl-height",
		Balancer:     &kafka.LeastBytes{}, // 指定分区的balancer模式为最小字节分布
		RequiredAcks: kafka.RequireAll,    // ack模式
		Async:        true,                // 异步
	}

	err := w.WriteMessages(context.Background(),
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func watchBlock() {

}

func loadProgress() {

}
