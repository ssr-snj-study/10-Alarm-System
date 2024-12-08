package device

import (
	"api/config"
	"api/model"
	"time"
)

func CreateDevice(data *model.Device) (int, error) {
	db := config.DB()
	device := &model.Device{
		UserId:         data.UserId,
		DeviceToken:    data.DeviceToken,
		LastLoggedInAt: time.Now(),
	}
	if err := db.Create(&device).Error; err != nil {
		return 0, err
	}
	return device.UserId, nil
}
