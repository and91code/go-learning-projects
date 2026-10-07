# DTO

O diretório `internal/dto` reúne os objetos usados para transportar dados entre a API e as camadas internas. DTO significa *Data Transfer Object*: cada tipo descreve os dados esperados em uma entrada ou devolvidos em uma resposta.

## O que pertence a este diretório

- estruturas de entrada para operações da API;
- estruturas de resposta expostas pela API;
- metadados de serialização JSON e validação de payload;
- conversões ou formatação simples necessárias para apresentar respostas.

Atualmente, os DTOs estão agrupados por contexto, como autenticação, listas e tarefas.

## Regras de uso

- Use DTOs para definir o formato de entrada e saída da API; não os trate como entidades de domínio ou modelos de persistência.
- Mantenha aqui os nomes JSON e as regras de validação ligadas ao formato recebido pela API.
- Não coloque acesso ao banco, chamadas HTTP ou orquestração de casos de uso nos DTOs.
- Não exponha dados internos ou sensíveis que não façam parte do contrato público da API.
- Quando os dados de domínio e o formato da resposta forem diferentes, mantenha a transformação explícita entre eles.

## Fluxo no projeto

O handler vincula o payload recebido a um DTO e encaminha os dados ao serviço. O serviço processa a operação e o resultado é representado por um DTO de resposta antes de ser enviado ao cliente.

## Exemplos de arquivos

- `auth_dto.go` descreve os dados de cadastro e login e a resposta de autenticação.
- `list_dto.go` descreve as entradas e respostas relacionadas a listas.

Em resumo, este diretório responde à pergunta: **“em qual formato os dados entram ou saem da API?”**
