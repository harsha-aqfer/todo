package pkg

import (
	"fmt"
	"strings"
	"time"
)

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

func (tr *TodoRequest) Validate() error {
	// task should be non-empty
	if tr.Task == "" {
		return fmt.Errorf("inadequate input parameters. Required field: task")
	}

	cat := tr.Category
	tr.Category = strings.ToLower(tr.Category)

	if tr.Category != "work" && tr.Category != "home" {
		return fmt.Errorf("unknown category value: %s", cat)
	}

	pr := tr.Priority
	tr.Priority = strings.ToLower(tr.Priority)

	if tr.Priority != "low" && tr.Priority != "medium" && tr.Priority != "high" {
		return fmt.Errorf("unknown priority value: %s", pr)
	}
	return nil
}
