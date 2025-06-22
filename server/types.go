package main

type Config struct {
	HourList []string `json:"hourList"`
	SpeakerUuid string `json:"speakerUuid"`
	StyleID string `json:"styleId"`
}

type APIResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
}
