# Service

O diretório `internal/service` contém os casos de uso da aplicação. Seus serviços coordenam as operações solicitadas pelos handlers, aplicam regras de negócio e usam as interfaces de repositório de `internal/types` para consultar ou persistir dados.

## O que pertence a este diretório

- operações como cadastro e login, criação e consulta de listas e criação de tarefas;
- regras que determinam se uma operação é válida para o negócio;
- coordenação entre entidades, DTOs e interfaces de repositório;
- retorno de entidades e valores internos definidos em `internal/types`.

Atualmente, os serviços de autenticação, listas e tarefas estão organizados em arquivos próprios.

## Regras de uso

- Mantenha a lógica dos casos de uso aqui, em vez de concentrá-la nos handlers ou nos repositórios.
- Receba dependências por construtores e use repositórios para acesso aos dados.
- Use exclusivamente os tipos e contratos internos de `internal/types`; não receba ou retorne DTOs.
- Propague erros para que a camada HTTP decida como apresentá-los ao cliente.
- Não manipule diretamente o contexto Gin, status HTTP ou respostas JSON.
- Não execute queries SQL nem acesse a conexão com o banco diretamente.
- Preserve as responsabilidades das camadas existentes: o handler cuida do protocolo HTTP e o repositório cuida da persistência.

## Fluxo no projeto

O handler valida o DTO e converte a entrada para `internal/types`. O serviço aplica as regras e coordena o repositório usando contratos de `internal/types`. O resultado interno retorna ao handler, que converte para o DTO da resposta HTTP.

## Exemplos de arquivos

- `auth_service.go` coordena cadastro e login, incluindo verificação de credenciais e emissão de token.
- `list_service.go` coordena as operações de criação, consulta, atualização e exclusão de listas.

Em resumo, este diretório responde à pergunta: **“o que o sistema deve fazer nesta operação?”**
