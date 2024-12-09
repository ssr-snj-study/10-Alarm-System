package config

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

var stKafka StKafka

type StKafka struct {
	StKafka *kafka.Consumer
}

func KafkaConsumer() StKafka {
	return stKafka
}

func init() {
	consumeTopic()
}

func consumeTopic() {
	// Kafka Consumer 설정
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092", // Kafka 브로커 주소
		"group.id":          "example-group",  // Consumer 그룹 ID
		"auto.offset.reset": "earliest",       // 메시지를 처음부터 읽음
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}

	// 구독할 토픽 설정
	topic := "example-topic"
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}

	log.Println("Consumer started. Waiting for messages...")

	stKafka.StKafka = consumer
}
