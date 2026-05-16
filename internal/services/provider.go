package services

import "github.com/UrielJaloto/music-genre-fetcher/internal/domain"

type MusicDataProvider interface {
	FetchGenre(path, artist, title string) domain.Result
}
