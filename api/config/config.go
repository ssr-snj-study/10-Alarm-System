package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB
var stCache StCache

var Ctx context.Context = context.Background()

type StCache struct {
	StCache *redis.Client
}

func init() {
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

func Cache() StCache {
	return stCache
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
}

func (c *StCache) GetRedisByKey(key string) string {
	// 값 가져오기
	ctx := context.Background()
	val, err := c.StCache.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func (c *StCache) InsertRedis(key, value string) {
	// 키-값 설정
	ctx := context.Background()

	err := c.StCache.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}
