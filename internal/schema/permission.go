package schema

type RouteList struct {
	ID     string   `json:"id"`
	PID    string   `json:"pid"`
	Name   string   `json:"name"`
	Url    []string `json:"url"`
	IsPage int64    `json:"is_page"`
}

type Node struct {
	ID       string   `json:"id"`
	PID      string   `json:"pid"`
	Name     string   `json:"name"`
	URL      []string `json:"url"`
	IsPage   int64    `json:"is_page"`
	Checked  int64    `json:"checked"`
	Disabled bool     `json:"disabled"`
	Children []Node   `json:"children"`
}
