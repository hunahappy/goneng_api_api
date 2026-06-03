package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 전체 앱 설정
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	MQTT     MQTTConfig     `yaml:"mqtt"`
	JWT      JWTConfig      `yaml:"jwt"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type MQTTConfig struct {
	Broker   string `yaml:"broker"`
	ClientID string `yaml:"goneng_api_api"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

// App 전역 설정 인스턴스
var App Config

// Load config/config.yaml 로드
func Load() {
	data, err := os.ReadFile("goneng_api_api.yaml")
	if err != nil {
		log.Fatalf("[Config] 설정 파일 읽기 오류: %v", err)
	}
	if err := yaml.Unmarshal(data, &App); err != nil {
		log.Fatalf("[Config] 설정 파싱 오류: %v", err)
	}
	if App.JWT.ExpireHours == 0 {
		App.JWT.ExpireHours = 24
	}
	log.Printf("[Config] 로드 완료 - DB: %s:%d/%s, Port: %d",
		App.Database.Host, App.Database.Port, App.Database.Name, App.Server.Port)
}
