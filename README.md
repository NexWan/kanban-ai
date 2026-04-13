# AI Kanban (Name TBD lol)

An AI-driven Kanban board designed to orchestrate tasks between humans and autonomous agents.

This project explores a Kanban-style task management system where tasks can be delegated to AI agents capable of analyzing, decomposing, and executing work while updating the board in real time.

The goal is to build a **collaborative workspace** where users and AI agents interact through a shared Kanban interface.

---

# Project Goals

* Build a functional Kanban board (boards, columns, cards)
* Allow delegation of tasks to AI agents
* Enable AI to create subtasks and move tasks across states
* Support multi-user collaboration through projects
* Provide extensible architecture for MCP tools and sub-agents
* Explore Go backend + Python AI orchestration architecture

---

# Core Concepts

### Boards

Kanban boards containing columns and cards.

### Columns

Represent workflow states (Backlog, In Progress, Review, Done).

### Cards

Tasks that can be:

* created by users
* assigned to users
* delegated to AI agents
* split into subtasks
* moved across columns

### AI Agents

Agents capable of:

* analyzing tasks
* generating subtasks
* updating board state
* invoking tools (via MCP)
* collaborating with users

### Projects (planned)

Projects group:

* users
* boards
* permissions

This enables collaborative workspaces.

---

# Architecture Overview

Frontend communicates with Go backend, which acts as the source of truth.
AI orchestration is handled by a separate Python service.

```
Frontend (Svelte/Vue/etc)
        |
        v
Go API (Boards / Cards / Columns)
        |
        v
PostgreSQL
        |
        v
Python AI Service (Agents / MCP / Tools)
```

---

# Tech Stack

## Backend

Language:

* Go

Libraries:

* pgx (PostgreSQL driver)
* go-migrate (database migrations)
* chi (router, planned)
* air (hot reload)

Responsibilities:

* boards
* columns
* cards
* projects (planned)
* users (planned)
* permissions (planned)
* API layer
* source of truth

---

## AI Service

Language:

* Python

Planned stack:

* FastAPI
* Pydantic
* MCP integration
* Agent orchestration layer

Responsibilities:

* task analysis
* subtask generation
* agent execution
* tool invocation
* AI automation

---

## Frontend

Planned:

* SvelteKit or Vue (TBD)
* TypeScript
* TailwindCSS
* Drag and Drop (kanban interactions)

Responsibilities:

* board UI
* task management
* delegation to AI
* project navigation
* real-time updates (planned)

---

## Database

PostgreSQL

Core entities:

* boards
* columns
* cards
* projects (planned)
* users (planned)
* project_members (planned)
* comments (planned)
* activity_logs (planned)

---

# Current Status

Implemented:

* Go backend
* PostgreSQL schema
* Boards CRUD
* Columns CRUD
* Cards CRUD
* Database migrations
* Hot reload with Air
* Basic API structure

In Progress:

* Frontend Kanban UI
* Board rendering
* Card creation
* Column rendering

Planned:

* Projects
* Users
* Authentication
* Permissions
* AI task delegation
* Agent execution
* MCP tools
* Real-time updates

---

# Development Setup

## Start infrastructure

PostgreSQL and Redis run via Docker:

```
make infra-up
```

## Run backend

```
cd apps/api
air
```

## Run frontend

```
cd apps/web
npm run dev
```

## Run AI service

```
cd apps/ai-service
uvicorn app.main:app --reload
```

---

# Roadmap

Phase 1 — Basic Kanban

* boards
* columns
* cards
* CRUD
* frontend UI

Phase 2 — Projects

* projects table
* project members
* board per project

Phase 3 — Auth

* users
* login
* sessions
* permissions

Phase 4 — AI Delegation

* delegate task to AI
* subtask generation
* status updates

Phase 5 — Agent Orchestration

* subagents
* MCP tools
* execution logs

---

# Vision

The long-term goal is to transform a traditional Kanban board into an **AI orchestration surface**, where tasks are not only tracked but actively executed by autonomous agents collaborating with users.

This project is both:

* a learning experiment
* a production-grade architecture exploration
* an AI workflow orchestration prototype
