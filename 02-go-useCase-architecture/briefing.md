# 📋 TaskFlow API & Web App

<p align="center">
  <img src="https://img.shields.io/badge/version-1.0.0-blue?style=for-the-badge" alt="Version">
  <img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/laravel-%23FF2D20.svg?style=for-the-badge&logo=laravel&logoColor=white" alt="Laravel">
  <img src="https://img.shields.io/badge/react-%2320232a.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB" alt="React">
  <img src="https://img.shields.io/badge/mysql%208.0-%2300f.svg?style=for-the-badge&logo=mysql&logoColor=white" alt="MySQL">
  <img src="https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
</p>

---

## 🎯 Propósito do Projeto

O **TaskFlow API & Web App** é uma aplicação full-stack de gerenciamento de tarefas (_Todo List_) desenvolvida com backend em **Go** e frontend em **React**.

O objetivo principal do projeto é servir como um ambiente prático para a **assimilação da linguagem Go**, aplicando o mindset e estabelecendo paralelos arquiteturais diretos com o ecossistema **Laravel** (MVC, Injeção de Dependência, DTOs e divisão em camadas).

---

## 🛠️ Tech Stack

| Camada             | Tecnologia                | Detalhes                                                           |
| :----------------- | :------------------------ | :----------------------------------------------------------------- |
| **Backend**        | `Go (Golang)`             | Utilizando driver nativo `go-sql-driver/mysql`                     |
| **Frontend**       | `React.js`                | Interface reativa e componentizada                                 |
| **Database**       | `MySQL 8.0`               | Banco de dados relacional com suporte a FKs e operações em cascata |
| **Infraestrutura** | `Docker & Docker Compose` | Padronização e isolamento do ambiente de desenvolvimento           |

---

## 📐 Padrões Arquiteturais & Paralelos com Laravel

A API adota uma **Layered Architecture** (_Handler → Service → Repository_) com isolamento estrito entre **DTOs de HTTP** (_Request/Response_) e **Entidades de Domínio/Banco de Dados**.

```text
  [ Client Request ]
         │
         ▼
 ┌───────────────┐  ──► Equivalente aos Controllers + Form Requests
 │    Handler    │      (Validação de payload e extração de contexto HTTP)
 └───────┬───────┘
         │ DTOs
         ▼
 ┌───────────────┐  ──► Equivalente aos Services / Actions
 │    Service    │      (Regras de negócio e orquestração de casos de uso)
 └───────┬───────┘
         │ Domain Entities
         ▼
 ┌───────────────┐  ──► Equivalente às Eloquent Queries / DB Layer
 │  Repository   │      (Persistência e comunicação direta com o MySQL)
 └───────────────┘
```
