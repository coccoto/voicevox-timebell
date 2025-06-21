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
	fmt.Println("{comment}")

	body, err := io.ReadAll(request.Body)
	if err != nil {
		sendErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("{comment}", err))
		return
	}
	defer request.Body.Close()

	// Validate JSON format
	var config Config
	if err := json.Unmarshal(body, &config); err != nil {
		sendErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("{comment}", err))
		return
	}

	if err := saveConfigFile(body); err != nil {
		sendErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("{comment}", err))
		return
	}
	sendSuccessResponse(writer, "{comment}")
}

func saveConfigFile(data []byte) error {
	configFilePath := filepath.Join(STORAGE_PATH, CONFIG_FILENAME)
	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		return fmt.Errorf("{comment}", err)
	}
	return nil
}
