package handlers

import (
	"encoding/json"
	"fmt"
	"goneng_api_api/db"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

// uploadDir : 저장 디렉토리 (서버 실행 위치 기준 상대경로)
// 필요시 환경변수 또는 config.yaml 에서 주입하도록 수정 가능
const uploadDir = "./upload_jpg"

// UploadJPG  POST /upload_jpg
//
// Content-Type : image/jpeg
// Body         : raw JPEG binary (ESP32 프레임버퍼 직접 전송)
//
// 성공 응답 (200):
//
//	{ "filename": "20060102_150405_000.jpg", "size": 123456 }
func UploadJPG(c echo.Context) error {
	// ── 디렉토리 보장 ──────────────────────────────────────────
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"error": "cannot create upload directory"})
	}

	// ── 파일명: 타임스탬프 + 밀리초 (초 단위 충돌 방지) ──────
	now := time.Now()
	filename := fmt.Sprintf("%s_%03d.jpg",
		now.Format("20060102_150405"),
		now.UnixMilli()%1000,
	)
	savePath := filepath.Join(uploadDir, filename)

	// ── 파일 생성 ─────────────────────────────────────────────
	f, err := os.Create(savePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"error": "cannot create file"})
	}
	defer f.Close()

	// ── 요청 바디를 파일로 스트리밍 복사 (메모리 버퍼 최소화) ─
	written, err := io.Copy(f, c.Request().Body)
	if err != nil {
		// 부분 기록된 파일 정리
		_ = os.Remove(savePath)
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"error": "failed to write file"})
	}

	if written == 0 {
		_ = os.Remove(savePath)
		return c.JSON(http.StatusBadRequest,
			map[string]string{"error": "empty body"})
	}

	c.Logger().Infof("jpg saved: %s (%d bytes)", filename, written)

	content := map[string]interface{}{
		"filename": filename,
	}

	// JSON으로 변환
	jsonData, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err)
	}

	_, _ = db.DB.Exec(
		`INSERT INTO 로그 (장치, 구분, 내용, 토픽) VALUES ($1, $2, $3, $4)`,
		"camera", "jpg", jsonData, "upload_jpg",
	)

	return c.JSON(http.StatusOK, map[string]any{
		"filename": filename,
		"size":     written,
	})
}
