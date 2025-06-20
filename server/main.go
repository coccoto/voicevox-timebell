package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func main() {
	http.HandleFunc("/start", startHandler)
	// HTTP Server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	speechMessage := "５時になりました"
	// Step 1: audio_query API
	audioQuery, err := requestAudioQuery(speechMessage, 1)
	if err != nil {
		return
	}
	// Step 2: synthesis API
	audioData, err := requeStsynthesis(audioQuery, 1)
	if err != nil {
		return
	}
	// Step 3: 音声データを /app/storage に保存する
	err = saveAudioFile(audioData, "/app/storage")
	if err != nil {
		return
	}
	// Step 4: 音声ファイルを再生する

	// Step 5: 音声ファイルを削除する
	err = os.Remove("/app/storage/voice.wav")
	if err != nil {
		return
	}
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
		return nil, nil
	}
	return io.ReadAll(response.Body)
}

func requeStsynthesis(audioQuery []byte, speaker int) ([]byte, error) {
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

func saveAudioFile(data []byte, dirPath string) error {
	err := os.MkdirAll(dirPath, 0777)
	if err != nil {
		return err
	}
	filePath := dirPath + "/voice.wav"
	err = os.WriteFile(filePath, data, 0777)
	if err != nil {
		return err
	}
	return nil
}