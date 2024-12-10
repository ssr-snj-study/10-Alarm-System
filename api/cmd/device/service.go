package device

import (
	"api/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func InsertDevice(c echo.Context) error {
	b := new(model.Device)
	if err := c.Bind(b); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}
	if err := CheckDevice(b); err != nil {
		return err
	}
	_, _ = CreateDevice(b)
	return nil
}
