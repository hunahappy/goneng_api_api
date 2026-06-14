package handlers

import (
	"fmt"
	"net/http"

	"goneng_api_api/db"
	"goneng_api_api/models"

	"github.com/labstack/echo/v4"
)

// GetSensor POST /post_get_sensor/:type
func GetAllSensors(c echo.Context) error {
	var req models.SensorRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "요청 형식 오류"})
	}

	// 기본값 처리
	if req.StartDate == "" {
		req.StartDate = "today"
	}
	if req.EndDate == "" {
		req.EndDate = "today"
	}
	if req.Unit == "" {
		req.Unit = "1"
	}

	// JSONB 필드에서 각 센서 값을 추출하여 시간 단위로 집계
	// 내용 필드: {"EC": 0.39, "수온": 26, "습도": 47.05, "온도": 27.55, "조도": 25.83}
	query := `
        SELECT
			TO_TIMESTAMP(
					FLOOR(EXTRACT(EPOCH FROM ts) / ($3 * 60)) * ($3 * 60)
				) AS bucket,				   
            AVG((내용->>'EC')::numeric) AS ec,
            AVG((내용->>'조도')::numeric) AS lux,
            AVG((내용->>'온도')::numeric) AS temp,
            AVG((내용->>'수온')::numeric) AS wt,
            AVG((내용->>'습도')::numeric) AS humi
        FROM public.로그
        WHERE 구분 = 'data' AND 토픽 = 'goneng/farm1/data/sensor/thes'
		  AND ts >= $1::date
          AND ts < ($2::date + INTERVAL '1 day')
        GROUP BY bucket
        ORDER BY bucket ASC
    `
	fmt.Printf("Executing query: %s\nWith params: start=%s, end=%s, unit=%s\n", query, req.StartDate, req.EndDate, req.Unit)

	rows, err := db.DB.Query(query, req.StartDate, req.EndDate, req.Unit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "DB 조회 오류: " + err.Error()})
	}
	defer rows.Close()

	// 결과를 담을 구조체 (모든 데이터를 포함하는 Map 형태)
	var results []map[string]interface{}

	for rows.Next() {
		var time string
		var ec, lux, temp, wt, humi float64

		// Scan 시 SQL에서 집계된 순서대로 매핑
		err := rows.Scan(&time, &ec, &lux, &temp, &wt, &humi)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"time": time,
			"ec":   ec,
			"lux":  lux,
			"temp": temp,
			"wt":   wt,
			"humi": humi,
		})

		// fmt.Printf("Fetched row: time=%s, ec=%.2f, lux=%.2f, temp=%.2f, wt=%.2f, humi=%.2f\n", time, ec, lux, temp, wt, humi)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    results,
		"count":   len(results),
		"message": "ok",
	})
}
