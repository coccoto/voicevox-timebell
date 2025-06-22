package main

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

func configRegisterHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Processing configRegisterHandler")

	body, err := io.ReadAll(request.Body)
	if err != nil {
		sendErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Failed to read request body: %v", err))
		return
	}
	defer request.Body.Close()

	if err := createFile(body, filepath.Join(STORAGE_PATH, CONFIG_FILENAME)); err != nil {
		sendErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("Failed to save config: %v", err))
		return
	}
	sendSuccessResponse(writer, "Finished configRegisterHandler")
}

func configReadHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Processing configReadHandler")

	var config Config
	if err := readJsonFile(filepath.Join(STORAGE_PATH, CONFIG_FILENAME), &config); err != nil {
		sendErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("Failed to read config: %v", err))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(config); err != nil {
		fmt.Println("Failed to encode config response:", err)
	}
}
