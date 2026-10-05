# Skill: Go Idiomatic Error Handling

## Descrição

Aplica as melhores práticas de tratamento de erros nativo em Go, com embrulho de erros (error wrapping) e equivalência ao manuseio de exceções do Laravel.

## Diretrizes de Código

1. **Tratamento Explícito:**
   - Todo método/função que retorna `error` deve ser verificado imediatamente.
   - Evite `panic()` em fluxo normal da API. Use `panic()` apenas na inicialização crítica da aplicação.

2. **Error Wrapping (Go 1.13+):**
   - Sempre adicione contexto ao erro antes de propagá-lo para a camada superior.
   - Use `%w` com `fmt.Errorf`:
     ```go
     if err != nil {
         return fmt.Errorf("repository.FindUserByID: %w", err)
     }
     ```

3. **Custom Domain Errors (Equivalente a Exceptions customizadas no Laravel):**
   - Defina erros de domínio reutilizáveis em pacotes apropriados:
     ```go
     var (
         ErrUserNotFound = errors.New("usuário não encontrado")
         ErrEmailInUse   = errors.New("e-mail já está em uso")
     )
     ```
   - Use `errors.Is(err, ErrUserNotFound)` na camada de Handler HTTP para retornar o status code correto (ex: HTTP 404).

4. **Validação de Inputs:**
   - Mapeie erros de validação para HTTP 422 (Unprocessable Entity), equivalente aos `FormRequests` do Laravel.
