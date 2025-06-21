package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	STORAGE_PATH = "/app/storage"
	VOICE_FILENAME = "voice.wav"
	CONFIG_FILENAME = "config.json"
	VOICEVOX_API_URL = "http://voicevox-engine:50021"
	SERVER_PORT = ":8080"
	AUDIO_DEVICE = "plughw:1,0"
)

// Config represents the timebell configuration
type Config struct {
	HourList []string `json:"hourList"`
	StyleID  string   `json:"styleId"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// corsMiddleware handles CORS headers for all requests
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "*")

		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusOK)
			return
		}
		next(writer, request)
	}
}

// sendErrorResponse sends a standardized error response
func sendErrorResponse(writer http.ResponseWriter, statusCode int, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	response := APIResponse{
		Status: "error",
		Error:  message,
	}
	json.NewEncoder(writer).Encode(response)
}

// sendSuccessResponse sends a standardized success response
func sendSuccessResponse(writer http.ResponseWriter, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	response := APIResponse{
		Status:  "success",
		Message: message,
	}
	json.NewEncoder(writer).Encode(response)
}

func main() {
	setupRoutes()
	startServer()
}

func setupRoutes() {
	http.HandleFunc("/api/announce", corsMiddleware(announceHandler))
	http.HandleFunc("/api/config", corsMiddleware(configHandler))
}

func startServer() {
	fmt.Printf("Starting server on port %s\n", SERVER_PORT)
	if err := http.ListenAndServe(SERVER_PORT, nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}

func announceHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Processing announce request")

	if err := processTimeAnnouncement(); err != nil {
		sendErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("Failed to process announcement: %v", err))
		return
	}

	sendSuccessResponse(writer, "Time announcement completed successfully")
}

func processTimeAnnouncement() error {
	speechMessage := getSpeechMessage()
	
	audioQuery, err := requestAudioQuery(speechMessage, 1)
	if err != nil {
		return fmt.Errorf("audio query failed: %w", err)
	}

	audioData, err := requestSynthesis(audioQuery, 1)
	if err != nil {
		return fmt.Errorf("synthesis failed: %w", err)
	}

	if err := saveAudioFile(audioData); err != nil {
		return fmt.Errorf("save audio failed: %w", err)
	}

	if err := playAudioFile(); err != nil {
		return fmt.Errorf("play audio failed: %w", err)
	}

	if err := cleanupAudioFile(); err != nil {
		fmt.Printf("Warning: failed to cleanup audio file: %v\n", err)
	}

	return nil
}

func getSpeechMessage() string {
	// 現在時刻を取得する
	hour := time.Now().Hour()
	return fmt.Sprintf("%d時になりました", hour)
}

func requestAudioQuery(speechMessage string, speaker int) ([]byte, error) {
	requestUrl := fmt.Sprintf("%s/audio_query?text=%s&speaker=%d", 
		VOICEVOX_API_URL, url.QueryEscape(speechMessage), speaker)
	
	response, err := http.Post(requestUrl, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make audio query request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("audio query API returned status: %s", response.Status)
	}
	
	return io.ReadAll(response.Body)
}

func requestSynthesis(audioQuery []byte, speaker int) ([]byte, error) {
	requestUrl := fmt.Sprintf("%s/synthesis?speaker=%d", VOICEVOX_API_URL, speaker)
	
	response, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(audioQuery))
	if err != nil {
		return nil, fmt.Errorf("failed to make synthesis request: %w", err)
	}
	defer response.Body.Close()
	
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("synthesis API returned status: %s", response.Status)
	}
	
	return io.ReadAll(response.Body)
}

func saveAudioFile(data []byte) error {
	if err := ensureStorageDirectory(); err != nil {
		return fmt.Errorf("failed to ensure storage directory: %w", err)
	}
	
	soundFilePath := filepath.Join(STORAGE_PATH, VOICE_FILENAME)
	if err := os.WriteFile(soundFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write audio file: %w", err)
	}
	
	return nil
}

func playAudioFile() error {
	soundFilePath := filepath.Join(STORAGE_PATH, VOICE_FILENAME)
	cmd := exec.Command("aplay", "-D", AUDIO_DEVICE, soundFilePath)
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to play audio file: %w", err)
	}
	
	return nil
}

func cleanupAudioFile() error {
	soundFilePath := filepath.Join(STORAGE_PATH, VOICE_FILENAME)
	return os.Remove(soundFilePath)
}

func ensureStorageDirectory() error {
	return os.MkdirAll(STORAGE_PATH, 0755)
}

func configHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Processing config request")

	body, err := io.ReadAll(request.Body)
	if err != nil {
		sendErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Failed to read request body: %v", err))
		return
	}
	defer request.Body.Close()

	// Validate JSON format
	var config Config
	if err := json.Unmarshal(body, &config); err != nil {
		sendErrorResponse(writer, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := saveConfigFile(body); err != nil {
		sendErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("Failed to save config: %v", err))
		return
	}

	sendSuccessResponse(writer, "Configuration saved successfully")
}

func saveConfigFile(data []byte) error {
	if err := ensureStorageDirectory(); err != nil {
		return fmt.Errorf("failed to ensure storage directory: %w", err)
	}

	configFilePath := filepath.Join(STORAGE_PATH, CONFIG_FILENAME)
	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
