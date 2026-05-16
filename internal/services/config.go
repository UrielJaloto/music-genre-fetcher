package services

import "github.com/UrielJaloto/music-genre-fetcher/internal/domain"

type ConfigProvider interface {
	LoadConfiguration(path string) (domain.Configuration, error)
}
