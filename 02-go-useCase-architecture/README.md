# 🚀 01 - Architecture: Go & Laravel Mindset

<p align="center">
  <img src="https://img.shields.io/badge/status-in%20development-yellow?style=for-the-badge" alt="Status">
  <img src="https://img.shields.io/badge/version-1.0.0-blue?style=for-the-badge" alt="Version">
  <img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/laravel-%23FF2D20.svg?style=for-the-badge&logo=laravel&logoColor=white" alt="Laravel">
  <img src="https://img.shields.io/badge/react-%2320232a.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB" alt="React">
  <img src="https://img.shields.io/badge/mysql-%2300f.svg?style=for-the-badge&logo=mysql&logoColor=white" alt="MySQL">
  <img src="https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
</p>

---

## 📌 Sobre o Projeto

O **Architecture Go & Laravel Mindset** é um projeto de estudo prático focado na **transição de mentalidade arquitetural** entre o ecossistema opinativo do **Laravel (PHP)** e o modelo explícito/performático do **Go (Golang)**.

A aplicação consiste em um **Gerenciador de Tarefas (Todo List)** completo e robusto, integrando um backend escalável em Go com um frontend moderno em React, utilizando isolamento via Docker.

---

## 🛠️ Tech Stack

| Camada       | Tecnologia         | Descrição                                                                             |
| :----------- | :----------------- | :------------------------------------------------------------------------------------ |
| **Backend**  | `Go (Golang)`      | Linguagem compilada, concorrência nativa, alta performance e simplicidade sintática.  |
| **Frontend** | `React.js`         | Biblioteca reativa para construção de SPA (_Single Page Application_) componentizada. |
| **Database** | `MySQL`            | Banco de dados relacional com forte integridade referencial (Foreign Keys).           |
| **Infra**    | `Docker & Compose` | Padronização e orquestração do ambiente de desenvolvimento em containers.             |

---

## ⚡ Funcionalidades do Sistema

- [x] **Autenticação & Segurança**
  - Cadastro e login de usuários
  - Autenticação via Token JWT
  - Proteção de rotas com Middleware de autorização
  - Atualização autenticada do perfil e da senha
- [x] **Gestão de Listas**
  - CRUD completo de listas por usuário
  - Exclusão em cascata (_Cascade Delete_) para tarefas e tags vinculadas
- [x] **Gestão de Tarefas**
  - CRUD de tarefas vinculadas a uma lista
  - Status da tarefa: `pending` \| `in_progress` \| `completed`
  - Prioridade: `low` \| `medium` \| `high`
  - Definição de data limite (`due_date`) e associação de múltiplas tags
- [x] **Categorização por Tags**
  - Criação de tags por lista
  - Atribuição de cores visuais (Hexadecimal)

---

## 🐳 Ambiente de desenvolvimento com Docker

Na pasta `01-go-laravel-architecture`, inicie o MySQL, o backend e o frontend com:

```sh
docker compose up --build
```

Os arquivos locais `backend/.env` e `frontend/.env` fornecem as configurações de
cada aplicação. Use `backend/.env-sample` e `frontend/.env-sample` como modelo
para criá-los em uma nova instalação.

A aplicação estará disponível em `http://localhost:5173` e a API em
`http://localhost:8080/api/v1`. Para encerrar os serviços, use `docker compose down`.
Os dados do MySQL ficam em `01-go-laravel-architecture/.data/mysql_data`.

O `JWT_SECRET` pode ser configurado no ambiente antes de subir o Compose; se
omitido, será usado um valor padrão apenas para desenvolvimento local.

---

## 🧠 Mindset Architectural: Laravel vs. Go

Abaixo está o mapeamento conceitual de como padrões comuns no Laravel são traduzidos para a arquitetura idiomática em Go:

| Conceito                     | 🔴 Laravel Mindset                                                                                 | 🔵 Go Mindset                                                                                              |
| :--------------------------- | :------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------- |
| **Arquitetura**              | **MVC Opinativo**<br>Convenções rigorosas e automáticas (Artisan, Service Providers).              | **Modular em Camadas**<br>Separação explícita (_Handler/Controller_ → _Service/Usecase_ → _Repository_).   |
| **ORM & Database**           | **Eloquent ORM**<br>Pattern ActiveRecord com relacionamentos dinâmicos e "mágicos".                | **SQL / ORM Leve**<br>SQL puro (`go-sql-driver/mysql`) ou `GORM`, com structs e tipos explícitos.          |
| **Migrations**               | **Artisan Migrations**<br>Nativas do framework (`php artisan migrate`).                            | **CLI Externa**<br>Ferramentas como `golang-migrate` gerenciando arquivos `.sql` (UP/DOWN).                |
| **Roteamento & Middlewares** | **Routes File**<br>Arquivo centralizado (`routes/api.php`) com grupos de middleware.               | **HTTP Routers**<br>Roteadores performáticos (`Chi` / `Gin`) encadeando funções no contexto da requisição. |
| **Validação de Dados**       | **Form Requests**<br>Validação automática por classes dedicadas antes do controller.               | **Struct Tags**<br>Parse de JSON no payload e validação via tags (`go-playground/validator`).              |
| **Injeção de Dependência**   | **Service Container (IoC)**<br>Injeção automática via _Auto-wiring_ no ciclo de vida da aplicação. | **Manual & Explícita**<br>Instanciação direta via construtores e repasse manual de structs.                |

---

## 🏗️ Estrutura Arquitetural em Go

Para manter o código Go limpo e desacoplado, o backend é estruturado da seguinte forma:

```text
.
├── backend/
│   ├── main.go          # Carrega configuração e inicia banco, rotas e servidor
│   └── internal/
│       ├── database/    # Conexão MySQL e configuração do pool
│       ├── logger/      # Logging da aplicação e das requisições HTTP
│       ├── routes/      # Middlewares, injeção de dependências e rotas Gin
│       ├── handler/     # Camada HTTP (Controllers / Delivery)
│       ├── usecase/     # Casos de uso organizados por operação
│       ├── service/     # Implementações de regras de negócio legadas
│       └── repository/  # Camada de acesso ao banco de dados (Persistence)
```

### Testes unitários do backend

Os testes unitários de Service e Handler usam mocks com `testify` e não acessam o banco de dados. Execute os comandos a partir da pasta `backend`:

```sh
go test -tags=unit -v ./...
go test -tags=unit -cover ./...
go test -tags=unit -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Sem `-tags=unit`, os testes marcados não são incluídos.

### Atualização dos dados do usuário

Os dois endpoints exigem `Authorization: Bearer <token>`:

- `PUT /api/v1/user/profile` recebe `{"name":"Ana Silva","email":"ana@example.com"}`.
- `PUT /api/v1/user/password` recebe `{"current_password":"senha-atual","new_password":"nova-senha"}`.

O nome deve ter entre 2 e 100 caracteres, o e-mail deve ser válido e a nova senha
deve ter pelo menos 6 caracteres, seguindo as regras de cadastro.
