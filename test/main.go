package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	// Kafka Producer 설정
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092", // Kafka 브로커 주소
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	topic := "example-topic"
	message := "Hello, Kafka with confluent-kafka-go/v2!"

	// 비동기 전달 결과 확인
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				fmt.Println(e.String())
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// 메시지 전송
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
	}, nil)

	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	// 메시지 처리가 완료될 때까지 대기
	producer.Flush(15 * 1000)
	fmt.Println("Producer finished successfully")
}

//
//package main
//
//import (
//"log"
//
//"github.com/confluentinc/confluent-kafka-go/v2/kafka"
//)
//
//func main() {
//	// Kafka Consumer 설정
//	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
//		"bootstrap.servers": "localhost:9092", // Kafka 브로커 주소
//		"group.id":          "example-group",  // Consumer 그룹 ID
//		"auto.offset.reset": "earliest",       // 메시지를 처음부터 읽음
//	})
//	if err != nil {
//		log.Fatalf("Failed to create consumer: %s", err)
//	}
//	defer consumer.Close()
//
//	// 구독할 토픽 설정
//	topic := "example-topic"
//	err = consumer.SubscribeTopics([]string{topic}, nil)
//	if err != nil {
//		log.Fatalf("Failed to subscribe to topic: %s", err)
//	}
//
//	log.Println("Consumer started. Waiting for messages...")
//
//	// 메시지 수신
//	for {
//		msg, err := consumer.ReadMessage(-1)
//		if err != nil {
//			log.Printf("Consumer error: %v", err)
//			continue
//		}
//		log.Printf("Received message: %s from topic: %s", string(msg.Value), *msg.TopicPartition.Topic)
//	}
//}
