package services

import "github.com/UrielJaloto/music-genre-fetcher/internal/domain"

type MusicRepository interface {
	LoadMusic(path string) ([]domain.Track, error)
}
