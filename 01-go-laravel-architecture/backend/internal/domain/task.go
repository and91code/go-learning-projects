package domain

import "time"

// Task representa a tabela 'tasks' no banco de dados.
type Task struct {
	ID          int64        `json:"id"`
	ListID      int64        `json:"list_id"`
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	DueDate     *time.Time   `json:"due_date,omitempty"` // Ponteiro permite tratar NULL
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`

	// Relacionamento N:N com Tags
	Tags []Tag `json:"tags,omitempty"`
}
