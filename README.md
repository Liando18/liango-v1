# LianGo 🚀

> **Personal Backend Starter Kit** — Built with Gin + GORM for production-ready Go APIs.

LianGo is a reusable, modular backend boilerplate that speeds up development across multiple projects while maintaining clean architecture, security, and scalability.

---

## Tech Stack

| Layer      | Technology                          |
| ---------- | ----------------------------------- |
| Language   | Go 1.21+                            |
| Framework  | Gin                                 |
| ORM        | GORM                                |
| Database   | PostgreSQL (default) / MySQL        |
| Auth       | JWT (Access + Refresh Token)        |
| Security   | bcrypt, rate limiter, CORS, headers |
| Hot Reload | Air                                 |

---

## Project Structure

```
liango/
├── main.go                         # Entry point
├── .env.example                    # Environment template
├── .air.toml                       # Hot reload config
│
├── config/
│   └── app.go                      # App config & .env loader
│
├── database/
│   ├── connection.go               # Multi-driver DB connection (postgres/mysql)
│   └── migrator.go                 # AutoMigrate helper
│
├── app/
│   ├── models/
│   │   ├── base_model.go           # BaseModel: UUID + timestamps + soft delete
│   │   ├── user.go                 # User + Token models
│   │   └── ExampleModel.go         # ← TEMPLATE: copy for new entities
│   │
│   ├── controllers/
│   │   ├── AuthController.go       # Register, Login, Refresh, Logout, Me
│   │   └── ExampleController.go    # ← TEMPLATE: full CRUD controller
│   │
│   ├── services/
│   │   ├── AuthService.go          # Auth business logic
│   │   └── ExampleService.go       # ← TEMPLATE: CRUD business logic
│   │
│   ├── repositories/
│   │   ├── UserRepository.go       # User DB queries
│   │   └── ExampleRepository.go    # ← TEMPLATE: CRUD DB queries
│   │
│   ├── routes/
│   │   ├── api.go                  # Master route registry
│   │   └── ExampleRoute.go         # ← TEMPLATE: resource routes
│   │
│   ├── middlewares/
│   │   ├── jwt.go                  # JWT + Role + Permission middleware
│   │   ├── cors.go                 # CORS middleware
│   │   └── security.go             # Rate limiter, API key, IP whitelist, timeout
│   │
│   ├── validations/
│   │   └── ExampleValidation.go    # ← TEMPLATE: request validation structs
│   │
│   ├── responses/
│   │   └── response.go             # Standard API response format
│   │
│   ├── helpers/
│   │   ├── jwt.go                  # JWT generate & parse helpers
│   │   ├── password.go             # bcrypt helpers
│   │   └── pagination.go           # Pagination helper
│   │
│   ├── utils/
│   │   └── logger.go               # Structured logger
│   │
│   └── constants/
│       └── roles.go                # Role & permission constants
│
├── cmd/
│   └── generator/
│       └── main.go                 # CLI code generator
│
└── storage/
    └── logs/                       # Log files
```

---

## Quick Start

### 1. Clone & Setup

```bash
git clone https://github.com/Liando18/liango-v1.git my-project
cd my-project
cp .env.example .env
# Edit .env with your config
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Run with Hot Reload

```bash
# Install air (first time)
go install github.com/cosmtrek/air@latest

# Run
air
```

### 4. Or run directly

```bash
go run main.go
```

---

## API Endpoints

### Welcome

```
GET /
→ { "message": "Welcome to LianGo version 1.0.0" }
```

### Auth

```
POST /api/v1/auth/register   Register new user
POST /api/v1/auth/login      Login (returns access + refresh token)
POST /api/v1/auth/refresh    Refresh access token
POST /api/v1/auth/logout     Revoke refresh token
GET  /api/v1/auth/me         Get current user (requires JWT)
```

### Example Resource (requires JWT)

```
GET    /api/v1/examples         List all (paginated)
GET    /api/v1/examples/:id     Get one
POST   /api/v1/examples         Create (requires create:any permission)
PUT    /api/v1/examples/:id     Update (requires update:any permission)
DELETE /api/v1/examples/:id     Delete (admin only)
```

---

## Database Configuration

In `.env`:

```env
# PostgreSQL
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=liango_db

# MySQL (switch by changing DB_DRIVER)
# DB_DRIVER=mysql
# DB_PORT=3306
```

No other code changes needed when switching databases.

---

## CLI Generator

Generate boilerplate files automatically:

```bash
# Single file
go run cmd/generator/main.go make:model Product
go run cmd/generator/main.go make:controller Product
go run cmd/generator/main.go make:service Product
go run cmd/generator/main.go make:repository Product

# Generate all CRUD files at once
go run cmd/generator/main.go make:crud Product
```

---

## Adding a New Resource

1. **Generate files:**

   ```bash
   go run cmd/generator/main.go make:crud Product
   ```

2. **Register the route** in `app/routes/api.go`:

   ```go
   RegisterProductRoutes(v1)
   ```

3. **Add model to migration** in `main.go`:

   ```go
   database.Migrate(&models.Product{})
   ```

4. Done! ✅

---

## Security Features

| Feature          | Implementation                          |
| ---------------- | --------------------------------------- |
| JWT Auth         | Access token (15m) + Refresh token (7d) |
| Token Rotation   | Old refresh token revoked on refresh    |
| Password Hashing | bcrypt (default cost)                   |
| Rate Limiting    | 100 req/min per IP                      |
| CORS             | Configurable via env                    |
| Security Headers | X-Frame-Options, CSP, HSTS, etc.        |
| IP Whitelist     | Optional via ALLOWED_IPS env            |
| API Key          | Optional via API_KEY env                |
| Request Timeout  | 30 seconds global timeout               |
| Panic Recovery   | Auto-recover with 500 response          |
| Soft Delete      | All models use GORM soft delete         |

---

## Roles & Permissions

Defined in `app/constants/roles.go`:

| Role   | Permissions                                  |
| ------ | -------------------------------------------- |
| admin  | create:any, read:any, update:any, delete:any |
| editor | create:any, read:any, update:own             |
| user   | read:own, update:own                         |

---

## Standard API Response Format

All responses follow this structure:

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "meta": {
    "page": 1,
    "per_page": 15,
    "total": 100,
    "total_pages": 7
  }
}
```

---

Developer by Aprilian Gevindo
