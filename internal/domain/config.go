package domain

import "time"

type Configuration struct {
	LastFMAPIKey string            `json:"lastfm_api_key"`
	InputFile    string            `json:"input_file,omitempty"`
	ExportPaths  map[string]string `json:"export_paths,omitempty"`
	Concurrency  int               `json:"concurrency,omitempty"`
	Pause        time.Duration     `json:"pause,omitempty"`
}
