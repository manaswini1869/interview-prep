package model

type Token struct {
	ID      string   `json:"id"`
	Revoked bool     `json:"revoked"`
	Scopes  []string `json:"scopes"`
}
