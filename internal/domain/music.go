package domain

type Track struct {
	Path   string `json:"path"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

type Result struct {
	Path   string `json:"path"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Genre  string `json:"genre"`
}
