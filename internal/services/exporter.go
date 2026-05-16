package services

import "github.com/UrielJaloto/music-genre-fetcher/internal/domain"

type DataExporter interface {
	ExportJSON(path string, data []domain.Result) error
	ExportCSV(path string, data []domain.Result) error
	ExportTXT(path string, data []domain.Result) error
}
