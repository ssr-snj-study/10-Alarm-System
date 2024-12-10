package cmd

import (
	"fmt"
	"msgWorker/config"
	"msgWorker/model"
)

func CreateMsg(data *model.ResponseData) {
	db := config.DB()
	msg := &model.Message{
		UserId:   data.UserId,
		Contents: data.Contents,
		SendTime: data.SendTime,
	}
	if err := db.Create(msg).Error; err != nil {
		fmt.Println(err)
	}
}
