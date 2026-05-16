package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/UrielJaloto/music-genre-fetcher/internal/domain"
)

type Application struct {
	repo     MusicRepository
	provider MusicDataProvider
	exporter DataExporter
	ui       UserInterface
}

func NewApplication(repo MusicRepository, provider MusicDataProvider, exporter DataExporter, ui UserInterface) *Application {
	return &Application{
		repo:     repo,
		provider: provider,
		exporter: exporter,
		ui:       ui,
	}
}

func (app *Application) Execute(inputFile, outputJSON, outputCSV, outputTXT string, concurrency int, pause time.Duration) {
	app.ui.ShowMessage("--- Starting Music Processing in Go ---")

	tracks, err := app.repo.LoadMusic(inputFile)
	if err != nil {
		app.ui.ShowFatalError("Error loading music", err)
		return
	}

	totalTracks := len(tracks)
	app.ui.ShowMessage(fmt.Sprintf("Total tracks loaded: %d", totalTracks))

	tracksChan := make(chan domain.Track, totalTracks)
	resultsChan := make(chan domain.Result, totalTracks)

	for _, t := range tracks {
		tracksChan <- t
	}
	close(tracksChan)

	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range tracksChan {
				res := app.provider.FetchGenre(t.Path, t.Artist, t.Title)
				resultsChan <- res
				time.Sleep(pause)
			}
		}()
	}

	wg.Wait()
	close(resultsChan)

	var finalResults []domain.Result
	for res := range resultsChan {
		finalResults = append(finalResults, res)
	}

	app.ui.ShowMessage(fmt.Sprintf("Total tracks processed: %d", len(finalResults)))

	if err := app.exporter.ExportJSON(outputJSON, finalResults); err != nil {
		app.ui.ShowError("Error saving JSON file", err)
	} else {
		app.ui.ShowMessage(fmt.Sprintf("JSON file saved: %s", outputJSON))
	}

	if err := app.exporter.ExportCSV(outputCSV, finalResults); err != nil {
		app.ui.ShowError("Error saving CSV file", err)
	} else {
		app.ui.ShowMessage(fmt.Sprintf("CSV file saved: %s", outputCSV))
	}

	if err := app.exporter.ExportTXT(outputTXT, finalResults); err != nil {
		app.ui.ShowError("Error saving TXT file", err)
	} else {
		app.ui.ShowMessage(fmt.Sprintf("TXT file saved: %s", outputTXT))
	}

	app.ui.ShowMessage("--- Processing completed ---")
}
