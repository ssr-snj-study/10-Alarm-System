package consumer

import (
	"context"
	"encoding/json"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/streadway/amqp"
	"google.golang.org/api/option"
	"log"
)

// NotificationMessage 메시지 데이터 구조체
type NotificationMessage struct {
	UserID         int    `json:"user_id"`
	Email          string `json:"email"`
	CountryCode    int    `json:"country_code"`
	PhoneNumber    string `json:"phone_number"`
	CreatedAt      string `json:"created_at"`
	DeviceToken    string `json:"device_token"`
	DeviceType     int    `json:"device_type"`
	LastLoggedInAt string `json:"last_logged_in_at"`
	Message        string `json:"message"`
}

// Firebase Admin App 초기화
func initFirebaseApp() *messaging.Client {
	opt := option.WithCredentialsFile("cjh-alarm-app-firebase-adminsdk-7t0f3-e0f83a6bf7.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase App: %v", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Messaging Client: %v", err)
	}

	return client
}

// FCM 알림 전송 함수
func sendNotification(client *messaging.Client, token, title, body string) {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: token,
	}

	// 메시지 전송
	response, err := client.Send(context.Background(), message)
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		return
	}
	log.Printf("Successfully sent notification: %s", response)
}

// HandleMessages RabbitMQ 메시지 처리 함수
func HandleMessages(msgs <-chan amqp.Delivery) {
	// Firebase 클라이언트 초기화
	fcmClient := initFirebaseApp()

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)

		// RabbitMQ 메시지를 NotificationMessage로 변환
		var notification NotificationMessage
		err := json.Unmarshal(d.Body, &notification)
		if err != nil {
			log.Printf("Failed to parse message: %v", err)
			d.Nack(false, false) // 메시지 재처리 안 함
			continue
		}

		// FCM 알림 제목과 내용 생성
		title := "alarm project test"
		body := notification.Message

		// Firebase FCM 알림 전송
		sendNotification(
			fcmClient,
			notification.DeviceToken,
			title,
			body,
		)

		// 메시지 성공적으로 처리 -> RabbitMQ에서 ACK
		d.Ack(false)
	}
}
