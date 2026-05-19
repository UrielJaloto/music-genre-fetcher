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

func (app *Application) Execute(configPath, inputFilePath string) error {
	config, err := app.deps.ConfigLoader.Load(configPath, inputFilePath)
	if err != nil {
		return fmt.Errorf("error loading configuration file: %w", err)
	}

	if config.LastFMAPIKey == "" {
		return errors.New("the API key (lastfm_api_key) is empty in the configuration file")
	}

	app.deps.UI.ShowMessage("--- Starting Music Processing in Go ---")

	tracks, err := app.deps.Repo.LoadMusic(config.InputFile)
	if err != nil {
		return fmt.Errorf("error loading music: %w", err)
	}

	app.deps.UI.ShowMessage(fmt.Sprintf("Total tracks loaded: %d", len(tracks)))

	finalResults := app.processTracksConcurrently(tracks, config)

	app.deps.UI.ShowMessage(fmt.Sprintf("Total tracks processed: %d", len(finalResults)))

	app.exportResults(finalResults, config.ExportPaths)

	app.deps.UI.ShowMessage("--- Processing completed ---")
	return nil
}

func (app *Application) processTracksConcurrently(tracks []domain.Track, config domain.Configuration) []domain.Result {
	tracksChan := make(chan domain.Track)
	resultsChan := make(chan domain.Result, len(tracks))

	go func() {
		for _, t := range tracks {
			tracksChan <- t
		}
		close(tracksChan)
	}()

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

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var finalResults []domain.Result
	for res := range resultsChan {
		finalResults = append(finalResults, res)
	}

	return finalResults
}

func (app *Application) exportResults(results []domain.Result, paths map[string]string) {
	for format, path := range paths {
		exporter, exists := app.deps.Exporters[format]
		if !exists {
			continue
		}

		err := exporter.Export(path, results)
		if err != nil {
			app.deps.UI.ShowError(fmt.Sprintf("Error saving %s file", format), err)
			continue
		}

		app.deps.UI.ShowMessage(fmt.Sprintf("%s file saved: %s", format, path))
	}
}
