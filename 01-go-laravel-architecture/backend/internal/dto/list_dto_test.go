package dto

import (
	"testing"
	"time"

	"github.com/seu-usuario/taskflow-backend/internal/types"
	"github.com/stretchr/testify/require"
)

// Cenário: CreateListRequest.ToDomain copia título e descrição.
func TestCreateListRequest_ToDomain(t *testing.T) {
	list := CreateListRequest{Title: "Casa", Description: "Tarefas domésticas"}.ToDomain()

	require.Equal(t, types.List{Title: "Casa", Description: "Tarefas domésticas"}, list)
}

// Cenário: UpdateListRequest.ToDomain copia título e descrição.
func TestUpdateListRequest_ToDomain(t *testing.T) {
	list := UpdateListRequest{Title: "Trabalho", Description: "Projetos"}.ToDomain()

	require.Equal(t, types.List{Title: "Trabalho", Description: "Projetos"}, list)
}

// Cenário: lista sem tarefas e tags gera resposta sem esses campos preenchidos.
func TestListResponseFromDomain_WithoutRelations(t *testing.T) {
	createdAt := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)

	response := ListResponseFromDomain(types.List{ID: 1, UserID: 7, Title: "Casa", CreatedAt: createdAt})

	require.Equal(t, ListResponse{ID: 1, UserID: 7, Title: "Casa", CreatedAt: "2026-08-01T10:00:00Z"}, response)
	require.Nil(t, response.Tasks)
	require.Nil(t, response.Tags)
}

// Cenário: lista com tarefas e tags mapeia cada relação para o DTO correspondente.
func TestListResponseFromDomain_WithRelations(t *testing.T) {
	list := types.List{
		ID:    1,
		Tasks: []types.Task{{ID: 10, ListID: 1, Title: "Lavar louça"}},
		Tags:  []types.Tag{{ID: 20, ListID: 1, Name: "urgente", Color: "#ff0000"}},
	}

	response := ListResponseFromDomain(list)

	require.Len(t, response.Tasks, 1)
	require.Equal(t, int64(10), response.Tasks[0].ID)
	require.Equal(t, []TagResponse{{ID: 20, Name: "urgente", Color: "#ff0000"}}, response.Tags)
}

// Cenário: ListResponsesFromDomain converte todos os itens e retorna slice vazio (não nil) quando não há listas.
func TestListResponsesFromDomain(t *testing.T) {
	responses := ListResponsesFromDomain([]types.List{{ID: 1}, {ID: 2}})
	empty := ListResponsesFromDomain(nil)

	require.Len(t, responses, 2)
	require.Equal(t, int64(2), responses[1].ID)
	require.NotNil(t, empty)
	require.Empty(t, empty)
}
