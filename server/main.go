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

func main() {
	http.HandleFunc("/start", startHandler)
	// HTTP Server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("startHandler")

	speechMessage := getSpeechMessage()
	if speechMessage == "" {
		return
	}
	// Step 1: audio_query API
	audioQuery, err := requestAudioQuery(speechMessage, 1)
	if err != nil {
		http.Error(w, "Failed to requestAudioQuery.", http.StatusInternalServerError)
		return
	}
	// Step 2: synthesis API
	audioData, err := requeStsynthesis(audioQuery, 1)
	if err != nil {
		http.Error(w, "Failed to requeStsynthesis.", http.StatusInternalServerError)
		return
	}
	// Step 3: 音声データを /app/storage に保存する
	err = saveAudioFile(audioData, "/app/storage")
	if err != nil {
		http.Error(w, "Failed to saveAudioFile.", http.StatusInternalServerError)
		return
	}
	// Step 4: 音声ファイルを再生する
	err = playAudioFile("/app/storage/voice.wav")
	if err != nil {
		fmt.Println("Failed to playAudioFile. Error:", err)
	}
	// Step 5: 音声ファイルを削除する
	err = os.Remove("/app/storage/voice.wav")
	if err != nil {
		http.Error(w, "Failed to delete voice.wav", http.StatusInternalServerError)
		return
	}
}

func getSpeechMessage() string {
	// 現在時刻を取得
	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()

	if minute == 0 {
		return fmt.Sprintf("%d時になりました", hour)
	}
	return "TEST"
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

func playAudioFile(filepath string) error {
	cmd := exec.Command("aplay", filepath)

	// 標準出力と標準エラー出力をバッファに保存
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("aplay failed: %v\nstdout: %s\nstderr: %s", err, out.String(), stderr.String())
	}

	return nil
}