package types

import "errors"

var (
	ErrUserNotFound       = errors.New("usuário não encontrado")
	ErrEmailExists        = errors.New("e-mail já está cadastrado")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrListNotFound       = errors.New("lista não encontrada")
	ErrTaskNotFound       = errors.New("tarefa não encontrada")
	ErrInvalidTask        = errors.New("invalid task")
)
