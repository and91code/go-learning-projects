# DTO

O diretório `internal/dto` reúne os objetos usados exclusivamente na fronteira HTTP. DTO significa *Data Transfer Object*: cada tipo descreve os dados esperados em uma entrada ou devolvidos em uma resposta. Os mappers convertem esses objetos para/de entidades e contratos de `internal/types`.

## O que pertence a este diretório

- estruturas de entrada para operações da API;
- estruturas de resposta expostas pela API;
- metadados de serialização JSON e validação de payload;
- mappers explícitos entre DTOs e tipos internos de `internal/types`.

Atualmente, os DTOs estão agrupados por contexto, como autenticação, listas e tarefas.

## Regras de uso

- Use DTOs para definir o formato de entrada e saída da API; não os trate como entidades de domínio ou modelos de persistência.
- Mantenha aqui os nomes JSON e as regras de validação ligadas ao formato recebido pela API.
- Para `due_date` em tarefas, são aceitos timestamps RFC3339 (por exemplo, `2026-08-01T00:00:00-03:00`) e `YYYY-MM-DD HH:MM:SS` (interpretado como UTC).
- Não coloque acesso ao banco, chamadas HTTP ou orquestração de casos de uso nos DTOs.
- Não exponha dados internos ou sensíveis que não façam parte do contrato público da API.
- Quando os dados de domínio e o formato da resposta forem diferentes, mantenha a transformação explícita entre eles.
- Não declare interfaces ou contratos de serviço/repositório aqui; esses tipos pertencem a `internal/types`.

## Fluxo no projeto

O handler vincula o payload recebido a um DTO, usa seu mapper para produzir o tipo interno e chama o serviço. O resultado interno é convertido em DTO de resposta antes de ser enviado ao cliente.

## Exemplos de arquivos

- `auth_dto.go` descreve os dados de cadastro e login e a resposta de autenticação.
- `list_dto.go` descreve as entradas e respostas relacionadas a listas.

Em resumo, este diretório responde à pergunta: **“em qual formato os dados entram ou saem da API?”**
