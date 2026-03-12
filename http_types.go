package main

type infoResponse struct {
	OK     bool   `json:"ok"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

type pathResponse struct {
	OK    bool     `json:"ok"`
	Root  string   `json:"root,omitempty"`
	Roots []string `json:"roots,omitempty"`
	Items []string `json:"items,omitempty"`
	Error string   `json:"error,omitempty"`
}
