# Skill: Go Layered Architecture & Laravel Parallels

## Descrição

Instrui o agente a gerar e organizar código backend em Go seguindo o padrão de 3 camadas com injeção de dependência explícita.

## Mapeamento de Camadas (Laravel vs. Go)

| Camada no Go           | Equivalente no Laravel   | Responsabilidade                                                         |
| :--------------------- | :----------------------- | :----------------------------------------------------------------------- |
| **Handler / Delivery** | Controller               | Recebe requisição HTTP, valida JSON/DTO, chama o Service e retorna JSON. |
| **Service / UseCase**  | Service Class / Action   | Contém as regras de negócio, transações de banco e orquestração.         |
| **Repository**         | Eloquent / Model Queries | Executa queries SQL diretas ou chamadas de ORM no MySQL.                 |

## Padrão de Implementação em Go

### 1. Injeção de Dependência Manual

Sempre defina dependências através de `structs` e interfaces no construtor:

```go
// UserHandler depende da interface do UserService
type UserHandler struct {
    service UserService
}

func NewUserHandler(service UserService) *UserHandler {
    return &UserHandler{service: service}
}
```

## Isolamento de DTOs (Data Transfer Objects) e Mapeamento de Camadas

### Regras de DTOs

1. **Nunca exponha entidades do banco diretamente na API:**
   - As `structs` do banco de dados (Repository) não devem conter tags de serialização JSON da resposta final.
   - Crie `structs` de Request e Response dedicadas no pacote de Handlers/HTTP.

2. **Validação de Entrada (Request DTOs):**
   - Mapeie dados de entrada com `json` tags e tags de validação (`validate`).
   - Equivalente aos `FormRequests` do Laravel.

3. **Mapeamento Explícito (Sem Mágica):**
   - Realize a conversão de DTO -> Entidade de Domínio e Entidade de Domínio -> DTO de forma explícita via métodos ou funções auxiliares.
   - Evite mapeadores reflexivos/automáticos.

### Exemplo do Padrão em Go

```go
// 1. DTO de Entrada (Handler/HTTP Layer)
type CreateTaskRequest struct {
    Title    string    `json:"title" validate:"required,min=3,max=100"`
    DueDate  time.Time `json:"due_date" validate:"required"`
    Priority string    `json:"priority" validate:"required,oneof=low medium high"`
}

// 2. DTO de Saída / Resposta (Handler/HTTP Layer)
type TaskResponse struct {
    ID        int64     `json:"id"`
    Title     string    `json:"title"`
    Status    string    `json:"status"`
    Priority  string    `json:"priority"`
    DueDate   string    `json:"due_date"`
    CreatedAt string    `json:"created_at"`
}

// 3. Função de Mapeamento (DTO -> Response)
func NewTaskResponse(task domain.Task) TaskResponse {
    return TaskResponse{
        ID:        task.ID,
        Title:     task.Title,
        Status:    string(task.Status),
        Priority:  string(task.Priority),
        DueDate:   task.DueDate.Format(time.RFC3339),
        CreatedAt: task.CreatedAt.Format(time.RFC3339),
    }
}
```
