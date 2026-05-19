package domain

import "time"

type Configuration struct {
	LastFMAPIKey string `json:"lastfm_api_key"`
	InputFile    string
	ExportPaths  map[string]string
	Concurrency  int
	Pause        time.Duration
}
