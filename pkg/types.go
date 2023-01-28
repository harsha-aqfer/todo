package pkg

import (
	"fmt"
	"github.com/harsha-aqfer/todo/internal/util"
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
	if tr.Task == "" {
		return fmt.Errorf("inadequate input parameters. Required field: task")
	}

	category := tr.Category
	tr.Category = strings.ToLower(tr.Category)

	categories := []string{"work", "home"}

	if !util.Contains(categories, tr.Category) {
		return fmt.Errorf("unknown category value: %s", category)
	}

	pr := tr.Priority
	tr.Priority = strings.ToLower(tr.Priority)

	priorities := []string{"low", "medium", "high"}

	if !util.Contains(priorities, tr.Priority) {
		return fmt.Errorf("unknown priority value: %s", pr)
	}
	return nil
}
