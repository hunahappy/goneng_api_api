package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"goneng_api_api/config"
	"goneng_api_api/db"
	mw "goneng_api_api/middleware"
	"goneng_api_api/models"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
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

	content := map[string]interface{}{
		"user_id": req.Username,
	}

	// JSON으로 변환
	jsonData, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err)
	}

	_, _ = db.DB.Exec(
		`INSERT INTO 로그 (장치, 구분, 내용, 토픽) VALUES ($1, $2, $3, $4)`,
		"api_api", "login", jsonData, "login",
	)

	return c.JSON(http.StatusOK, models.LoginResponse{
		Token:    signed,
		Username: req.Username,
		Message:  "로그인 성공",
	})
}
