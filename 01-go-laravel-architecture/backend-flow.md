# Fluxo do backend: arquitetura em camadas com Services

Este documento descreve o processamento de requisições no projeto
`01-go-laravel-architecture`. O backend organiza a lógica de aplicação em
services por domínio. Cada service reúne as operações desse domínio e coordena
repositories que acessam o MySQL. As operações são síncronas: a requisição
aguarda a execução até receber a resposta.

## Visão geral

```text
Cliente HTTP
  -> Gin Router e middlewares
  -> Handler
  -> DTO / tipos de domínio
  -> Service do domínio
  -> Repository
  -> MySQL
  -> resposta HTTP
```

## Inicialização e injeção de dependências

1. `backend/main.go` carrega a configuração e abre a conexão com o MySQL.
2. `internal/routes/routes.go` configura o Gin, registra middlewares globais e
   monta o grupo de endpoints `/api/v1`.
3. `internal/module/module.go` instancia os repositories e os services e injeta
   cada service no handler correspondente.
4. Os registradores de rotas encaminham as requisições aos handlers.

O projeto não usa Redis, fila de jobs, worker ou SSE neste fluxo.

## Fluxo de leitura

Exemplo: `GET /api/v1/tasks?list_id=12`.

1. O middleware JWT valida o token e coloca o `userID` no contexto Gin.
2. O handler valida o parâmetro `list_id`.
3. O handler chama a operação de consulta do `TaskService`.
4. O service aplica as regras do domínio e chama o repository.
5. O repository executa a consulta SQL no MySQL.
6. O resultado é convertido para DTO e retornado ao cliente como JSON.

```text
GET /api/v1/tasks?list_id=12
  -> AuthMiddleware
  -> TaskHandler.IndexByList
  -> TaskService.GetTasksByList(ctx, listID, userID)
  -> TaskRepository.FindByListIDAndUserID
  -> MySQL
  <- HTTP 200 + JSON de tarefas
```

Listas e tags seguem o mesmo percurso: handler → service → repository → MySQL.
Registro e login são endpoints públicos; as rotas dos recursos e do perfil
exigem autenticação.

## Fluxo de escrita

Exemplo: `POST /api/v1/tasks` com os dados de uma tarefa.

1. O middleware autentica o usuário.
2. O handler valida o JSON, converte o DTO para o tipo de domínio e recupera o
   `userID`.
3. O handler chama a operação de criação do `TaskService`.
4. O service valida título, prioridade, tags e pertencimento da lista; aplica
   valores iniciais de status e timestamps.
5. O service solicita a persistência ao `TaskRepository`.
6. O repository grava a tarefa e seus vínculos no MySQL.
7. O handler transforma o resultado em DTO e responde após a gravação.

```text
POST /api/v1/tasks + JSON
  -> TaskHandler.Store
  -> TaskService.CreateTask(ctx, userID, task)
  -> regras de domínio
  -> TaskRepository.Create
  -> MySQL
  <- HTTP 201 + DTO da tarefa criada
```

Atualizações e exclusões percorrem o mesmo caminho síncrono, chamando a
operação correspondente do service e do repository. Se a validação ou
persistência falhar, o handler converte o erro para o código HTTP apropriado.

## Formato de entrada e saída

Os DTOs em `backend/internal/dto/` definem os formatos JSON da API e convertem
as entradas em modelos de domínio de `backend/internal/types/`.

Exemplo simplificado para criar uma tarefa:

```json
{
  "list_id": 12,
  "title": "Preparar apresentação",
  "description": "Revisar os slides",
  "priority": "medium",
  "due_date": null,
  "tag_ids": []
}
```

A gravação concluída retorna a representação pública do recurso em JSON
(normalmente HTTP `201 Created` para criação, `200 OK` para atualização e
`204 No Content` para exclusão). Leituras retornam HTTP `200 OK`; erros de
validação, autenticação ou persistência são traduzidos pelo handler para
respostas HTTP de erro.

## Organização por camada

| Camada | Responsabilidade |
| --- | --- |
| `internal/routes` | Define URLs, métodos HTTP e grupos de rotas protegidos ou públicos. |
| `internal/middleware` | Valida JWT e disponibiliza a identidade do usuário. |
| `internal/handler` | Trata HTTP, valida/transforma DTOs, chama services e monta respostas. |
| `internal/dto` | Define os formatos JSON de entrada e saída. |
| `internal/service` | Agrupa operações de negócio por domínio e coordena repositories. |
| `internal/repository` | Encapsula consultas e comandos SQL. |
| `internal/database` | Abre e configura a conexão com MySQL. |
| `internal/module` | Instancia e injeta repositories, services e handlers. |

## Arquivos principais

- Inicialização: `backend/main.go`
- Rotas: `backend/internal/routes/`
- Injeção de dependências: `backend/internal/module/module.go`
- Handlers HTTP: `backend/internal/handler/`
- Services por domínio: `backend/internal/service/`
- Persistência: `backend/internal/repository/`
- DTOs: `backend/internal/dto/`
- Modelos, contratos e erros: `backend/internal/types/`

## Comparação com os projetos seguintes

Este projeto concentra as operações de cada domínio em services. Em
`02-go-useCase-architecture`, essas operações são separadas em casos de uso
atômicos com `Execute`, mantendo o percurso síncrono. Em
`03-go-redis-architecture`, as leituras continuam síncronas e as gravações de
listas, tarefas e tags são enfileiradas para processamento assíncrono por um
worker Redis.
