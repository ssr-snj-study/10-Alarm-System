package main

import (
	"alarm_work/config"
	"alarm_work/consumer"
	"alarm_work/logger"
	"alarm_work/rabbitmq"
	"log"
)

func main() {
	// 1. 환경 변수 로드
	config.LoadConfig()
	rabbitmqURL := config.GetEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

	// 2. 로깅 초기화
	logger.InitLogger()

	log.Println("Application starting...")

	// 3. RabbitMQ 연결
	conn, ch, err := rabbitmq.Connect(rabbitmqURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	// 4. 큐 선언
	q, err := ch.QueueDeclare(
		"Android",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// RabbitMQ 메시지 소비
	msgs, err := ch.Consume(
		q.Name,
		"",
		false, // 수동 ACK
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// 6. 메시지 처리 시작
	log.Println("Waiting for messages...")
	consumer.HandleMessages(msgs)
}
