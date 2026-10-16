# Fluxo do backend: Use Cases, Redis e notificações

Este documento descreve o backend do projeto `03-go-redis-architecture`.
O MySQL é a persistência principal de usuários, listas, tarefas e tags. Redis
enfileira gravações assíncronas, armazena notificações e entrega eventos em
tempo real por Pub/Sub e Server-Sent Events (SSE).

## Visão geral

```text
LEITURA:
Cliente -> Gin/Auth -> Handler -> Use Case -> Repository -> MySQL -> resposta

ESCRITA ASSÍNCRONA:
Cliente -> Gin/Auth -> Handler -> Redis Producer -> taskflow:queue
                                                    |
                                                    v
                                      Worker -> Use Case -> Repository -> MySQL
                                        |
                                        +-> notificação persistida no Redis
                                        +-> Pub/Sub user:<user_id>:events
                                                   |
                                                   v
                                           Broker SSE -> navegador
```

## Inicialização e injeção de dependências

1. `backend/main.go` carrega as configurações, conecta ao MySQL e cria o cliente
   Redis usando `REDIS_URL`.
2. A conexão Redis é verificada na inicialização. A aplicação reporta erro se
   não conseguir conectar às dependências necessárias.
3. `internal/routes/routes.go` compõe handlers e casos de uso, inicia o worker
   da fila e o broker SSE em goroutines e registra os endpoints.
4. Ao encerrar o processo, o contexto de cancelamento interrompe as rotinas de
   fundo.
5. `internal/module/module.go` liga repositories aos casos de uso, estes aos
   handlers, e configura Producer, Worker e Broadcaster.

## Fluxo de leitura

Exemplo: `GET /api/v1/tasks?list_id=12`.

As leituras não entram na fila. O handler valida parâmetros e identidade do
usuário, executa o caso de uso de leitura, que consulta o repository e o MySQL,
e responde na própria requisição.

```text
GET /api/v1/tasks?list_id=12
  -> AuthMiddleware
  -> TaskHandler.IndexByList
  -> ListTasks.Execute(ctx, listID, userID)
  -> TaskRepository -> MySQL
  <- HTTP 200 + JSON de tarefas
```

## Fluxo de escrita assíncrona

Exemplo: `PUT /api/v1/tasks/42` com os dados atualizados.

1. O middleware JWT autentica o usuário e deixa seu ID no contexto Gin.
2. O handler valida o payload, converte o DTO e cria um job.
3. O Producer serializa o job como JSON e executa `LPUSH` em
   `taskflow:queue`.
4. Após o Redis aceitar o job, o handler responde `202 Accepted`. Isso indica
   que o trabalho foi aceito, não que a gravação já terminou.
5. O Worker consome a fila com `BRPOPLPUSH`, movendo o job para
   `taskflow:processing`.
6. O dispatcher invoca o caso de uso correspondente. O caso de uso aplica
   regras de negócio e persiste via repository no MySQL.
7. Em sucesso, o worker persiste uma notificação no Redis, publica o evento
   para o usuário e confirma o job removendo-o de `processing`.

```text
PUT /api/v1/tasks/42 + JSON
  -> TaskHandler.Update
  -> Producer.Enqueue(Job)
  -> LPUSH taskflow:queue
  <- HTTP 202 {"status":"accepted"}

BRPOPLPUSH taskflow:queue -> taskflow:processing
  -> Worker / Dispatcher
  -> UpdateTask.Execute(ctx, 42, userID, task)
  -> TaskRepository -> MySQL
  -> grava notificação Redis -> PUBLISH user:<user_id>:events
  -> LREM taskflow:processing (ack)
```

Falhas permanentes e jobs inválidos são enviados para `taskflow:dead-letter`.
Falhas transitórias são retentadas até o limite configurado pelo worker.
Como a persistência no MySQL e a confirmação do job no Redis não são uma
transação distribuída, um job pode ser repetido em algumas interrupções; os
efeitos precisam considerar idempotência.

## Formato do job da fila

O job é serializado em JSON e guardado em uma lista Redis. Campos:

| Campo | Tipo | Uso |
| --- | --- | --- |
| `action` | string | Ação, como `TASK_CREATE`, `LIST_UPDATE` ou `TAG_DELETE`. |
| `user_id` | integer | Usuário autenticado que solicitou a operação. |
| `entity_id` | integer, opcional | ID do registro em update/delete. |
| `payload` | JSON, opcional | Dados serializados necessários para create/update. |
| `attempts` | integer, opcional | Número de tentativas já realizadas. |

Exemplo:

```json
{
  "action": "TASK_UPDATE",
  "user_id": 7,
  "entity_id": 42,
  "payload": {
    "list_id": 12,
    "title": "Preparar apresentação",
    "description": "Revisar os slides",
    "status": "in_progress",
    "priority": "medium",
    "tags": []
  },
  "attempts": 0
}
```

## Formato e ciclo de vida das notificações

Após o caso de uso concluir com sucesso, o worker cria uma mensagem com o
resultado da operação. A notificação é persistida por usuário em um Hash Redis
`taskflow:user:<user_id>:notifications`, indexado pela Sorted Set
`taskflow:user:<user_id>:notifications:index`.

```json
{
  "id": "1760000000000000000-42",
  "message": "Tarefa atualizada com sucesso.",
  "resource_type": "task",
  "resource_id": 42,
  "action": "updated",
  "timestamp": "2026-10-15T12:00:00Z",
  "unread": true
}
```

`resource_type` informa a área (`list`, `task` ou `tag`); `resource_id` é o ID
persistido pelo caso de uso para criação, ou o ID alvo para atualização e
exclusão. O cliente pode usar esses campos para abrir a área e destacar o
registro correspondente. Notificações novas começam com `unread: true`.

Os endpoints autenticados permitem:

| Método e caminho | Ação |
| --- | --- |
| `GET /api/v1/notifications` | Listar notificações do usuário. |
| `PATCH /api/v1/notifications/:id/read` | Marcar uma notificação como lida. |
| `DELETE /api/v1/notifications/:id` | Excluir uma notificação. |

A bolinha de notificação no frontend representa a existência de pelo menos uma
notificação não lida.

## Pub/Sub e SSE

O Broadcaster publica o corpo JSON da notificação no canal Redis
`user:<user_id>:events`. O broker em `internal/transport/sse` assina
`user:*:events`, distribui cada evento somente aos clientes SSE locais daquele
usuário, sem bloquear conexões lentas.

O navegador abre `GET /api/v1/events?token=<JWT>` porque o `EventSource` padrão
não permite definir cabeçalho `Authorization`. Essa rota valida o token e o
servidor mantém a resposta como `text/event-stream`. Cada mensagem é enviada
com o evento `taskflow`; comentários periódicos mantêm a conexão ativa.

Pub/Sub é entrega em tempo real, não histórico durável. O frontend carrega o
histórico pelo endpoint de notificações e usa SSE para novas mensagens.

Ao receber uma notificação, o frontend atualiza somente o contexto relacionado:

- `resource_type: "list"` recarrega listas.
- `resource_type: "task"` recarrega tarefas.
- `resource_type: "tag"` recarrega tags.

## Camadas e principais arquivos

| Responsabilidade | Localização |
| --- | --- |
| Configuração, conexão MySQL/Redis e shutdown | `backend/main.go` |
| Rotas e inicialização das rotinas de fundo | `backend/internal/routes/` |
| Construção do grafo e worker | `backend/internal/module/module.go` |
| Handlers e enfileiramento HTTP | `backend/internal/handler/` |
| Contrato do job e Producer | `backend/internal/queue/producer.go` |
| Worker, dispatcher, retries e dead-letter | `backend/internal/queue/worker.go` |
| Persistência relacional | `backend/internal/repository/` |
| Regras da aplicação | `backend/internal/usecase/` |
| Publicação Redis Pub/Sub | `backend/internal/pubsub/broadcaster.go` |
| Distribuição de SSE por usuário | `backend/internal/transport/sse/broker.go` |
| Formato HTTP de entrada e saída | `backend/internal/dto/` |

## Comparação com `02-go-useCase-architecture`

O projeto `02` usa o mesmo princípio de casos de uso atômicos, mas processa
leituras e gravações de forma síncrona: o cliente aguarda o use case e o MySQL,
e a gravação responde com o resultado da operação. O projeto `03` preserva o
caminho síncrono de leitura e passa gravações de listas, tarefas e tags por uma
fila Redis; por isso, as gravações respondem `202 Accepted` e o resultado final
é comunicado pelo worker e por notificações/SSE.
