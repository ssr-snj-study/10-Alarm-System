package config

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

func startTopic() {
	// Kafka Producer 설정
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092", // Kafka 브로커 주소
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()
}
