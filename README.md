# Todo Golang API

This is a simple **Todo APP** built with Go, Chi and GORM in a bid to learn Golang while following clean architecture principles.  
This project demonstrates modular design with authentication, event-driven patterns, and real-time updates via Server-Sent Events (SSE).

---

## Features

- User authentication (JWT-based)
- Todo management (CRUD operations)
- Custom Event Bus Implementation
- Repository Pattern for data abstraction
- Server-Sent Events (SSE) for real-time updates
- Token blacklist support for logout/invalidation
- Config-driven setup with `.env`
- Unit and integration testing support

---

## Getting Started

### Clone the repository

```bash
git clone https://github.com/horlerdipo/todo-golang.git
cd todo-golang
go mod tidy
cp .env.example .env
```

### Run the application

You can either run this application by 
```bash
go run cmd/api/main.go
```

or using air
```bash
air
```

## Architectural Decisions
### Custom Event Bus
    
I implemented a custom and minimal Event Bus in ```internal/events```.

Enables loose coupling between services (e.g., TodoService can publish events without knowing who listens).

Simplifies real-time event propagation (for SSE).

### Repository Pattern
Repositories define interfaces for persistence operations.
Actual implementations (using GORM) live in ```internal/repos```
This allows:

Easier testing with in-memory repos
Future-proofing (swap DB implementation without touching the services)
Separation of concerns (The services are not aware of the database, just the repository)

### Server Sent Events
I implemented SSE to enable real-time streaming of updates.
When Todos are created, updated, or deleted, events are pushed to connected clients and this was integrated with the event bus for seamless data flow.
I implemented this as a lightweight alternative to Websockets.

## Possible Improvements
- Expand unit test coverage across services

- Add Dockerfile and Docker Compose setup

- CI/CD pipeline integration with GitHub Actions

- Implement Waterfall for proper event driven architecture and persistence

- Update ```golangci-lint.yml```

## Testing

Run all test 
```bash
go test ./...
```