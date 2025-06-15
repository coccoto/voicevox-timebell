package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

// logRequest は受信したリクエストをログに記録します
func logRequest(handlerName string) {
	log.Printf("%s: Request received", handlerName)
}

// decodeRequestBody は JSON リクエストボディをデコードします
func decodeRequestBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

// saveResponseToFile は HTTP レスポンスボディをファイルに保存します
func saveResponseToFile(body io.Reader, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, body)
	return err
}

// ConfigMutex は設定へのスレッドセーフなアクセスを提供します
var ConfigMutex sync.RWMutex
