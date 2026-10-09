# Types

O pacote `internal/types` é a fonte única para entidades, enums, erros compartilhados e contratos entre repositórios e serviços. Ele representa os dados e as operações internas do sistema, sem depender de DTOs, handlers, banco de dados ou frameworks HTTP.

## Organização

- `user.go`: usuário, dados internos de autenticação e contratos de usuário/autenticação.
- `list.go`: entidade e contratos de lista.
- `task.go`: entidades de tarefa e tag, enums de status/prioridade e contratos de tarefa.
- `errors.go`: erros de domínio reconhecidos entre as camadas.

## Fluxo dos contratos

O repositório recebe e retorna as entidades deste pacote. O serviço consome as interfaces de repositório e expõe as interfaces de serviço definidas aqui. O handler converte os DTOs HTTP para tipos internos antes de chamar o serviço e converte os resultados internos de volta para DTOs.

## Regras de uso

- Defina cada entidade, enum e interface compartilhada em `internal/types`; não replique esses contratos em `repository`, `service` ou `handler`.
- Contratos de repositório e serviço devem usar entidades/tipos deste pacote, nunca tipos HTTP.
- Não importe `dto`, `handler`, Gin ou dependências de banco neste pacote.
- Mantenha validação JSON e conversores de entrada/saída HTTP em `internal/dto`.
- Mantenha SQL e detalhes de persistência em `internal/repository`.
- Mantenha a composição concreta de handlers em `internal/module`. O contêiner de aplicação fica nesse pacote porque depende dos handlers; colocá-lo em `types` criaria dependência circular e acoplaria os tipos internos à camada HTTP.

## Regra de referência

Para mudanças no modelo interno ou nos contratos entre camadas, este diretório é a fonte de verdade. Para o formato público da API, a fonte de verdade permanece em `internal/dto`.
