# Skill: Repositório Go e Manipulação SQL com Database/sql

## Diretrizes de Implementação

Esta skill define o padrão obrigatório para a escrita de repositórios persistentes utilizando Go puro e `database/sql` (sem ORMs pesados).

### 1. Assinatura e Contexto

- Todo método de repositório deve receber `ctx context.Context` como o primeiro argumento.
- Propague o contexto nas chamadas de driver de banco (`ExecContext`, `QueryContext`, `QueryRowContext`).

### 2. Tratamento de Campos Anuláveis (Null Types)

- Bancos relacionais frequentemente possuem colunas que aceitam `NULL`. No Go, utilize os tipos do pacote `database/sql`:
  - `sql.NullString` para strings opcionais.
  - `sql.NullTime` para datas e horários opcionais.
- Sempre verifique a propriedade `.Valid` antes de atribuir o valor à struct de domínio.

### 3. Tratamento de Erros e Linhas Afetadas

- Sempre encapsule os erros retornados utilizando `fmt.Errorf` com o verbo `%w` para preservar a árvore de erros:
  `return fmt.Errorf("taskRepository.FindByID: %w", err)`
- Mapeie `sql.ErrNoRows` para um erro de domínio personalizado (ex: `ErrTaskNotFound`).
- Em operações de `UPDATE` ou `DELETE`, valide o retorno de `RowsAffected()` para garantir que a linha realmente existia.

### Exemplo de Padrão:

```go
func (r *mysqlTaskRepository) FindByID(ctx context.Context, id int64) (*domain.Task, error) {
    query := `SELECT id, title, description, due_date FROM tasks WHERE id = ?`

    var t domain.Task
    var desc sql.NullString
    var dueDate sql.NullTime

    err := r.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.Title, &desc, &dueDate)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrTaskNotFound
        }
        return nil, fmt.Errorf("taskRepository.FindByID: %w", err)
    }

    if desc.Valid { t.Description = desc.String }
    if dueDate.Valid { t.DueDate = &dueDate.Time }

    return &t, nil
}
```
