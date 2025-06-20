package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Config はアプリケーションの設定を表します
type Config struct {
	Times   []int  `json:"times"`
	Speaker string `json:"speaker"`
}

// PlayRequest は手動再生リクエストを表します
type PlayRequest struct {
	Speaker string `json:"speaker"`
	Hour    int    `json:"hour"`
}

// App はアプリケーションの状態を保持します
type App struct {
	config     Config
	configMux  sync.RWMutex
	configFile string
}

// NewApp は新しいアプリケーションインスタンスを作成します
func NewApp(configFile string) *App {
	return &App{
		configFile: configFile,
		config:     Config{},
	}
}

func main() {
	configFile := os.Getenv("CONFIG_FILE_PATH")
	if configFile == "" {
		configFile = "config.json" // デフォルト値
	}
	app := NewApp(configFile)

	// 既存の設定を読み込み
	if err := app.loadConfig(); err != nil {
		log.Printf("Warning: Could not load config: %v", err)
	}

	setupRoutes(app)
	app.startScheduler()
	app.startServer()
}

// setupRoutes は HTTP ルートを設定します
func setupRoutes(app *App) {
	http.HandleFunc("/api/save", app.saveConfigHandler)
	http.HandleFunc("/api/config", app.getConfigHandler)
	http.HandleFunc("/api/play", app.playHandler)
	http.HandleFunc("/", app.serveFrontendHandler)
}

// startScheduler は時報スケジューラーを開始します
func (app *App) startScheduler() {
	go app.scheduleTimeBell()
}

// startServer は HTTP サーバーを開始します
func (app *App) startServer() {
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// saveConfigHandler は設定保存リクエストを処理します
func (app *App) saveConfigHandler(w http.ResponseWriter, r *http.Request) {
	app.logRequest("saveConfigHandler")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newConfig Config
	if err := app.decodeRequestBody(r, &newConfig); err != nil {
		log.Printf("Failed to parse JSON: %v", err)
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	app.configMux.Lock()
	app.config = newConfig
	app.configMux.Unlock()

	if err := app.saveConfigToFile(); err != nil {
		log.Printf("Failed to save config: %v", err)
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// getConfigHandler は設定取得リクエストを処理します
func (app *App) getConfigHandler(w http.ResponseWriter, r *http.Request) {
	app.configMux.RLock()
	config := app.config
	app.configMux.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// serveFrontendHandler はフロントエンドの HTML ファイルを配信します
func (app *App) serveFrontendHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "frontend/index.html")
}

// scheduleTimeBell は時報スケジューラーを実行します
func (app *App) scheduleTimeBell() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		currentHour := time.Now().Hour()
		currentMinute := time.Now().Minute()

		// 正時のみ実行（分が0の時）
		if currentMinute != 0 {
			continue
		}

		app.configMux.RLock()
		times := app.config.Times
		speaker := app.config.Speaker
		app.configMux.RUnlock()

		for _, hour := range times {
			if hour == currentHour {
				log.Printf("Playing time bell for %d:00", currentHour)
				go app.playVoice(speaker, currentHour)
				break
			}
		}
	}
}

// playVoice は VOICEVOX API を使用して音声を生成・再生します
func (app *App) playVoice(speaker string, hour int) {
	logRequest("playVoice")
	if speaker == "" {
		log.Println("No speaker configured")
		return
	}

	text := fmt.Sprintf("%d時をお知らせします", hour)
	queryURL := fmt.Sprintf("http://voicevox-engine:50021/audio_query?speaker=%s&text=%s", speaker, url.QueryEscape(text))
	queryResp, err := http.Post(queryURL, "application/json", nil)
	if err != nil {
		log.Printf("VOICEVOX audio_query API call failed: %v", err)
		return
	}
	defer queryResp.Body.Close()
	if queryResp.StatusCode != http.StatusOK {
		log.Printf("VOICEVOX audio_query API returned status: %d", queryResp.StatusCode)
		body, _ := io.ReadAll(queryResp.Body)
		log.Printf("VOICEVOX audio_query API error response: %s", string(body))
		return
	}

	queryBody, err := io.ReadAll(queryResp.Body)
	if err != nil {
		log.Printf("Failed to read audio_query response: %v", err)
		return
	}

	synthURL := fmt.Sprintf("http://voicevox-engine:50021/synthesis?speaker=%s", speaker)
	synthResp, err := http.Post(synthURL, "application/json", bytes.NewReader(queryBody))
	if err != nil {
		log.Printf("VOICEVOX synthesis API call failed: %v", err)
		return
	}
	defer synthResp.Body.Close()
	if synthResp.StatusCode != http.StatusOK {
		log.Printf("VOICEVOX synthesis API returned status: %d", synthResp.StatusCode)
		body, _ := io.ReadAll(synthResp.Body)
		log.Printf("VOICEVOX synthesis API error response: %s", string(body))
		return
	}

	filename := fmt.Sprintf("/tmp/output_%d.wav", time.Now().UnixNano())
	if err := app.saveResponseToFile(synthResp.Body, filename); err != nil {
		log.Printf("Failed to save audio file: %v", err)
		return
	}

	// 複数の音声再生方法を試す
	if err := app.playAudio(filename); err != nil {
		log.Printf("Failed to play audio: %v", err)
	}

	if err := os.Remove(filename); err != nil {
		log.Printf("Failed to remove temporary file: %v", err)
	}
}

// playHandler は手動再生リクエストを処理します
func (app *App) playHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Invalid request method"})
		return
	}
	var request PlayRequest
	if err := app.decodeRequestBody(r, &request); err != nil {
		log.Printf("Failed to parse play request: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Failed to parse JSON"})
		return
	}
	go app.playVoice(request.Speaker, request.Hour)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "playing"})
}

// ユーティリティ関数

// logRequest は受信したリクエストをログに記録します
func (app *App) logRequest(handlerName string) {
	log.Printf("%s: Request received", handlerName)
}

// decodeRequestBody は JSON リクエストボディをデコードします
func (app *App) decodeRequestBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

// loadConfig は設定をファイルから読み込みます
func (app *App) loadConfig() error {
	file, err := os.Open(app.configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	app.configMux.Lock()
	err = decoder.Decode(&app.config)
	app.configMux.Unlock()
	return err
}

// saveConfigToFile は設定をファイルに保存します
func (app *App) saveConfigToFile() error {
	file, err := os.Create(app.configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	app.configMux.RLock()
	err = encoder.Encode(app.config)
	app.configMux.RUnlock()
	return err
}

// saveResponseToFile は HTTP レスポンスボディをファイルに保存します
func (app *App) saveResponseToFile(body io.Reader, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, body)
	return err
}

// playAudio は複数の方法で音声ファイルを再生を試みます
func (app *App) playAudio(filename string) error {
	// 試す音声再生コマンドのリスト
	commands := [][]string{
		{"aplay", filename},                    // ALSA
		{"paplay", filename},                   // PulseAudio
		{"ffplay", "-nodisp", "-autoexit", filename}, // FFmpeg
		{"mpv", "--no-video", "--quiet", filename},   // mpv
	}

	for _, cmd := range commands {
		if app.tryPlayCommand(cmd) {
			log.Printf("Successfully played audio using: %s", cmd[0])
			return nil
		}
	}

	return fmt.Errorf("all audio playback methods failed")
}

// tryPlayCommand は指定されたコマンドで音声再生を試みます
func (app *App) tryPlayCommand(cmdArgs []string) bool {
	// コマンドが存在するかチェック
	if _, err := exec.LookPath(cmdArgs[0]); err != nil {
		log.Printf("Command %s not found: %v", cmdArgs[0], err)
		return false
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to play with %s: %v, output: %s", cmdArgs[0], err, string(output))
		return false
	}

	return true
}
