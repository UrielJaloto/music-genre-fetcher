package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/UrielJaloto/music-genre-fetcher/internal/infrastructure"
	"github.com/UrielJaloto/music-genre-fetcher/internal/presentation"
	"github.com/UrielJaloto/music-genre-fetcher/internal/services"
)

const ConfigFile = "env/config.json"

func main() {
	inputFile := flag.String("input", "env/input/mp3tag_data.txt", "")
	flag.Parse()

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

	if err := app.Execute(ConfigFile, *inputFile); err != nil {
		fmt.Printf("FATAL ERROR: %v\n", err)
		os.Exit(1)
	}
}
