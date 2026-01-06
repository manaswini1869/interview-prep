package model

type Rollout struct {
	ID        string            `json:"id"`
	Regions   map[string]string `json:"regions"` // region -> status
	Completed bool              `json:"completed"`
}

type StatusUpdate struct {
	Region string `json:"region"`
	Status string `json:"status"`
}
