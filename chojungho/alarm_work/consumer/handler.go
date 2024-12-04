package consumer

import (
	"log"

	"github.com/streadway/amqp"
)

func HandleMessages(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		// 메시지 처리 로직 추가
	}
}
