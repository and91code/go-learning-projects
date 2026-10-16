# Fluxo do backend: Use Cases atômicos

Este documento descreve o processamento de requisições no projeto
`02-go-useCase-architecture`. A arquitetura separa as ações da aplicação em
casos de uso atômicos, mantendo a execução das operações síncrona: a requisição
aguarda o caso de uso, o repository e o MySQL antes de receber a resposta.

## Visão geral

```text
Cliente HTTP
  -> Gin Router e middlewares
  -> Handler
  -> DTO / tipos de domínio
  -> Use Case específico (Execute)
  -> Repository
  -> MySQL
  -> resposta HTTP
```

## Inicialização e injeção de dependências

1. `backend/main.go` carrega a configuração e abre a conexão com o MySQL.
2. `internal/routes/routes.go` configura o Gin, registra middlewares globais e
   monta o grupo de rotas `/api/v1`.
3. `internal/module/module.go` cria os repositories, instancia os casos de uso
   e os injeta nos handlers.
4. As rotas encaminham cada endpoint ao handler correspondente.

O projeto não usa Redis para executar o fluxo de requisições.

## Fluxo de leitura

Exemplo: `GET /api/v1/tasks?list_id=12`.

1. O middleware JWT valida o token e disponibiliza `userID` no contexto.
2. O handler valida e converte `list_id`.
3. O handler chama o caso de uso de listagem de tarefas com `Execute`.
4. O caso de uso chama o repository, que consulta o MySQL.
5. O resultado de domínio é convertido para DTO e retornado como JSON.

```text
GET /api/v1/tasks?list_id=12
  -> AuthMiddleware
  -> TaskHandler.IndexByList
  -> ListTasks.Execute(ctx, listID, userID)
  -> TaskRepository.FindByListIDAndUserID
  -> MySQL
  <- HTTP 200 + JSON de tarefas
```

As rotas de leitura de listas, tarefas e tags seguem o mesmo padrão. As rotas
públicas de autenticação, como registro e login, não exigem o middleware JWT;
rotas de perfil e os recursos do usuário são protegidos.

## Fluxo de escrita

Exemplo: `POST /api/v1/tasks` com uma tarefa.

1. Gin aplica o middleware de autenticação e o handler identifica o usuário.
2. O handler valida o JSON de entrada e converte o DTO para o tipo de domínio.
3. O handler chama o caso de uso atômico adequado.
4. O caso de uso valida regras de negócio e solicita a persistência ao
   repository.
5. O repository executa SQL no MySQL.
6. Após o sucesso, o handler devolve a resposta da operação.

```text
POST /api/v1/tasks + JSON
  -> TaskHandler.Store
  -> CreateTask.Execute(ctx, userID, task)
  -> validação / TaskRepository.Create
  -> MySQL
  <- HTTP 201 + DTO da tarefa criada
```

As operações de listas, tags e tarefas possuem casos de uso separados para as
ações suportadas, incluindo criação, listagem/consulta, atualização e exclusão.
Registro e login também usam casos de uso. Perfil e senha são tratados pelos
casos de uso de autenticação correspondentes.

## Formato de entrada e saída

O cliente envia JSON nos endpoints de criação e atualização. Os DTOs em
`backend/internal/dto/` definem os formatos HTTP e fazem a conversão para tipos
de domínio em `backend/internal/types/`.

Exemplo simplificado de criação de tarefa:

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

Uma consulta retorna JSON com os campos públicos definidos pelo DTO; erros de
validação, autorização ou persistência são convertidos pelo handler em códigos
HTTP e respostas de erro.

## Organização por camada

| Camada | Responsabilidade |
| --- | --- |
| `internal/routes` | Define URLs, métodos HTTP e grupos protegidos/públicos. |
| `internal/middleware` | Valida JWT e extrai a identidade do usuário. |
| `internal/handler` | Valida o protocolo HTTP, chama casos de uso e monta respostas. |
| `internal/dto` | Define formatos JSON da API e converte dados para o domínio. |
| `internal/usecase/<domínio>` | Executa uma ação de aplicação por método `Execute`. |
| `internal/repository` | Encapsula consultas e comandos SQL. |
| `internal/database` | Abre e configura conexão com MySQL. |
| `internal/module` | Monta o grafo de dependências da aplicação. |

## Arquivos principais

- Bootstrap: `backend/main.go`
- Rotas e composição HTTP: `backend/internal/routes/`
- Dependências: `backend/internal/module/module.go`
- Handlers: `backend/internal/handler/`
- Casos de uso: `backend/internal/usecase/`
- Persistência: `backend/internal/repository/`
- Modelos e contratos: `backend/internal/types/`
