package config

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var stKafka StKafka
var database *gorm.DB

type StKafka struct {
	StKafka *kafka.Consumer
}

func KafkaConsumer() StKafka {
	return stKafka
}

func init() {
	consumeTopic()
	databaseInit()
}

func databaseInit() {
	var e error
	//host := os.Getenv("DB_HOST")
	//user := os.Getenv("DB_USER")
	//password := os.Getenv("DB_PASSWORD")
	//dbName := os.Getenv("DB_NAME")
	//port := os.Getenv("DB_PORT")
	host := "127.0.0.1"
	user := "snj"
	password := "snj"
	dbName := "snj_db"
	port := 5432

	connectInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", host, user, password, dbName, port)
	database, e = gorm.Open(postgres.Open(connectInfo), &gorm.Config{})
	if e != nil {
		panic(e)
	}
	//sqlSet, err := database.DB()
	//if err != nil {
	//	panic("failed to get database")
	//}
	//sqlSet.SetConnMaxLifetime(time.Hour)
	//sqlSet.SetMaxOpenConns(50)

}

func DB() *gorm.DB {
	return database
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
