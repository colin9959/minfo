// Package transport 定义 HTTP 传输层使用的响应结构。

package transport

// InfoResponse 表示信息类接口共用的 JSON 响应。
type InfoResponse struct {
	OK     bool   `json:"ok"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
	Logs   string `json:"logs,omitempty"`
}

// PathItem 表示路径联想接口返回的一条候选路径。
type PathItem struct {
	Path  string `json:"path"`
	IsDir bool   `json:"isDir,omitempty"`
	Size  int64  `json:"size,omitempty"`
}

// PathResponse 表示路径联想接口的 JSON 响应。
type PathResponse struct {
	OK    bool       `json:"ok"`
	Root  string     `json:"root,omitempty"`
	Roots []string   `json:"roots,omitempty"`
	Items []PathItem `json:"items,omitempty"`
	Error string     `json:"error,omitempty"`
}
