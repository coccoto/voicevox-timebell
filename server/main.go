package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"
)

const (
	STORAGE_PATH = "/app/storage"
	VOICE_FILENAME = "voice.wav"
)

func main() {
	http.HandleFunc("/start", startHandler)
	// HTTP Server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start server. Error:", err)
		return
	}
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("startHandler called")

	// 音声メッセージを取得する
	speechMessage := getSpeechMessage()
	// request audio_query API
	audioQuery, err := requestAudioQuery(speechMessage, 1)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to requestAudioQuery. Error: %v", err), http.StatusInternalServerError)
		return
	}
	// request synthesis API
	audioData, err := requestSynthesis(audioQuery, 1)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to requestSynthesis. Error: %v", err), http.StatusInternalServerError)
		return
	}
	// 音声ファイルを STORAGE_PATH に保存する
	err = saveAudioFile(audioData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to saveAudioFile. Error: %v", err), http.StatusInternalServerError)
		return
	}
	// 音声ファイルを再生する
	err = playAudioFile()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to playAudioFile. Error: %v", err), http.StatusInternalServerError)
		return
	}
	// 音声ファイルを削除する
	err = os.Remove(STORAGE_PATH + "/" + VOICE_FILENAME)
	if err != nil {
		fmt.Printf("Failed to remove voice.wav. Error: %v", err)
	}
	fmt.Println("startHandler completed")
}

func getSpeechMessage() string {
	// 現在時刻を取得する
	hour := time.Now().Hour()
	return fmt.Sprintf("%d時になりました", hour)
}

func requestAudioQuery(speechMessage string, speaker int) ([]byte, error) {
	requestUrl := fmt.Sprintf("http://voicevox-engine:50021/audio_query?text=%s&speaker=%d", url.QueryEscape(speechMessage), speaker)
	// API Request
	response, err := http.Post(requestUrl, "", nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("audio_query API returned status: %s", response.Status)
	}
	return io.ReadAll(response.Body)
}

func requestSynthesis(audioQuery []byte, speaker int) ([]byte, error) {
	requestUrl := fmt.Sprintf("http://voicevox-engine:50021/synthesis?speaker=%d", speaker)
	// API Request
	response, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(audioQuery))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	
	if response.StatusCode != http.StatusOK {
		return nil, nil
	}
	return io.ReadAll(response.Body)
}

func saveAudioFile(data []byte) error {
	err := os.MkdirAll(STORAGE_PATH, 0777)
	if err != nil {
		return err
	}
	// 音声ファイルの保存先を指定する
	soundFilePath := STORAGE_PATH + "/" + VOICE_FILENAME
	err = os.WriteFile(soundFilePath, data, 0777)
	if err != nil {
		return err
	}
	return nil
}

func playAudioFile() error {
	// USB キャプチャデバイス (card 1) から音声を出力する
	cmd := exec.Command("aplay", "-D", "plughw:1,0", STORAGE_PATH + "/" + VOICE_FILENAME)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
