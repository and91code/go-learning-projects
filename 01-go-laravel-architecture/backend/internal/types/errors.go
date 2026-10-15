package types

import "errors"

var (
	ErrUserNotFound           = errors.New("usuário não encontrado")
	ErrEmailExists            = errors.New("e-mail já está cadastrado")
	ErrInvalidCredentials     = errors.New("credenciais inválidas")
	ErrInvalidCurrentPassword = errors.New("senha atual incorreta")
	ErrListNotFound           = errors.New("lista não encontrada")
	ErrTagNotFound            = errors.New("tag não encontrada")
	ErrTaskNotFound           = errors.New("tarefa não encontrada")
	ErrInvalidTag             = errors.New("invalid tag")
	ErrInvalidTask            = errors.New("invalid task")
)
