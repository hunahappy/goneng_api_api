package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"goneng_api_api/db"
	"goneng_api_api/models"
)

// GetLog POST /post_get_log
func GetLog(c echo.Context) error {
	var req models.LogRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			map[string]string{"message": "요청 형식 오류"})
	}

	query := `
		SELECT
			TO_CHAR(log_time, 'YYYY-MM-DD HH24:MI:SS'),
			COALESCE(device, ''),
			COALESCE(gubun, ''),
			COALESCE(content, '')
		FROM log_table
		WHERE 1=1
	`
	args := []interface{}{}
	idx := 1

	if req.Gubun != "" {
		query += fmt.Sprintf(" AND gubun = $%d", idx)
		args = append(args, req.Gubun)
		idx++
	}
	if req.Device != "" {
		query += fmt.Sprintf(" AND device = $%d", idx)
		args = append(args, req.Device)
		idx++
	}
	if req.Text != "" {
		query += fmt.Sprintf(" AND content ILIKE $%d", idx)
		args = append(args, "%"+req.Text+"%")
		idx++
	}

	query += " ORDER BY log_time DESC LIMIT 1000"

	rows, err := db.DB.Query(query, args...)
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
