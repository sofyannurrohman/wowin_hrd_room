# HRD Room - Architecture Documentation

## Overview
HRD Room is a comprehensive assessment platform built to streamline recruitment. It features an HR Dashboard for management and a Participant Portal for exams.

## Technology Stack

### Backend
- **Language**: Go (Golang)
- **Framework/Libraries**:
  - `github.com/labstack/echo/v4`: Web framework.
  - `github.com/jmoiron/sqlx`: Database toolkit.
  - `github.com/golang-migrate/migrate`: Migration tool.
- **Database**: PostgreSQL
- **Architecture**: Clean Architecture / Ports and Adapters.

### Frontend
- **Framework**: Vue 3 (Composition API)
- **Build Tool**: Vite
- **State Management**: Pinia
- **Styling**: TailwindCSS & Shadcn/UI
- **Language**: TypeScript

### DevOps
- **Containerization**: Docker & Docker Compose
- **Reverse Proxy**: Nginx

---

## Backend Architecture

The backend follows Clean Architecture principles, ensuring separation of concerns and testability.

### Layers:
1.  **Domain (`internal/domain`)**: Contains core entities (structs) and repository/usecase interfaces. No dependencies on other layers.
2.  **Usecase (`internal/usecase`)**: Business logic implementation. Depends only on the domain layer.
3.  **Repository (`internal/repository`)**: Implementation of data access interfaces (PostgreSQL). Handles SQL queries and mapping.
4.  **Delivery (`internal/delivery/http`)**: HTTP handlers and routing. Interacts with usecases to handle requests and return responses.
5.  **Middleware (`internal/middleware`)**: Cross-cutting concerns like Auth, CORS, and Logging.

---

## Frontend Architecture

The frontend is structured as a modern Vue 3 Single Page Application (SPA).

### Key Directories:
- `src/api`: Axios instances and services for communicating with the backend.
- `src/stores`: Pinia stores for application state (Auth, Sessions, etc.).
- `src/pages`: Component-based views mapped to routes.
- `src/components`:
  - `ui/`: Shadcn components.
  - `shared/`: Reusable application-specific components.
- `src/composables`: Logic extracted from components for reuse.

---

## Core Features Flow

### 1. Recruitment Setup
HR creates **Job Positions** and adds **Questions** to the bank. Questions are then grouped into **Modules**.

### 2. Examination Session
HR creates a **Session** by selecting a Job Position, Modules, and a Schedule. A unique **Token** is generated.

### 3. Participant Journey
Applicants sign up, then use the Token to join a session. The system monitors the session for **violations** (e.g., tab switching).

### 4. Grading & Analytics
Objective questions are auto-graded. Subjective questions are manually graded by HR. Results are aggregated, and analytics are provided via the dashboard.

---

## Deployment
The project uses `docker-compose` to orchestrate:
- `backend`: Go API.
- `frontend`: Vite app served via Nginx (in production) or dev server.
- `db`: PostgreSQL instance.
- `nginx`: Reverse proxy routing `/api` to the backend and other routes to the frontend.
