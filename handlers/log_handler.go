package handlers

import (
	"fmt"
	"net/http"

	"goneng_api_api/db"
	"goneng_api_api/models"

	"github.com/labstack/echo/v4"
)

// GetLog POST /post_get_log
func GetLog(c echo.Context) error {
	var req models.LogRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			map[string]string{"message": "요청 형식 오류"})
	}

	// 기본값 처리
	if req.StartDate == "" {
		req.StartDate = "today"
	}
	if req.EndDate == "" {
		req.EndDate = "today"
	}

	query := `
		SELECT
			TO_CHAR(ts, 'YYYY-MM-DD HH24:MI:SS'),
			장치,
			구분,
			내용
		FROM public.로그
		WHERE ts >= $1::date
          AND ts < ($2::date + INTERVAL '1 day')
		ORDER BY ts DESC
	`

	fmt.Printf("Executing query: %s\nWith params: start=%s, end=%s\n", query, req.StartDate, req.EndDate)

	rows, err := db.DB.Query(query, req.StartDate, req.EndDate)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			map[string]string{"message": "DB 조회 오류: " + err.Error()})
	}
	defer rows.Close()

	logs := make([]models.LogEntry, 0, 64)
	for rows.Next() {
		var e models.LogEntry
		if err := rows.Scan(&e.Time, &e.Device, &e.Gubun, &e.Text); err != nil {
			continue
		}
		logs = append(logs, e)
	}

	return c.JSON(http.StatusOK, models.LogResponse{
		Data:    logs,
		Count:   len(logs),
		Message: "ok",
	})
}
