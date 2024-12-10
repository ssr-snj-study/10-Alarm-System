package config

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

var database *gorm.DB
var stCache StCache
var stKafka StKafka

type StCache struct {
	StCache *redis.Client
	context context.Context
}

type StKafka struct {
	StKafka *kafka.Producer
}

func init() {
	databaseInit()
	startTopic()
	alarmInit()
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

func Cache() StCache {
	return stCache
}

func KafkaProducer() StKafka {
	return stKafka
}

func init() {
	alarmInit()
}

func alarmInit() {
	//connectInfo := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	connectInfo := fmt.Sprintf("%s:%d", "127.0.0.1", 6380)
	stCache.StCache = redis.NewClient(&redis.Options{
		Addr: connectInfo, // Redis 서버 주소
		//Password: os.Getenv("REDIS_PASSWORD"), // 비밀번호가 없다면 빈 문자열
		Password: "snj", // 비밀번호가 없다면 빈 문자열
	})
	stCache.context = context.Background()
}

func (c *StCache) GetRedisByKey(key string) string {
	// 값 가져오기
	val, err := c.StCache.Get(c.context, key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func (c *StCache) InsertRedis(key, value string) {
	// 키-값 설정
	err := c.StCache.Set(c.context, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func (c *StCache) IncrementKey(key string) int {
	keyVal, err := c.StCache.Incr(c.context, key).Result()
	if err != nil {
		return 0
	}
	return int(keyVal)
}

func (c *StCache) IncrementWithTTL(key string, time time.Duration) error {
	_, err := c.StCache.Incr(c.context, key).Result()
	if err != nil {
		return err
	}

	_, err = c.StCache.Expire(c.context, key, time).Result()
	if err != nil {
		return err
	}

	return nil
}

func startTopic() {
	// Kafka Producer 설정
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":   "localhost:9092", // Kafka 브로커 주소
		"api.version.request": false,
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	stKafka.StKafka = producer
}

func (k *StKafka) ProduceMsg(message string) {
	topic := "example-topic"
	// 메시지 전송
	err := k.StKafka.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
	}, nil)

	go func() {
		for e := range k.StKafka.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	k.StKafka.Flush(15 * 1000)
}
