package transport

type InfoResponse struct {
	OK     bool   `json:"ok"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
	Logs   string `json:"logs,omitempty"`
}

type PathItem struct {
	Path  string `json:"path"`
	IsDir bool   `json:"isDir,omitempty"`
	Size  int64  `json:"size,omitempty"`
}

type PathResponse struct {
	OK    bool       `json:"ok"`
	Root  string     `json:"root,omitempty"`
	Roots []string   `json:"roots,omitempty"`
	Items []PathItem `json:"items,omitempty"`
	Error string     `json:"error,omitempty"`
}
