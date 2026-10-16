# Handler

O diretório `internal/handler` contém a camada de entrada HTTP da aplicação. Seus handlers recebem requisições, interpretam os dados enviados pelo cliente, chamam interfaces de casos de uso específicos e transformam o resultado em respostas HTTP.

## O que pertence a este diretório

- handlers das operações HTTP, organizados por contexto;
- leitura e vinculação de parâmetros, contexto e payloads;
- chamada aos casos de uso apropriados;
- conversão de erros em códigos e corpos de resposta HTTP;
- envio de respostas usando os recursos do Gin.

Atualmente, há handlers para autenticação, listas e tarefas.

## Regras de uso

- Mantenha aqui as preocupações do protocolo HTTP, como status, parâmetros, JSON e contexto da requisição.
- Encaminhe cada operação ao seu caso de uso; evite concentrar regras de negócio nos handlers ou depender de serviços monolíticos.
- Use DTOs para estruturar os dados recebidos e enviados pela API e seus mappers para converter entre `internal/dto` e `internal/types`.
- Injete interfaces de caso de uso específicas.
- Não execute queries SQL nem acesse o banco diretamente.
- Trate erros retornados pelos serviços de forma explícita e converta-os para respostas HTTP adequadas.
- Registre rotas na camada de rotas, sem misturar o registro de endpoints com a implementação dos handlers.

## Fluxo no projeto

Uma rota encaminha a requisição ao handler. O handler interpreta a entrada e chama o caso de uso da operação. Em seguida, converte o resultado ou erro em uma resposta HTTP para o cliente.

## Exemplos de arquivos

- `auth_handler.go` recebe os dados de cadastro e login e retorna os resultados de autenticação.
- `list_handler.go` trata operações HTTP relacionadas a listas e encaminha a identidade do usuário ao serviço.

Em resumo, este diretório responde à pergunta: **“como a aplicação recebe uma requisição HTTP e devolve uma resposta?”**
