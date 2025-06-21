package main

import (
	"bytes"
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
	VOICE_FILENAME = "voice.wav"
	VOICEVOX_API_URL = "http://voicevox-engine:50021"
	AUDIO_DEVICE = "plughw:1,0"
)

func alertHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Processing alertHandler")

	if err := processAlert(); err != nil {
		sendErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("Failed to process alert: %v", err))
		return
	}
	sendSuccessResponse(writer, "Completed alertHandler")
}

func processAlert() error {
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
		fmt.Println("Warning: cleanup failed:", err)
	}
	return nil
}

func getSpeechMessage() string {
	hour := time.Now().Hour()
	return fmt.Sprintf("%d時をお知らせします。", hour)
}

func requestAudioQuery(speechMessage string, speaker int) ([]byte, error) {
	requestUrl := fmt.Sprintf("%s/audio_query?text=%s&speaker=%d", VOICEVOX_API_URL, url.QueryEscape(speechMessage), speaker)
	
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
	if err := os.Remove(soundFilePath); err != nil {
		return fmt.Errorf("failed to remove audio file: %w", err)
	}
	return nil
}
