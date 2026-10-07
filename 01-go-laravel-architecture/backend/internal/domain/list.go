package domain

import "time"

// List representa a tabela 'lists' no banco de dados.
type List struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`

	// Relacionamentos carregados sob demanda (Equivalente ao eager loading do Eloquent)
	Tasks []Task `json:"tasks,omitempty"`
	Tags  []Tag  `json:"tags,omitempty"`
}
