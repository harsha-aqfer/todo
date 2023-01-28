package pkg

import "time"

type TodoRequest struct {
	Task     string `json:"task"`
	Done     bool   `json:"done,omitempty"`
	Category string `json:"category,omitempty"`
	Priority string `json:"priority,omitempty"`
}

type TodoResponse struct {
	Id          int64      `json:"id"`
	Task        string     `json:"task"`
	Category    string     `json:"category"`
	Priority    string     `json:"priority"`
	CreatedAt   *time.Time `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
