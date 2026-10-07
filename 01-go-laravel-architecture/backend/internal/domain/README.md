# Domain

Este diretório representa o núcleo do modelo de negócio da aplicação. Ele reúne as entidades e os tipos que descrevem o que a aplicação realmente entende como usuário, lista, tarefa e tag.

## Objetivo

O pacote `internal/domain` define as estruturas centrais da aplicação, sem depender de HTTP, banco de dados ou camada de apresentação. Sua função é manter um modelo estável do problema de negócio.

Em termos simples: aqui não se decide como a API recebe dados, nem como o banco salva dados. Aqui se define o que existe no sistema e como essas coisas se relacionam.

## O que fica neste diretório

As entidades deste pacote normalmente representam:

- usuários
- listas de tarefas
- tarefas
- etiquetas ou tags
- valores de domínio, como status e prioridade

Esses elementos são a base para as camadas que vêm depois, como DTOs, serviços e repositórios.

## Regras de uso do diretório

### 1) O domínio deve ser focado no negócio

Tudo que estiver aqui deve refletir conceitos do sistema e não detalhes técnicos.

Use este diretório para representar:

- entidades do negócio
- relacionamentos entre essas entidades
- enums e valores permitidos
- campos importantes do ciclo de vida da informação

Evite:

- acesso direto ao banco
- queries SQL
- validação de request HTTP
- lógica de framework
- logs e infraestrutura

### 2) Este pacote deve ser simples e estável

O diretório deve ser legível, previsível e pouco acoplado a tecnologia.

Isso significa que:

- cada arquivo geralmente concentra uma entidade ou um conjunto de tipos relacionados;
- os nomes devem refletir o conceito do domínio;
- campos e enums devem ser coerentes com o modelo do negócio;
- o código deve ser pensado como referência do sistema, não como implementação de transporte de dados.

### 3) O domínio é a base das outras camadas

As demais camadas usam o que está aqui como referência:

- DTOs convertem dados de entrada/saída para o exterior;
- serviços aplicam regras de negócio sobre as entidades;
- repositórios persistem e consultem essas entidades no banco.

Ou seja, o domínio é o “contrato do negócio” da aplicação.

### 4) Relações devem refletir o problema real

As entidades e campos neste pacote devem representar relações que fazem sentido para o sistema, por exemplo:

- uma tarefa pertence a uma lista;
- uma lista pertence a um usuário;
- uma tarefa pode ter várias tags;
- um usuário possui dados de autenticação e criação.

Essas relações ajudam a manter a consistência do modelo e a clareza da arquitetura.

### 5) Campos essenciais e opcionais devem ser bem distinguidos

Neste diretório é comum observar uma separação entre:

- dados obrigatórios do domínio;
- dados opcionais;
- dados de auditoria e timestamps;
- valores que podem estar ausentes em algumas situações.

A ideia é expressar a realidade do negócio de forma clara e segura.

### 6) O domínio não deve ser confundido com infraestrutura

Se um arquivo aqui começa a conhecer detalhes do banco, do framework, do protocolo HTTP ou da infraestrutura da aplicação, ele já está saindo do escopo correto.

O domínio deve responder à pergunta: “o que é isso no sistema?”

E não: “como isso é salvo?”, “como isso chega via HTTP?”, “como isso é consultado no banco?”

## Estrutura atual observada

O diretório atual se organiza em arquivos que representam os principais blocos do modelo:

- `user.go` — entidade de usuário
- `list.go` — entidade de lista de tarefas
- `task.go` — entidade de tarefa
- `tag.go` — entidade de tag
- `task_types.go` — tipos e valores permitidos para status e prioridade

## Exemplos de arquivos internos

### Exemplo 1: entidade de usuário

O arquivo de usuário representa a pessoa que interage com o sistema e guarda os dados mínimos de identidade e criação.

Ele deve conter dados de cadastro, referência de criação e, quando necessário, dados sensíveis tratados com cuidado.

### Exemplo 2: entidade de tarefa

O arquivo de tarefa representa o item principal do fluxo de trabalho da aplicação. Ele normalmente traz informações como:

- título
- descrição
- status
- prioridade
- vínculo com uma lista
- data de vencimento
- timestamps de criação e atualização

Essas informações expressam o ciclo de vida da tarefa dentro do domínio.

## Regra de ouro

O diretório `internal/domain` deve ser entendido como a camada de definição do modelo de negócio.

Ele deve conter:

- entidades de negócio
- tipos de domínio
- regras de consistência semântica
- relacionamentos fundamentais

E deve evitar:

- lógica de persistence
- lógica de transporte
- validação de entrada externa
- detalhes de infraestrutura

Se uma informação não pertence ao núcleo do negócio, ela não deve estar aqui.
