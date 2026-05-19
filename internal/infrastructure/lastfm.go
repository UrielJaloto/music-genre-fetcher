package infrastructure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/UrielJaloto/music-genre-fetcher/internal/domain"
)

type LastFMResponse struct {
	TopTags struct {
		Tag []struct {
			Name string `json:"name"`
		} `json:"tag"`
	} `json:"toptags"`
	Error int `json:"error"`
}

type LastFMProvider struct {
	client *http.Client
}

func NewLastFMProvider() *LastFMProvider {
	return &LastFMProvider{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *LastFMProvider) FetchGenre(path, artist, title, apiKey string) domain.Result {
	genre := p.fetchTags(artist, title, apiKey, "track.gettoptags")

	if genre == "Not Found" {
		genre = p.fetchTags(artist, "", apiKey, "artist.gettoptags")
	}

	return domain.Result{Path: path, Artist: artist, Title: title, Genre: genre}
}

func (p *LastFMProvider) fetchTags(artist, title, apiKey, method string) string {
	baseURL := "http://ws.audioscrobbler.com/2.0/"

	params := url.Values{}
	params.Add("method", method)
	params.Add("artist", artist)

	if title != "" {
		params.Add("track", title)
	}

	params.Add("api_key", apiKey)
	params.Add("format", "json")
	params.Add("autocorrect", "1")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	isYear := regexp.MustCompile(`^[0-9]{4}$`)

	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			return "Not Found"
		}

		req.Header.Set("User-Agent", "GenreFetcherGo/1.0.0")

		resp, err := p.client.Do(req)
		if err != nil {
			return "Not Found"
		}

		if resp.StatusCode == 429 || resp.StatusCode == 503 {
			resp.Body.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		var lfmResp LastFMResponse
		err = json.NewDecoder(resp.Body).Decode(&lfmResp)
		resp.Body.Close()

		if err != nil {
			return "Not Found"
		}

		if lfmResp.Error > 0 {
			return "Not Found"
		}

		if len(lfmResp.TopTags.Tag) == 0 {
			return "Unknown"
		}

		for _, tag := range lfmResp.TopTags.Tag {
			if !isYear.MatchString(tag.Name) {
				return tag.Name
			}
		}

		return "Unknown"
	}

	return "Not Found"
}
