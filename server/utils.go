package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	VOICEVOX_API_URL = "http://voicevox-engine:50021"
	STORAGE_PATH = "/app/storage"
	VOICE_FILENAME = "voice.wav"
	CONFIG_FILENAME = "config.json"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "*")

		// OPTION method handling (CORS preflight request)
		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusOK)
			return
		}
		next(writer, request)
	}
}

func sendErrorResponse(writer http.ResponseWriter, statusCode int, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	response := APIResponse{
		Status: "error",
		Message: message,
	}
	json.NewEncoder(writer).Encode(response)
}

func sendSuccessResponse(writer http.ResponseWriter, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	response := APIResponse{
		Status: "success",
		Message: message,
	}
	json.NewEncoder(writer).Encode(response)
}

func createFile(data []byte, filePath string) error {
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func readJsonFile(filePath string, v interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}
	return nil
}
