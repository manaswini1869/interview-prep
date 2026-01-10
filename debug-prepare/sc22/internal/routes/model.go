package routes

type Route struct {
	ID       string `json:"id"`
	Pattern  string `json:"pattern"`
	Priority int    `json:"priority"`
	Enabled  bool   `json:"enabled"`
}
