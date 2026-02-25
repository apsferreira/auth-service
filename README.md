# 🔐 Auth Service - Central Authentication & Authorization

## 📋 VISÃO GERAL

Serviço centralizado de **autenticação** e **autorização** para o ecossistema UCG (Unified Cloud Governance).

**Stack:**
- Go 1.21+ + Fiber v2
- PostgreSQL 15+
- JWT (access + refresh tokens)
- RBAC (Role-Based Access Control)

---

## 🚀 INÍCIO RÁPIDO

### **Pré-requisitos:**
- Docker & Docker Compose
- Go 1.21+ (opcional, para desenvolvimento local)

### **Executar:**

```bash
# 1. Clone o repositório
git clone <url>
cd auth-service

# 2. Configure variáveis de ambiente
cp .env.example .env
# Edite .env conforme necessário

# 3. Suba os containers
make up

# 4. Execute migrations
make migrate-up

# 5. Seed data (usuários mock)
make seed

# 6. Teste o endpoint
curl http://localhost:3002/health
```

---

## 📊 MOCK DATA

O Auth Service vem com dados mockados para testes:

**Tenants:**
- `apsferreira` (plan: premium)

**Users:**
```json
{
  "email": "admin@apsferreira.com",
  "password": "admin123",
  "role": "admin",
  "tenant": "apsferreira"
}
```

```json
{
  "email": "user@apsferreira.com",
  "password": "user123",
  "role": "user",
  "tenant": "apsferreira"
}
```

---

## 🔑 ENDPOINTS

### **Auth:**
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/validate` - Validar JWT
- `POST /api/v1/auth/logout` - Logout

### **Users:**
- `GET /api/v1/users` - Listar usuários
- `POST /api/v1/users` - Criar usuário
- `GET /api/v1/users/:id` - Buscar usuário
- `PUT /api/v1/users/:id` - Atualizar usuário
- `DELETE /api/v1/users/:id` - Deletar usuário

### **Health:**
- `GET /health` - Health check

---

## 🔗 INTEGRAÇÃO COM MYLIBRARY

```bash
# No MyLibrary backend, adicione no .env:
AUTH_SERVICE_URL=http://auth-service:3002
JWT_SECRET=your-secret-key

# No MyLibrary frontend:
AUTH_API_URL=http://localhost:3002
```

---

## 📝 DOCUMENTAÇÃO

- [Context](./docs/CONTEXT.md) - Arquitetura e design
- [API Reference](./docs/API.md) - Detalhes dos endpoints
- [Deployment](./docs/DEPLOYMENT.md) - Guia de deploy

---

## 🛠️ DESENVOLVIMENTO

```bash
# Rodar testes
make test

# Ver logs
make logs

# Parar containers
make down

# Reset completo (CAUTION!)
make clean
```

---

**Powered by UCG** 🚀

