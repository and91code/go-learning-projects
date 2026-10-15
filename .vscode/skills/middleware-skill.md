# Skill: Criação de Middlewares no Gin (Go)

## Diretrizes de Implementação

Esta skill orienta a construção de middlewares seguros, performáticos e idiomáticos para o framework Gin.

### 1. Interrupção de Requisições (Abort)

- Quando uma validação falhar (ex: Token JWT ausente, expirado ou inválido), utilize imediatamente `c.AbortWithStatusJSON(http.StatusX, gin.H{"error": "mensagem"})`.
- O uso de `AbortWithStatusJSON` garante que os próximos handlers da rota **não** serão executados.

### 2. Injeção de Dados no Contexto

- Armazene dados extraídos ou validados no middleware (como `userID` ou permissões) utilizando constantes tipadas para as chaves do contexto.
- Evite strings mágicas: utilize `c.Set(stringKey, value)` e recupere com segurança utilizando asserção de tipo no handler.

### 3. Práticas de Segurança

- Sempre valide o algoritmo de assinatura de tokens (ex: forçar HMAC `HS256` em tokens JWT).
- Valide a estrutura e o prefixo do cabeçalho de autorização (ex: verificar se começa com `Bearer `).

### Exemplo de Padrão:

```go
const UserIDKey = "userID"

func CustomAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
            return
        }

        // Lógica de validação do token...
        userID := int64(123) // extraído do token

        c.Set(UserIDKey, userID)
        c.Next() // Prossegue para o próximo handler
    }
}
```
