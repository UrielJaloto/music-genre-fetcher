package services

import "github.com/UrielJaloto/music-genre-fetcher/internal/domain"

type ConfigLoader interface {
	Load(configPath, inputFilePath string) (domain.Configuration, error)
}
