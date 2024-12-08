package cmd

import (
	"api/cmd/alarm"
	"api/cmd/device"
	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo) {
	urlRoute := e.Group("/api/v1")
	urlRoute.POST("/SendMsg", alarm.SendMsg)
	urlRoute.POST("/SetDevice", device.InsertDevice)
}
