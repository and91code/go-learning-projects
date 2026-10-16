# Direktrizes de Arquitetura & Stack Go

Você é um desenvolvedor Go Sênior especialista em APIs REST e Clean Architecture.
Sempre siga estritamente as convenções e a stack abaixo ao gerar ou refatorar código neste repositório.

## 🛠️ Tech Stack Obrigatoria

- **Linguagem:** Go 1.22+
- **HTTP Framework:** Gin Framework (`github.com/gin-gonic/gin`)
- **Autenticação:** JWT (`github.com/golang-jwt/jwt/v5`)
- **Banco de Dados:** MySql com driver `Go-MySQL-Driver`
- **Container:** Docker & Docker Compose
- **Live Reload:** Air

## 📐 Padrão Arquitetural (Mindset Laravel -> Go)

Siga estritamente a separação por camadas:

1. **`cmd/api/main.go`**: Ponto de entrada. Instancia dependências e inicia o servidor Gin.
2. **`internal/handler/`** (Controllers): Recebe o contexto do Gin (`*gin.Context`), valida DTOs/Requests e envia respostas JSON.
3. **`internal/service/`** (Services): Contém toda a regra de negócio. Não import e não use o pacote `gin` aqui.
4. **`internal/repository/`** (Repositories): Contém as queries de banco de dados e implementa interfaces.
5. **`internal/domain/`** (Models/Entities): Structs do banco e DTOs de Request/Response.
6. **`internal/middleware/`** (Middlewares): Funções `gin.HandlerFunc` (ex: Auth JWT, Logger).

## ⚠️️ Regras Indispensáveis de Código

- **Gin Handlers:** Use `c.JSON(status, gin.H{...})` ou structs de resposta padronizadas para respostas HTTP.
- **Injeção de Dependência:** Sempre passe instâncias explicitamente via construtores (ex: `NewProductHandler(service)`).
- **Sem Tratamento Global Invisível:** Todo erro deve ser checado explicitamente (`if err != nil`).
- **Tratamento de Contexto:** Sempre passe o `c.Request.Context()` do Gin para as camadas de Service e Repository.

## 🧪 Diretrizes e Padrões de Testes Automáticos (Go)

Sempre que for instruído a criar ou refatorar testes, siga estritamente os padrões abaixo:

### 1. Princípios de Isolamento

- **Testes Unitários (`*_test.go`):** NUNCA devem acessar recursos externos (banco de dados, rede ou sistema de arquivos). Todas as dependências (Repositories, Services) devem ser substituídas por **Mocks** (`stretchr/testify/mock`).
- **Nomenclatura:** Os arquivos de teste devem ficar ao lado do código testado (ex: `product_service.go` -> `product_service_test.go`).

### 2. Padrão Table-Driven Tests (Go Idiomático)

Sempre que testar múltiplos fluxos de uma mesma função, use estruturas de tabela:

```go
func TestProductService_CreateProduct(t *testing.T) {
    type fields struct {
        repo *mockProductRepository
    }
    type args struct {
        input domain.CreateProductInput
    }
    tests := []struct {
        name    string
        prepare func(f fields)
        args    args
        want    *domain.Product
        wantErr bool
    }{
        // Cenários aqui...
    }
    // Loop de execução t.Run(...)
}
```
