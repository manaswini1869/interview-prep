package main

type WorkerScript struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Content   string `json:"content"` // The actual JS code
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type CreateRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}
