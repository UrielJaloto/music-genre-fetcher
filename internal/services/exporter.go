package services

import "github.com/UrielJaloto/music-genre-fetcher/internal/domain"

type Exporter interface {
	Export(path string, data []domain.Result) error
}
