package main

import (
	"encoding/json"
	"net/http"
)

const (
	STORAGE_PATH = "/app/storage"
	CONFIG_FILENAME = "config.json"
	VOICE_FILENAME = "voice.wav"
	VOICEVOX_API_URL = "http://voicevox-engine:50021"
	AUDIO_DEVICE = "plughw:1,0"
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
