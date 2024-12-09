package main

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"msgWorker/cmd"
	"msgWorker/config"
	"msgWorker/model"
	"time"
)

func main() {

	consumer := config.KafkaConsumer()
	// 메시지 수신
	for {
		msg, err := consumer.StKafka.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v", err)
			continue
		}
		var msgVal model.Message
		_ = json.Unmarshal([]byte(msg.Value), &msgVal)
		log.Printf("Received message: %s from topic: %s", string(msg.Value), *msg.TopicPartition.Topic)
		if err := cmd.SendMessageToFCM(msgVal.DeviceToken, msgVal.Contents, msgVal.Contents); err != nil {
			log.Printf("Error sending FCM message: %v", err)
			log.Println("Message processing failed. Retrying...")
			retryProcessing(consumer.StKafka, &msgVal, msg)
		}
		_, err = consumer.StKafka.CommitMessage(msg)
		if err != nil {
			log.Printf("Failed to commit offset: %v", err)
		} else {
			log.Println("Offset committed successfully.")
		}
	}
}

func retryProcessing(consumer *kafka.Consumer, msgVal *model.Message, msg *kafka.Message) {
	retries := 3                     // 최대 재시도 횟수
	retryInterval := 2 * time.Second // 재시도 간격

	for i := 1; i <= retries; i++ {
		log.Printf("Retrying message (attempt %d/%d): %s\n", i, retries, string(msg.Value))
		time.Sleep(retryInterval)

		// 재처리 시도
		err := cmd.SendMessageToFCM(msgVal.DeviceToken, msgVal.Contents, msgVal.Contents)
		if err == nil {
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
