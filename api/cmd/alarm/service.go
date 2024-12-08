package alarm

import (
	"api/config"
	"api/model"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SendMsg(c echo.Context) error {
	b := new(model.Message)
	if err := c.Bind(b); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}
	var err error
	// 인증 먼저 시도

	// cache에서 값을 먼저 찾는다.
	userId := checkCacheUser(b)
	if userId < 1 {
		userId, err = checkUser(b)
		if err != nil {
			data := map[string]interface{}{
				"message": err.Error(),
			}
			return c.JSON(http.StatusInternalServerError, data)
		}
		if userId < 1 {
			data := map[string]interface{}{
				"message": errors.New("user_id is not available"),
			}
			return c.JSON(http.StatusInternalServerError, data)
		}
	}
	// 계정별 전송률 제한
	if checkSendOk(userId) {
		//TODO 대기를 할지 다름에 보내달라고 할지 결정하기
	} else {
		// queue에 값 입력
		fmt.Println("insert message queue")
		producer := config.KafkaProducer()
		producer.ProduceMsg(b.Contents)
	}

	response := map[string]interface{}{}

	return c.JSON(http.StatusOK, response)
}
