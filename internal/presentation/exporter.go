package presentation

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/UrielJaloto/music-genre-fetcher/internal/domain"
)

func ensureDirectory(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, os.ModePerm)
}

type JSONExporter struct{}

func NewJSONExporter() *JSONExporter {
	return &JSONExporter{}
}

func (e *JSONExporter) Export(path string, data []domain.Result) error {
	if err := ensureDirectory(path); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, jsonData, 0644)
}

type CSVExporter struct{}

func NewCSVExporter() *CSVExporter {
	return &CSVExporter{}
}

func (e *CSVExporter) Export(path string, data []domain.Result) error {
	if err := ensureDirectory(path); err != nil {
		return err
	}

	csvFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)

	err = writer.Write([]string{"path", "artist", "title", "genre"})
	if err != nil {
		return err
	}

	for _, r := range data {
		err := writer.Write([]string{r.Path, r.Artist, r.Title, r.Genre})
		if err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}

type TXTExporter struct{}

func NewTXTExporter() *TXTExporter {
	return &TXTExporter{}
}

func (e *TXTExporter) Export(path string, data []domain.Result) error {
	if err := ensureDirectory(path); err != nil {
		return err
	}

	txtFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer txtFile.Close()

	for _, r := range data {
		line := fmt.Sprintf("%s|||%s|||%s|||%s\n", r.Path, r.Title, r.Artist, r.Genre)
		if _, err := txtFile.WriteString(line); err != nil {
			return err
		}
	}
	return nil
}
