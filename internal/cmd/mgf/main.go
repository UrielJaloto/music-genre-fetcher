package main

import (
	"flag"
	"time"

	"github.com/UrielJaloto/music-genre-fetcher/internal/infrastructure"
	"github.com/UrielJaloto/music-genre-fetcher/internal/presentation"
	"github.com/UrielJaloto/music-genre-fetcher/internal/services"
)

const (
	MaxConcurrency = 5
	PauseMs        = 200 * time.Millisecond
	OutputJSON     = "env/output/genre_results.json"
	OutputCSV      = "env/output/genre_results.csv"
	OutputTXT      = "env/output/ai_results.txt"
	ConfigFile     = "env/config.json"
)

func main() {
	inputFile := flag.String("input", "env/input/mp3tag_data.txt", "")
	flag.Parse()

	ui := presentation.NewCLIUserInterface()

	configProvider := infrastructure.NewLocalConfigProvider()
	config, err := configProvider.LoadConfiguration(ConfigFile)
	if err != nil {
		ui.ShowFatalError("Error loading configuration file", err)
	}

	if config.LastFMAPIKey == "" {
		ui.ShowFatalError("The API key (lastfm_api_key) is empty in the configuration file", nil)
	}

	repo := infrastructure.NewLocalFileRepository()
	provider := infrastructure.NewLastFMProvider(config.LastFMAPIKey)
	exporter := presentation.NewExporter()

	app := services.NewApplication(repo, provider, exporter, ui)

	app.Execute(*inputFile, OutputJSON, OutputCSV, OutputTXT, MaxConcurrency, PauseMs)
}
