package main

import (
	"fmt"
	"io"
	"net/http"
)

func voicevoxSpeakersHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Processing voicevox speakers request")

	fmt.println(VOICEVOX_API_URL)
	response, err := http.Get(fmt.Sprintf("%s/speakers", VOICEVOX_API_URL))
	if err != nil {
		sendErrorResponse(writer, http.StatusInternalServerError, "Failed to fetch speakers from VoiceVox API")
		return
	}
	defer response.Body.Close()

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(response.StatusCode)
	if _, err := io.Copy(writer, response.Body); err != nil {
		fmt.Println("Failed to copy response body:", err)
	}
}
