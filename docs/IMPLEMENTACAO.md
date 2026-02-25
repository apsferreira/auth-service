# 🛠️ GUIA DE IMPLEMENTAÇÃO - Auth Service

## 📋 ESTRUTURA DE ARQUIVOS

```
backend/
├── cmd/api/main.go                 # Entry point
├── internal/
│   ├── domain/
│   │   ├── user.go                 # User entity
│   │   ├── tenant.go               # Tenant entity
│   │   ├── role.go                 # Role entity
│   │   └── permission.go           # Permission entity
│   ├── repository/
│   │   ├── user_repository.go      # Interface
│   │   ├── tenant_repository.go    # Interface
│   │   └── postgres/               # Implementations
│   ├── service/
│   │   ├── auth_service.go         # Login, logout, refresh
│   │   ├── jwt_service.go          # JWT generation/validation
│   │   └── user_service.go         # User CRUD
│   ├── handler/
│   │   ├── auth_handler.go         # /api/v1/auth/*
│   │   ├── user_handler.go         # /api/v1/users/*
│   │   └── health_handler.go       # /health
│   └── middleware/
│       ├── cors.go                 # CORS config
│       └── logger.go               # Logging
├── pkg/
│   ├── config/config.go            # Load env vars
│   ├── database/postgres.go        # DB connection
│ environment variables
│   └── jwt/jwt.go                  # JWT utilities
└── migrations/                     # SQL migrations
```

---

## 🔐 FLUXO DE LOGIN

### **1. Receber credenciais (email + password)**

### **2. Buscar user no banco (com bcrypt.CompareHashAndPassword)**

### **3. Buscar roles e permissions do user**

### **4. Gerar JWT (access + refresh)**

### **5. Salvar refresh token no banco**

### **6. Retornar tokens para o client**

---

## 🎯 IMPLEMENTAÇÃO PASSO A PASSO

### **Step 1: Domain Models**
Criar structs em `internal/domain/`:
- `User`, `Tenant`, `Role`, `Permission`

### **Step 2: Repository Interfaces**
Criar interfaces em `internal/repository/`:
- `UserRepository`, `TenantRepository`

### **Step 3: PostgreSQL Implementation**
Implementar queries em `internal/repository/postgres/`

### **Step 4: JWT Service**
Criar `internal/service/jwt_service.go` com:
- `GenerateAccessToken(user)`
- `GenerateRefreshToken(user)`
- `ValidateToken(token)`

### **Step 5: Auth Service**
Criar `internal/service/auth_service.go` com:
- `Login(email, password)`
- `Refresh(token)`
- `Validate(token)`
- `Logout(token)`

### **Step 6: Handlers**
Criar HTTP handlers em `internal/handler/`:
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/validate`

### **Step 7: Main**
Configurar Fiber em `cmd/api/main.go`

---

## 📝 CHECKLIST DE IMPLEMENTAÇÃO

- [ ] Domain models criados
- [ ] Repository interfaces definidas
- [ ] PostgreSQL implementations
- [ ] JWT service (generate + validate)
- [ ] Auth service (login, refresh, validate)
- [ ] HTTP handlers
- [ ] Middleware (CORS, logger)
- [ ] Main.go configurado
- [ ] Migrations testadas
- [ ] Seed data funcionando
- [ ] Health endpoint
- [ ] Testes unitários

---

## 🧪 TESTING

```bash
# 1. Start services
make up

# 2. Run migrations
make migrate-up

# 3. Seed data
make seed

# 4. Test login
curl -X POST http://localhost:3002/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@apsferreira.com","password":"admin123"}'

# 5. Test validate (use token from step 4)
curl -X POST http://localhost:3002/api/v1/auth/validate \
  -H "Authorization: Bearer <token>"
```

---

**Próximo:** Implementação completa... 🚀

