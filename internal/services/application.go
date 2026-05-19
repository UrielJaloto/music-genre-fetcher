package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/UrielJaloto/music-genre-fetcher/internal/domain"
)

type ApplicationDependencies struct {
	ConfigLoader ConfigLoader
	Repo         Repository
	Provider     Provider
	Exporters    map[string]Exporter
	UI           UI
}

type Application struct {
	deps ApplicationDependencies
}

func NewApplication(deps ApplicationDependencies) *Application {
	return &Application{
		deps: deps,
	}
}

func (app *Application) Execute(configPath string) {
	config, err := app.deps.ConfigLoader.Load(configPath)
	if err != nil {
		app.deps.UI.ShowFatalError("Error loading configuration file", err)
		return
	}

	if config.LastFMAPIKey == "" {
		app.deps.UI.ShowFatalError("The API key (lastfm_api_key) is empty in the configuration file", errors.New("empty api key"))
		return
	}

	app.deps.UI.ShowMessage("--- Starting Music Processing in Go ---")

	tracks, err := app.deps.Repo.LoadMusic(config.InputFile)
	if err != nil {
		app.deps.UI.ShowFatalError("Error loading music", err)
		return
	}

	totalTracks := len(tracks)
	app.deps.UI.ShowMessage(fmt.Sprintf("Total tracks loaded: %d", totalTracks))

	tracksChan := make(chan domain.Track, totalTracks)
	resultsChan := make(chan domain.Result, totalTracks)

	for _, t := range tracks {
		tracksChan <- t
	}
	close(tracksChan)

	var wg sync.WaitGroup

	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range tracksChan {
				res := app.deps.Provider.FetchGenre(t.Path, t.Artist, t.Title, config.LastFMAPIKey)
				resultsChan <- res
				time.Sleep(config.Pause)
			}
		}()
	}

	wg.Wait()
	close(resultsChan)

	var finalResults []domain.Result
	for res := range resultsChan {
		finalResults = append(finalResults, res)
	}

	app.deps.UI.ShowMessage(fmt.Sprintf("Total tracks processed: %d", len(finalResults)))

	for format, path := range config.ExportPaths {
		if exporter, exists := app.deps.Exporters[format]; exists {
			if err := exporter.Export(path, finalResults); err != nil {
				app.deps.UI.ShowError(fmt.Sprintf("Error saving %s file", format), err)
			} else {
				app.deps.UI.ShowMessage(fmt.Sprintf("%s file saved: %s", format, path))
			}
		}
	}

	app.deps.UI.ShowMessage("--- Processing completed ---")
}
