package models

type LogEntry struct {
	Time   string `json:"time"`
	Device string `json:"device"`
	Gubun  string `json:"gubun"`
	Text   string `json:"text"`
}

type LogRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Gubun     string `json:"gubun"`
	Device    string `json:"device"`
	Text      string `json:"text"`
}

type LogResponse struct {
	Data    []LogEntry `json:"data"`
	Count   int        `json:"count"`
	Message string     `json:"message"`
}
