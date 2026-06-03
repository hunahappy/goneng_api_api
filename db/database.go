package db

import (
	"database/sql"
	"fmt"
	"log"

	"goneng_api_api/config"

	_ "github.com/lib/pq"
)

// DB 전역 PostgreSQL 연결
var DB *sql.DB

// Init DB 연결 및 스키마 초기화
func Init() {
	cfg := config.App.Database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("[DB] Open 오류: %v", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	if err = DB.Ping(); err != nil {
		log.Fatalf("[DB] Ping 오류: %v", err)
	}
	log.Printf("[DB] PostgreSQL 연결 성공: %s:%d/%s", cfg.Host, cfg.Port, cfg.Name)

	//initSchema()
}

func initSchema() {
	sqls := []string{
		// 사용자 테이블
		`CREATE TABLE IF NOT EXISTS users (
			id         SERIAL PRIMARY KEY,
			username   VARCHAR(50)  UNIQUE NOT NULL,
			password   VARCHAR(255) NOT NULL,
			created_at TIMESTAMP   DEFAULT NOW()
		)`,
		// 센서 시계열 테이블
		`CREATE TABLE IF NOT EXISTS sensor_log (
			id          SERIAL PRIMARY KEY,
			log_time    TIMESTAMP       DEFAULT NOW(),
			device      VARCHAR(50),
			gubun       VARCHAR(50),
			sensor_type VARCHAR(20)     NOT NULL,
			value       DOUBLE PRECISION NOT NULL
		)`,
		// 일반 로그 테이블
		`CREATE TABLE IF NOT EXISTS log_table (
			id       SERIAL PRIMARY KEY,
			log_time TIMESTAMP DEFAULT NOW(),
			device   VARCHAR(50),
			gubun    VARCHAR(50),
			content  TEXT
		)`,
		// 인덱스
		`CREATE INDEX IF NOT EXISTS idx_sensor_time ON sensor_log(log_time)`,
		`CREATE INDEX IF NOT EXISTS idx_sensor_type ON sensor_log(sensor_type, log_time)`,
		`CREATE INDEX IF NOT EXISTS idx_log_time    ON log_table(log_time)`,
		// 기본 관리자 계정 (비밀번호: admin123)
		`INSERT INTO users (username, password)
		 VALUES ('admin', 'admin123')
		 ON CONFLICT (username) DO NOTHING`,
		// 샘플 센서 데이터
		`INSERT INTO sensor_log (log_time, device, gubun, sensor_type, value)
		 SELECT NOW() - (n * INTERVAL '1 minute'),
		        '센서', '데이터',
		        CASE (n % 5)
		          WHEN 0 THEN 'ec'
		          WHEN 1 THEN 'lux'
		          WHEN 2 THEN 'temp'
		          WHEN 3 THEN 'wt'
		          ELSE 'humi'
		        END,
		        ROUND((RANDOM() * 100)::NUMERIC, 2)
		 FROM generate_series(1, 200) n
		 ON CONFLICT DO NOTHING`,
	}

	for _, s := range sqls {
		if _, err := DB.Exec(s); err != nil {
			log.Printf("[DB] 스키마 경고: %v", err)
		}
	}
	log.Println("[DB] 스키마 초기화 완료")
}
