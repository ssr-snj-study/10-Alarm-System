package model

import "time"

type Message struct {
	Id          int       `json:"id,omitempty"`
	DeviceToken string    `json:"device_token,omitempty"`
	UserId      int       `json:"user_id,omitempty"`
	SendTime    time.Time `json:"send_time,omitempty"`
	Contents    string    `json:"contents,omitempty"`
	Receiver    string    `json:"receiver,omitempty"`
}
