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

func CheckDevice(data *model.Device) error {
	db := config.DB()
	device := &model.Device{}
	if res := db.Where("device_token = ?", data.DeviceToken).Find(device); res.Error != nil {
		return res.Error
	}
	return nil
}
