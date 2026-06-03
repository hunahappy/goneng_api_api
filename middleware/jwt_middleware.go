package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"goneng_api_api/config"
)

// Claims JWT 클레임 구조체
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWT Echo JWT 인증 미들웨어 반환
func JWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized,
					map[string]string{"message": "Authorization 헤더 없음"})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				return echo.NewHTTPError(http.StatusUnauthorized,
					map[string]string{"message": "Authorization 형식 오류 (Bearer <token>)"})
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(
				parts[1], claims,
				func(t *jwt.Token) (interface{}, error) {
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, echo.NewHTTPError(http.StatusUnauthorized,
							"잘못된 서명 방식")
					}
					return []byte(config.App.JWT.Secret), nil
				},
			)
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized,
					map[string]string{"message": "유효하지 않은 토큰"})
			}

			c.Set("username", claims.Username)
			return next(c)
		}
	}
}
