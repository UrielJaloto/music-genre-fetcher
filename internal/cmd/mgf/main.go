package main

import (
	"github.com/UrielJaloto/music-genre-fetcher/internal/infrastructure"
	"github.com/UrielJaloto/music-genre-fetcher/internal/presentation"
	"github.com/UrielJaloto/music-genre-fetcher/internal/services"
)

const ConfigFile = "env/config.json"

func main() {
	ui := presentation.NewCLIUserInterface()
	configLoader := infrastructure.NewLocalConfigProvider()
	repo := infrastructure.NewLocalFileRepository()
	provider := infrastructure.NewLastFMProvider()

	exporters := map[string]services.Exporter{
		"json": presentation.NewJSONExporter(),
		"csv":  presentation.NewCSVExporter(),
		"txt":  presentation.NewTXTExporter(),
	}

	deps := services.ApplicationDependencies{
		ConfigLoader: configLoader,
		Repo:         repo,
		Provider:     provider,
		Exporters:    exporters,
		UI:           ui,
	}

	app := services.NewApplication(deps)

	app.Execute(ConfigFile)
}
