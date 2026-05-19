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
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if isFirstLine {
			isFirstLine = false
			if r.isHeaderLine(line) {
				continue
			}
		}

		track, valid := r.parseTrackLine(line)
		if valid {
			tracks = append(tracks, track)
		}
	}

	return tracks, scanner.Err()
}

func (r *LocalFileRepository) isHeaderLine(line string) bool {
	lowerLine := strings.ToLower(line)
	return strings.Contains(lowerLine, "title") && strings.Contains(lowerLine, "artist")
}

func (r *LocalFileRepository) parseTrackLine(line string) (domain.Track, bool) {
	parts := strings.Split(line, "|||")

	if len(parts) >= 3 {
		return domain.Track{
			Path:   strings.TrimSpace(parts[0]),
			Title:  strings.TrimSpace(parts[1]),
			Artist: strings.TrimSpace(parts[2]),
		}, true
	}

	if len(parts) == 2 {
		return domain.Track{
			Path:   "",
			Title:  strings.TrimSpace(parts[0]),
			Artist: strings.TrimSpace(parts[1]),
		}, true
	}

	return domain.Track{}, false
}
