package main

import (
	"fmt"
	"log"

	"goneng_api_api/config"
	"goneng_api_api/db"
	"goneng_api_api/mqtt"
	"goneng_api_api/router"
)

func main() {
	log.Println("========================================")
	log.Println("  goneng_api_api 서버 시작")
	log.Println("========================================")

	// 1. 설정 파일 로드
	config.Load()

	// 2. PostgreSQL 연결
	db.Init()

	// 3. MQTT 클라이언트 초기화
	mqtt.Init()

	// 4. Echo 라우터 설정
	e := router.Setup()

	// 5. 서버 시작
	addr := fmt.Sprintf(":%d", config.App.Server.Port)
	log.Printf("[Server] HTTP 서버 시작: http://0.0.0.0%s", addr)
	if err := e.Start(addr); err != nil {
		log.Fatalf("[Server] 서버 오류: %v", err)
	}
}
