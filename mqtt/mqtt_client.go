package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"goneng_api_api/config"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// Client 전역 MQTT 클라이언트
var Client paho.Client

type controlPayload struct {
	Device  string                 `json:"장치"`
	Gubun   string                 `json:"구분"`
	Content map[string]interface{} `json:"내용"`
}

// Init MQTT 클라이언트 초기화 및 브로커 연결
func Init() {
	cfg := config.App.MQTT
	opts := paho.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetClientID(cfg.ClientID)

	if cfg.UserName != "" {
		opts.SetUsername(cfg.UserName)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}

	opts.SetConnectTimeout(5 * time.Second)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(true)

	opts.OnConnect = func(_ paho.Client) {
		log.Printf("[MQTT] 브로커 연결 완료: %s", cfg.Broker)
	}
	opts.OnConnectionLost = func(_ paho.Client, err error) {
		log.Printf("[MQTT] 연결 끊김 (자동 재연결): %v", err)
	}

	Client = paho.NewClient(opts)
	token := Client.Connect()
	if token.WaitTimeout(5 * time.Second) {
		if err := token.Error(); err != nil {
			log.Printf("[MQTT] 초기 연결 경고: %v (자동 재연결 대기)", err)
			return
		}
	}
	log.Println("[MQTT] 초기화 완료")
}

// SendAircon 에어컨 제어 메세지 발행
// action: "켜기" | "끄기"
func SendAircon(action string) error {
	if Client == nil {
		return fmt.Errorf("MQTT 클라이언트 미초기화")
	}
	if !Client.IsConnected() {
		return fmt.Errorf("MQTT 브로커 미연결")
	}

	payload := controlPayload{
		Device: "server_api_api",
		Gubun:  "control",
		Content: map[string]interface{}{
			"제어": action,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("JSON 직렬화 실패: %w", err)
	}

	const topic = "goneng/farm1/control/ir/air"
	token := Client.Publish(topic, 1, false, data)
	if !token.WaitTimeout(3 * time.Second) {
		return fmt.Errorf("MQTT 발행 타임아웃")
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("MQTT 발행 오류: %w", err)
	}

	log.Printf("[MQTT] 발행 완료 topic=%s payload=%s", topic, data)
	return nil
}
