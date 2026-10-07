# 🧪 Todo Apps Playground & Go Learning Journey

Este repositório centraliza meus projetos de estudos e práticas focados em **Go (Golang)**, **React** e arquitetura de software. O objetivo principal é consolidar conceitos de backend, padrões de projeto e APIs robustas, aplicando a bagagem e o mindset do ecossistema **Laravel** ao ecossistema **Go**.

---

## 🎯 Objetivos do Repositório

- **Transição e Assimilação:** Desenvolver soluções em Go comparando estruturas e padrões com o Laravel (MVC, Dependency Injection, ORMs e Migrations).
- **Prática Progressiva:** Criar variações de aplicações Todo List evoluindo a complexidade arquitetural a cada projeto.
- **Portfólio de Evolução:** Registrar o histórico de aprendizado e tomadas de decisão técnicas.

---

## 📁 Estrutura do Monorepo

```text
.
├── 01-go-laravel-architecture/   # Projeto inicial: API REST Go + React
│   ├── backend/                  # API em Go (Rotas, Handlers, Models, Database)
│   └── frontend/                 # Client em React
│
├── 02-go-clean-architecture/     # (Futuro) Aplicação de Clean Arch / DDD em Go
└── README.md                     # Visão geral do repositório
```

---

## 🧩 Entendimento da arquitetura de domínio

A camada `internal/domain` é o coração conceitual da aplicação. Ela define o modelo de negócio e reúne as entidades que representam o que o sistema entende como realidade.

### O que ela representa

Em projetos como este, o domínio concentra a linguagem do negócio, por exemplo:

- usuário
- lista de tarefas
- tarefa
- tag
- status e prioridade

Essa camada serve como base para as outras camadas do sistema, como DTOs, serviços e repositórios.

### Regras de uso do diretório

O diretório de domínio deve ser usado para expressar o problema de negócio, não detalhes técnicos.

Ele deve conter:

- entidades do negócio
- relacionamentos entre entidades
- enums e valores permitidos do sistema
- dados essenciais ao ciclo de vida das informações

Ele deve evitar:

- acesso direto ao banco
- queries SQL
- validação de requisições HTTP
- lógica de framework
- código voltado a infraestrutura ou logging

### Visão prática

O domínio responde à pergunta: “o que é isso no sistema?”

E não: “como isso é salvo?”, “como isso chega na API?” ou “como isso é consultado no banco?”

Essa separação é importante porque permite que a aplicação mantenha uma linguagem clara e um modelo estável, mesmo quando camadas externas mudam.

### Exemplo de modelo de negócio

O domínio do projeto representa conceitos como:

- um usuário que possui listas
- uma lista que contém várias tarefas
- uma tarefa que tem status, prioridade e possível data de vencimento
- uma tag que organiza visualmente ou categoriza tarefas

Esses relacionamentos formam a base da aplicação e ajudam a manter o código coerente entre backend, regras de negócio e persistência.

### Regra de ouro

O conteúdo do domínio deve ser enxuto, estável e orientado ao negócio.

Se uma informação não pertence ao núcleo da aplicação, ela não deve estar ali.

---
