package model

type Route struct {
	ID       string
	Pattern  string
	Priority int
	Enabled  bool
}
