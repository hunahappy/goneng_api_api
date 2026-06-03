package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"goneng_api_api/config"
	"goneng_api_api/db"
	mw "goneng_api_api/middleware"
	"goneng_api_api/models"
)

// Login POST /login
func Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			map[string]string{"message": "요청 형식 오류"})
	}
	if req.Username == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest,
			map[string]string{"message": "아이디와 비밀번호를 입력하세요"})
	}

	var storedPw string
	err := db.DB.QueryRow(
		"SELECT password FROM users WHERE username = $1",
		req.Username,
	).Scan(&storedPw)

	if err != nil || storedPw != req.Password {
		return echo.NewHTTPError(http.StatusUnauthorized,
			map[string]string{"message": "아이디 또는 비밀번호 오류"})
	}

	exp := time.Now().Add(time.Duration(config.App.JWT.ExpireHours) * time.Hour)
	claims := &mw.Claims{
		Username: req.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.App.JWT.Secret))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			map[string]string{"message": "토큰 생성 실패"})
	}

	return c.JSON(http.StatusOK, models.LoginResponse{
		Token:    signed,
		Username: req.Username,
		Message:  "로그인 성공",
	})
}
