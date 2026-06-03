package models

import "time"

type SensorData struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

type SensorRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Unit      string `json:"unit"`   // "1" | "10" | "60"
	Device    string `json:"device"`
	Gubun     string `json:"gubun"`
}

type SensorResponse struct {
	Data    []SensorData `json:"data"`
	Count   int          `json:"count"`
	Message string       `json:"message"`
}
