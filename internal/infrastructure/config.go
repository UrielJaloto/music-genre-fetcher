package infrastructure

import (
	"encoding/json"
	"os"

	"github.com/UrielJaloto/music-genre-fetcher/internal/domain"
)

type LocalConfigProvider struct{}

func NewLocalConfigProvider() *LocalConfigProvider {
	return &LocalConfigProvider{}
}

func (p *LocalConfigProvider) LoadConfiguration(path string) (domain.Configuration, error) {
	var config domain.Configuration

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}
