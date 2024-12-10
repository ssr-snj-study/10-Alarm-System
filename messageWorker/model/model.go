package model

import "time"

type Message struct {
	Id       int       `json:"id,omitempty"`
	UserId   int       `json:"user_id,omitempty"`
	SendTime time.Time `json:"send_time,omitempty"`
	Contents string    `json:"contents,omitempty"`
	Receiver string    `json:"receiver,omitempty"`
}

func (Message) TableName() string {
	return "message"
}

type ResponseData struct {
	DeviceToken string    `json:"device_token,omitempty"`
	Contents    string    `json:"contents,omitempty"`
	Title       string    `json:"title,omitempty"`
	UserId      int       `json:"user_id,omitempty"`
	SendTime    time.Time `json:"send_time,omitempty"`
}
