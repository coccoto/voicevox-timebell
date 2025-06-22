package main

import (
	"fmt"
	"io"
	"net/http"
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

	if err := createFile(body, filepath.Join(STORAGE_PATH, CONFIG_FILENAME)); err != nil {
		sendErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("Failed to save config: %v", err))
		return
	}
	sendSuccessResponse(writer, "Finished configHandler")
}
