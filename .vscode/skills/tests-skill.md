# Skill: Testes Unitários em Go (Testing & Testify)

## Diretrizes de Implementação

Esta skill orienta a criação de testes unitários limpos, focados na camada de **Service** (regras de negócio), utilizando a biblioteca padrão `testing` em conjunto com `github.com/stretchr/testify/assert` e mocks.

### 1. Estrutura dos Testes

- Crie o arquivo de teste no mesmo diretório com o sufixo `_test.go` (ex: `task_service_test.go`).
- Utilize a abordagem baseada em tabelas (_Table-Driven Tests_) para testar múltiplos cenários (sucesso, erros, validações) em uma única função de teste.

### 2. Isolamento com Mocks

- Para testar a camada de serviço sem tocar no banco de dados real, crie structs mockadas que implementem as interfaces do repositório definidas em `internal/types/`.
- Utilize assertions limpas com `assert.NoError(t, err)`, `assert.Equal(t, expected, actual)` e `assert.ErrorIs(t, err, expectedErr)`.

### Exemplo de Padrão (Table-Driven Test):

```go
func TestTaskService_Create(t *testing.T) {
    // Cenários de teste
    tests := []struct {
        name          string
        input         dto.CreateTaskRequest
        mockBehavior  func(m *MockTaskRepo, lr *MockListRepo)
        expectedError error
    }{
        {
            name: "Sucesso ao criar tarefa",
            input: dto.CreateTaskRequest{Title: "Estudar Go", ListID: 1},
            mockBehavior: func(tr *MockTaskRepo, lr *MockListRepo) {
                lr.On("FindByIDAndUserID", mock.Anything, int64(1), int64(10)).Return(&domain.List{}, nil)
                tr.On("Create", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(nil)
            },
            expectedError: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Execução do teste...
        })
    }
}
```

```

```
