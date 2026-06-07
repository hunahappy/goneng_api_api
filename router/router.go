package router

import (
	"net/http"

	"goneng_api_api/handlers"
	mw "goneng_api_api/middleware"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

// Setup Echo 라우터 설정 및 반환
func Setup() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	// ── 미들웨어 ──────────────────────────────────────────────
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet, http.MethodPost,
			http.MethodPut, http.MethodDelete, http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	// ── 공개 라우트 (인증 불필요) ─────────────────────────────
	e.POST("/login", handlers.Login)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// ── 보호 라우트 (JWT 필요) ────────────────────────────────
	api := e.Group("", mw.JWT())

	// 센서 데이터 조회
	api.POST("/get_all_sensors", handlers.GetAllSensors)

	// 에어컨 제어
	api.POST("/post_set_on_aircon", handlers.AirconOn)
	api.POST("/post_set_off_aircon", handlers.AirconOff)
	api.POST("/upload_jpg", handlers.UploadJPG)

	// 로그 조회
	api.POST("/post_get_log", handlers.GetLog)

	return e
}
