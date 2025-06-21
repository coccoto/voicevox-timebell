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
	CONFIG_FILENAME = "config.json"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "*")

		// プリフライトリクエスト (OPTIONS)
		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusOK)
			return
		}
		next(writer, request)
	}
}

func main() {
	http.HandleFunc("/api/announce", corsMiddleware(announceHandler))
	http.HandleFunc("/api/config", corsMiddleware(configHandler))
	// HTTP Server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start http server. Error:", err)
		return
	}
}

func announceHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("announceHandler called")

	// 音声メッセージを取得する
	speechMessage := getSpeechMessage()
	// request audio_query API
	audioQuery, err := requestAudioQuery(speechMessage, 1)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to requestAudioQuery. Error: %v", err), http.StatusInternalServerError)
		return
	}
	// request synthesis API
	audioData, err := requestSynthesis(audioQuery, 1)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to requestSynthesis. Error: %v", err), http.StatusInternalServerError)
		return
	}
	// 音声ファイルを STORAGE_PATH に保存する
	err = saveAudioFile(audioData)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to saveAudioFile. Error: %v", err), http.StatusInternalServerError)
		return
	}
	// 音声ファイルを再生する
	err = playAudioFile()
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to playAudioFile. Error: %v", err), http.StatusInternalServerError)
		return
	}
	// 音声ファイルを削除する
	err = os.Remove(STORAGE_PATH + "/" + VOICE_FILENAME)
	if err != nil {
		fmt.Printf("Failed to remove voice.wav. Error: %v", err)
	}
	fmt.Println("announceHandler completed")
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

func configHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("configHandler called")

	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to read request body. Error: %v", err), http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	err = os.MkdirAll(STORAGE_PATH, 0777)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to create storage directory. Error: %v", err), http.StatusInternalServerError)
		return
	}

	configFilePath := STORAGE_PATH + "/" + CONFIG_FILENAME
	err = os.WriteFile(configFilePath, body, 0777)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to save config file. Error: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println("configHandler completed")
}
