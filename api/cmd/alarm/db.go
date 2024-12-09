package alarm

import (
	"api/config"
	"api/model"
	"strconv"
)

func checkCacheUser(data *model.Message) string {
	cache := config.Cache()
	token := cache.GetRedisByKey(strconv.Itoa(data.UserId))
	return token
}

func inCacheDeviceToken(userId int, token string) {
	cache := config.Cache()
	cache.InsertRedis(strconv.Itoa(userId), token)
}

func checkUser(data *model.Message) (string, error) {
	db := config.DB()
	device := &model.Device{}
	if res := db.Where("user_id = ?", data.UserId).Find(device); res.Error != nil {
		return "", res.Error
	}
	return device.DeviceToken, nil
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
