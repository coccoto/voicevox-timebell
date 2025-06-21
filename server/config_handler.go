package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func configHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Processing config request")

	body, err := io.ReadAll(request.Body)
	if err != nil {
		sendErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Failed to read request body: %v", err))
		return
	}
	defer request.Body.Close()

	// Validate JSON format
	var config Config
	if err := json.Unmarshal(body, &config); err != nil {
		sendErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Invalid JSON format: %v", err))
		return
	}

	if err := saveConfigFile(body); err != nil {
		sendErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("Failed to save config: %v", err))
		return
	}
	sendSuccessResponse(writer, "Completed configHandler")
}

func saveConfigFile(data []byte) error {
	configFilePath := filepath.Join(STORAGE_PATH, CONFIG_FILENAME)
	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
