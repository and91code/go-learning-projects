# Service

O diretório `internal/service` contém os casos de uso da aplicação. Seus serviços coordenam as operações solicitadas pelos handlers, aplicam regras de negócio e usam os repositórios para consultar ou persistir dados.

## O que pertence a este diretório

- operações como cadastro e login, criação e consulta de listas e criação de tarefas;
- regras que determinam se uma operação é válida para o negócio;
- coordenação entre entidades, DTOs e interfaces de repositório;
- tradução dos resultados da operação para os dados que serão devolvidos ao handler.

Atualmente, os serviços de autenticação, listas e tarefas estão organizados em arquivos próprios.

## Regras de uso

- Mantenha a lógica dos casos de uso aqui, em vez de concentrá-la nos handlers ou nos repositórios.
- Receba dependências por construtores e use repositórios para acesso aos dados.
- Propague erros para que a camada HTTP decida como apresentá-los ao cliente.
- Não manipule diretamente o contexto Gin, status HTTP ou respostas JSON.
- Não execute queries SQL nem acesse a conexão com o banco diretamente.
- Preserve as responsabilidades das camadas existentes: o handler cuida do protocolo HTTP e o repositório cuida da persistência.

## Fluxo no projeto

O handler valida e encaminha a operação ao serviço. O serviço aplica as regras e coordena o repositório. O resultado retorna ao handler, que monta a resposta HTTP.

## Exemplos de arquivos

- `auth_service.go` coordena cadastro e login, incluindo verificação de credenciais e emissão de token.
- `list_service.go` coordena as operações de criação, consulta, atualização e exclusão de listas.

Em resumo, este diretório responde à pergunta: **“o que o sistema deve fazer nesta operação?”**
