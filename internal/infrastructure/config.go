package infrastructure

import (
	"encoding/json"
	"os"
	"time"

	"github.com/UrielJaloto/music-genre-fetcher/internal/domain"
)

type LocalConfigProvider struct{}

func NewLocalConfigProvider() *LocalConfigProvider {
	return &LocalConfigProvider{}
}

func (p *LocalConfigProvider) Load(configPath, inputFilePath string) (domain.Configuration, error) {
	var config domain.Configuration

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	config.InputFile = inputFilePath
	config.ExportPaths = map[string]string{
		"json": "env/output/genre_results.json",
		"csv":  "env/output/genre_results.csv",
		"txt":  "env/output/ai_results.txt",
	}
	config.Concurrency = 5
	config.Pause = 200 * time.Millisecond

	return config, nil
}
