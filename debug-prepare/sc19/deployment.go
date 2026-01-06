package main

import "time"

type Deployment struct {
	ID        string    `json:"id"`
	Script    string    `json:"script"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
