# HRD Room - Skill Documentation

## Overview
This skill provides instructions for developing, maintaining, and extending the HRD Room platform.

## Common Development Tasks

### 1. Backend Development (Go)

#### Adding a New Feature
1.  **Define Domain Entity**: Add the struct and interfaces in `backend/internal/domain`.
2.  **Implement Repository**: Add database methods in `backend/internal/repository`.
3.  **Implement Usecase**: Implement business logic in `backend/internal/usecase`.
4.  **Implement Handler**: Create the HTTP handler in `backend/internal/delivery/http`.
5.  **Register Routes**: Wire up the dependency injection and register routes in `backend/cmd/api/main.go` (or similar).

#### Database Migrations
Migrations are stored in `backend/migrations`.
- To create a new migration: Create `.up.sql` and `.down.sql` files with a timestamp/version prefix.

### 2. Frontend Development (Vue)

#### State Management
Use Pinia stores in `frontend/src/stores`. Maintain a clear separation between API calls (in `src/api`) and state logic.

#### Styling & UI
- Use **TailwindCSS** for layout and spacing.
- Use **Shadcn/UI** components for consistent design.
- Components are located in `frontend/src/components`.

### 3. Running the Project
The project is designed to run with Docker Compose.

```bash
# Start all services
docker-compose up --build

# Backend only logs
docker-compose logs -f backend
```

---

## Best Practices

### Backend
- **Error Handling**: Use the Echo framework's error handling. Ensure all repository methods return meaningful errors.
- **Dependency Injection**: Manually wire dependencies in the main entry point to keep layers decoupled.
- **Validation**: Use a validation library (like `go-playground/validator`) in the delivery layer.

### Frontend
- **Type Safety**: Use TypeScript for all new components and stores.
- **Components**: Break down large pages into small, reusable components.
- **Async Handling**: Use `loading` states in Pinia stores when fetching data.

### Commit Guidelines
Follow conventional commits (`feat:`, `fix:`, `chore:`, etc.).

---

## Troubleshooting

### Database Connection Issues
If the backend cannot connect to the database, ensure:
1.  The `db` container is healthy.
2.  The `.env` file matches the `docker-compose.yml` environment variables.
3.  The port `5433` (host) or `5432` (internal) is not blocked.

### Frontend Build Errors
If `npm run build` fails, check:
1.  TypeScript errors in `src`.
2.  Vite configuration and environment variables.
3.  Missing dependencies in `package.json`.
