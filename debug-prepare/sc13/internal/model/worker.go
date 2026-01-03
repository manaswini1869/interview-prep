package model

type Worker struct {
	ID      string `json:"id"`
	Script  string `json:"script"`
	Active  bool   `json:"active"`
	Version int    `json:"version"`
}
