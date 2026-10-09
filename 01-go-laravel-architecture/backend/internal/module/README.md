# Module

O diretório `internal/module` compõe o grafo de dependências do backend. `Build(db)` instancia os repositórios, injeta-os nos serviços e retorna os handlers concretos que as rotas registram.

## Fluxo de composição

`database connection → repositories → services → handlers → routes`

O tipo `AppModules` declara os handlers disponíveis para o roteador. Como seus campos são ponteiros para handlers concretos, a chamada de cada rota é verificada estaticamente pelo compilador.

## Regras de uso

- Centralize aqui a construção e a injeção das dependências da aplicação.
- Construa cada repositório uma vez e compartilhe-o com os serviços que dependem dele.
- Retorne os handlers esperados pelo roteador, sem construir dependências dentro das rotas.
- Mantenha as entidades e contratos de serviço/repositório em `internal/types`.
- Não mova `AppModules` para `internal/types`: ele precisa referenciar handlers e isso faria o pacote de tipos depender da camada HTTP.

## Uso no projeto

`routes.SetupRouter` chama `module.Build(db)` e passa `Auth`, `Lists` e `Tasks` aos respectivos registradores de rota. Assim, qualquer incompatibilidade entre o retorno da fábrica e as rotas é detectada na compilação.
