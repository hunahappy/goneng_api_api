package handlers

import (
	"log"
	"net/http"

	"goneng_api_api/mqtt"

	"github.com/labstack/echo/v4"
)

// AirconOn POST /post_set_on_aircon
func AirconOn(c echo.Context) error {
	return airconCmd(c, "켜기")
}

// AirconOff POST /post_set_off_aircon
func AirconOff(c echo.Context) error {
	return airconCmd(c, "끄기")
}

func airconCmd(c echo.Context, action string) error {
	username, _ := c.Get("username").(string)

	if err := mqtt.SendAircon(action); err != nil {
		log.Printf("[Control] 에어컨 %s MQTT 실패: %v", action, err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			map[string]string{"message": "MQTT 전송 실패: " + err.Error()})
	}

	log.Printf("[Control] 에어컨 %s - 사용자: %s", action, username)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "에어컨 " + action + " 명령 전송 완료",
		"action":  action,
	})
}
