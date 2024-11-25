package alarm

import (
	"github.com/labstack/echo/v4"
)

func AlarmRoute(e *echo.Echo) {
	urlRoute := e.Group("/api/v1")
	urlRoute.POST("/SendMsg", SendMsg)
}
