package domain

import "time"

// User representa a tabela 'users' no banco de dados.
type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // O hífen é para impedir que o hash da senha seja vazado em serializações JSON
	CreatedAt    time.Time `json:"created_at"`
}
