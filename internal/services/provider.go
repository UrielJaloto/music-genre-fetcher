package services

import "github.com/UrielJaloto/music-genre-fetcher/internal/domain"

type Provider interface {
	FetchGenre(path, artist, title, apiKey string) domain.Result
}
