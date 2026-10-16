package dto

import (
	"testing"

	"github.com/seu-usuario/taskflow-backend/internal/types"
	"github.com/stretchr/testify/require"
)

// Cenário: CreateTagRequest.ToDomain copia lista, nome e cor.
func TestCreateTagRequest_ToDomain(t *testing.T) {
	tag := CreateTagRequest{ListID: 1, Name: "urgente", Color: "#ff0000"}.ToDomain()

	require.Equal(t, types.Tag{ListID: 1, Name: "urgente", Color: "#ff0000"}, tag)
}

// Cenário: UpdateTagRequest.ToDomain copia lista, nome e cor.
func TestUpdateTagRequest_ToDomain(t *testing.T) {
	tag := UpdateTagRequest{ListID: 2, Name: "casa", Color: "#00ff00"}.ToDomain()

	require.Equal(t, types.Tag{ListID: 2, Name: "casa", Color: "#00ff00"}, tag)
}

// Cenário: resposta resumida expõe apenas id, nome e cor (sem list_id).
func TestTagResponseFromDomain(t *testing.T) {
	response := TagResponseFromDomain(types.Tag{ID: 5, ListID: 1, Name: "urgente", Color: "#ff0000"})

	require.Equal(t, TagResponse{ID: 5, Name: "urgente", Color: "#ff0000"}, response)
}

// Cenário: resposta detalhada inclui também o list_id.
func TestTagDetailResponseFromDomain(t *testing.T) {
	response := TagDetailResponseFromDomain(types.Tag{ID: 5, ListID: 1, Name: "urgente", Color: "#ff0000"})

	require.Equal(t, TagDetailResponse{ID: 5, ListID: 1, Name: "urgente", Color: "#ff0000"}, response)
}

// Cenário: TagResponsesFromDomain converte todos os itens e retorna slice vazio (não nil) sem tags.
func TestTagResponsesFromDomain(t *testing.T) {
	responses := TagResponsesFromDomain([]types.Tag{{ID: 1}, {ID: 2}})
	empty := TagResponsesFromDomain(nil)

	require.Len(t, responses, 2)
	require.Equal(t, int64(2), responses[1].ID)
	require.NotNil(t, empty)
	require.Empty(t, empty)
}
