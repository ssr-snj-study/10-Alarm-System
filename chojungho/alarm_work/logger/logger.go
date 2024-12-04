package logger

import (
	"log"
)

func InitLogger() {
	// 기본 로깅 설정: 콘솔 출력
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Logger initialized for console output")
}
