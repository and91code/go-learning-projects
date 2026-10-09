# Domain (legado)

As entidades e contratos do domínio foram centralizados em [`../types`](../types/README.md) para que repositórios, serviços e handlers compartilhem as mesmas tipagens.

Este diretório não deve receber novas definições de entidades, enums ou interfaces. Ao criar ou alterar tipos compartilhados, use `internal/types`; os DTOs continuam em `internal/dto` e descrevem exclusivamente a entrada e a saída HTTP.
