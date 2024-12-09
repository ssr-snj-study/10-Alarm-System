package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/google"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	fcmURL      = "https://fcm.googleapis.com/v1/projects/%s/messages:send"
	projectID   = "test-879c6"             // Firebase 프로젝트 ID
	credentials = "serviceAccountKey.json" // 서비스 계정 키 파일 경로
)

type FCMMessage struct {
	Message FCMMessageBody `json:"message"`
}

type FCMMessageBody struct {
	Token        string            `json:"token"`        // 대상 디바이스의 토큰
	Notification *Notification     `json:"notification"` // 알림 내용
	Data         map[string]string `json:"data"`         // 추가 데이터 (옵션)
}

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func getAccessToken() (string, error) {
	// 서비스 계정 키 파일을 사용하여 Google OAuth2 인증
	data, err := ioutil.ReadFile(credentials)
	if err != nil {
		return "", fmt.Errorf("failed to read credentials file: %v", err)
	}

	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		return "", fmt.Errorf("failed to parse credentials: %v", err)
	}

	// 토큰 가져오기
	tokenSource := conf.TokenSource(context.Background())
	token, err := tokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("failed to get token: %v", err)
	}

	return token.AccessToken, nil
}

func SendMessageToFCM(token, title, body string) error {
	// FCM 메시지 생성
	message := FCMMessage{
		Message: FCMMessageBody{
			Token: token,
			Notification: &Notification{
				Title: title,
				Body:  body,
			},
			Data: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// FCM URL 구성
	url := fmt.Sprintf(fcmURL, projectID)

	// Access Token 가져오기
	accessToken, err := getAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %v", err)
	}

	// HTTP 요청 생성
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonMessage))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// HTTP 요청 보내기
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// 응답 처리
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send message: %v, response: %s", resp.Status, string(respBody))
	}

	log.Println("Message sent successfully!")
	return nil
}
