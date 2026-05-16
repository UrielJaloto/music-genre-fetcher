package infrastructure

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/UrielJaloto/music-genre-fetcher/internal/domain"
)

type LocalFileRepository struct{}

func NewLocalFileRepository() *LocalFileRepository {
	return &LocalFileRepository{}
}

func (r *LocalFileRepository) LoadMusic(path string) ([]domain.Track, error) {
	if strings.HasSuffix(path, ".json") {
		return r.loadJSON(path)
	}
	return r.loadTXT(path)
}

func (r *LocalFileRepository) loadJSON(path string) ([]domain.Track, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tracks []domain.Track
	err = json.Unmarshal(data, &tracks)

	return tracks, err
}

func (r *LocalFileRepository) loadTXT(path string) ([]domain.Track, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tracks []domain.Track
	scanner := bufio.NewScanner(file)
	isFirstLine := true

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		if isFirstLine {
			isFirstLine = false
			if strings.Contains(strings.ToLower(line), "title") && strings.Contains(strings.ToLower(line), "artist") {
				continue
			}
		}

		parts := strings.Split(line, "|||")
		if len(parts) >= 3 {
			tracks = append(tracks, domain.Track{
				Path:   strings.TrimSpace(parts[0]),
				Title:  strings.TrimSpace(parts[1]),
				Artist: strings.TrimSpace(parts[2]),
			})
		} else if len(parts) == 2 {
			tracks = append(tracks, domain.Track{
				Path:   "",
				Title:  strings.TrimSpace(parts[0]),
				Artist: strings.TrimSpace(parts[1]),
			})
		}
	}

	return tracks, scanner.Err()
}
