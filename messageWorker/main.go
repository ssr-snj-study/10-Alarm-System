package main

import (
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	// Kafka Consumer 설정
	config := &kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092",
		"group.id":           "example-group",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false, // 자동 Offset Commit 비활성화
	}
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// 토픽 구독
	topic := "example-topic"
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	log.Println("Consumer started. Waiting for messages...")

	for {
		// 메시지 읽기
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v", err)
			continue
		}

		log.Printf("Message received: %s\n", string(msg.Value))

		// 메시지 처리
		success := processMessage(msg)
		if success {
			// 메시지 처리 성공: Offset Commit
			_, err := consumer.CommitMessage(msg)
			if err != nil {
				log.Printf("Failed to commit offset: %v", err)
			} else {
				log.Println("Offset committed successfully.")
			}
		} else {
			// 메시지 처리 실패: Retry 로직 실행
			log.Println("Message processing failed. Retrying...")
			retryProcessing(consumer, msg)
		}
	}
}

// 메시지 처리 함수
func processMessage(msg *kafka.Message) bool {
	// 예제: 메시지 값이 "fail"이면 처리 실패
	if string(msg.Value) == "fail" {
		return false
	}

	// 처리 로직 (성공 시 true 반환)
	log.Printf("Processing message: %s\n", string(msg.Value))
	return true
}

// 재시도 로직
func retryProcessing(consumer *kafka.Consumer, msg *kafka.Message) {
	retries := 3                     // 최대 재시도 횟수
	retryInterval := 2 * time.Second // 재시도 간격

	for i := 1; i <= retries; i++ {
		log.Printf("Retrying message (attempt %d/%d): %s\n", i, retries, string(msg.Value))
		time.Sleep(retryInterval)

		// 재처리 시도
		success := processMessage(msg)
		if success {
			// 성공 시 Offset Commit
			_, err := consumer.CommitMessage(msg)
			if err != nil {
				log.Printf("Failed to commit offset after retry: %v", err)
			} else {
				log.Println("Offset committed successfully after retry.")
			}
			return
		}
	}

	// 모든 재시도 실패
	log.Printf("Failed to process message after %d attempts: %s\n", retries, string(msg.Value))
}
