package model

import "time"

type UserTb struct {
	UserId      int       `json:"userId,omitempty"`
	Email       string    `json:"email,omitempty"`
	CountryCode int       `json:"country_Code,omitempty"`
	PhoneNumber int       `json:"phone_Number,omitempty"`
	CreatedArt  time.Time `json:"created_Art"`
}

func (UserTb) TableName() string {
	return "user_tb"
}

type Device struct {
	Id             int       `json:"id,omitempty"`
	DeviceToken    string    `json:"device_Token,omitempty"`
	UserId         int       `json:"user_Id,omitempty"`
	LastLoggedInAt time.Time `json:"last_Logged_In_At"`
}

func (Device) TableName() string {
	return "device"
}

type Message struct {
	Id          int       `json:"id,omitempty"`
	DeviceToken string    `json:"device_Token,omitempty"`
	UserId      int       `json:"user_Id,omitempty"`
	SendTime    time.Time `json:"send_time"`
	Contents    string    `json:"contents"`
	Receiver    string    `json:"receiver"`
}

func (Message) TableName() string {
	return "message"
}
