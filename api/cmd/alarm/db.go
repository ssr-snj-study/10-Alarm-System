package alarm

import (
	"api/config"
	"api/model"
	"fmt"
	"strconv"
	"time"
)

func checkCacheUser(data *model.Message) string {
	cache := config.Cache()
	token := cache.GetRedisByKey(strconv.Itoa(data.UserId))
	return token
}

func inCacheDeviceToken(userId int, token string) {
	cache := config.Cache()
	cache.InsertRedis(strconv.Itoa(userId), token)
	err := cache.IncrementWithTTL(token, time.Duration(1*time.Minute))
	if err != nil {
		fmt.Println(err)
	}
}

func checkUser(data *model.Message) (string, error) {
	db := config.DB()
	device := &model.Device{}
	if res := db.Where("user_id = ?", data.UserId).Find(device); res.Error != nil {
		return "", res.Error
	}
	return device.DeviceToken, nil
}

func getSendAmount(token string) int {
	cache := config.Cache()
	sendCnt := cache.IncrementKey(token)

	return sendCnt
}
