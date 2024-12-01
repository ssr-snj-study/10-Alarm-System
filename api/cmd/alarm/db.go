package alarm

import (
	"api/config"
	"api/model"
	"fmt"
	"strconv"
)

func checkCacheUser(data *model.Message) int {
	cache := config.Cache()
	userId := cache.GetRedisByKey(data.DeviceToken)
	intUserId, _ := strconv.Atoi(userId)
	return intUserId
}

func checkUser(data *model.Message) (int, error) {
	db := config.DB()
	device := &model.Device{}
	fmt.Println("test11 ", data.DeviceToken)
	if res := db.Where("device_token = ?", data.DeviceToken).Find(device); res.Error != nil {
		fmt.Println("test12 ", res.Error)
		return 0, res.Error
	}
	return device.UserId, nil
}

func checkSendOk(id int) bool {
	cache := config.Cache()
	if exists, _ := cache.StCache.Exists(config.Ctx, strconv.Itoa(id)).Result(); exists > 0 {
		return true
	}
	return false
}

func insertMessageQueue() {

}
