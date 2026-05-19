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

func (p *LocalConfigProvider) Load(configPath string) (domain.Configuration, error) {
	var config domain.Configuration

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	if len(config.ExportPaths) == 0 {
		config.ExportPaths = map[string]string{
			"json": "env/output/genre_results.json",
			"csv":  "env/output/genre_results.csv",
			"txt":  "env/output/ai_results.txt",
		}
	}

	if config.Concurrency <= 0 {
		config.Concurrency = 5
	}

	if config.Pause <= 0 {
		config.Pause = 200 * time.Millisecond
	}

	return config, nil
}
