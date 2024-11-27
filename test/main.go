package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func main() {
	// Kafka 컨슈머 생성
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092", // Kafka 브로커 주소
		"group.id":          "example-group",  // 컨슈머 그룹 ID
		"auto.offset.reset": "earliest",       // 메시지를 처음부터 읽음
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// 구독할 토픽 지정
	topic := "example-topic"
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topics: %v", err)
	}

	log.Println("Consumer started. Waiting for messages...")

	// 메시지 수신
	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v", err)
			continue
		}
		log.Printf("Received message: %s on topic %s", string(msg.Value), *msg.TopicPartition.Topic)
	}
}
