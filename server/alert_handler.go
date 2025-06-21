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
)

func alertHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("{comment}")

	if err := processAlert(); err != nil {
		sendErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("{comment}", err))
		return
	}
	sendSuccessResponse(writer, "{comment}")
}

func processAlert() error {
	speechMessage := getSpeechMessage()
	
	audioQuery, err := requestAudioQuery(speechMessage, 1)
	if err != nil {
		return fmt.Errorf("{comment}", err)
	}

	audioData, err := requestSynthesis(audioQuery, 1)
	if err != nil {
		return fmt.Errorf("{comment}", err)
	}

	if err := saveAudioFile(audioData); err != nil {
		return fmt.Errorf("{comment}", err)
	}

	if err := playAudioFile(); err != nil {
		return fmt.Errorf("{comment}", err)
	}

	if err := cleanupAudioFile(); err != nil {
		fmt.Println("{comment}")
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
		return nil, fmt.Errorf("{comment}", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("{comment}", response.Status)
	}
	return io.ReadAll(response.Body)
}

func requestSynthesis(audioQuery []byte, speaker int) ([]byte, error) {
	requestUrl := fmt.Sprintf("%s/synthesis?speaker=%d", VOICEVOX_API_URL, speaker)
	
	response, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(audioQuery))
	if err != nil {
		return nil, fmt.Errorf("{comment}", err)
	}
	defer response.Body.Close()
	
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("{comment}", response.Status)
	}
	return io.ReadAll(response.Body)
}

func saveAudioFile(data []byte) error {
	soundFilePath := filepath.Join(STORAGE_PATH, VOICE_FILENAME)
	if err := os.WriteFile(soundFilePath, data, 0644); err != nil {
		return fmt.Errorf("{comment}", err)
	}
	return nil
}

func playAudioFile() error {
	soundFilePath := filepath.Join(STORAGE_PATH, VOICE_FILENAME)
	cmd := exec.Command("aplay", "-D", "plughw:1,0", soundFilePath)
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("{comment}", err)
	}
	return nil
}

func cleanupAudioFile() error {
	soundFilePath := filepath.Join(STORAGE_PATH, VOICE_FILENAME)
	return os.Remove(soundFilePath)
}
