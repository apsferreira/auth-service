# 📚 CONTEXT - Auth Service

## 🎯 PROPÓSITO

**Auth Service** é um microserviço centralizado de **autenticação** e **autorização** que fornece:

1. **Autenticação de usuários** (email + password)
2. **Geração e validação de JWT tokens**
3. **Sistema de permissões RBAC** (Role-Based Access Control)
4. **Multi-tenancy** (isolamento por organização)
5. **Sessões e auditoria**

---

## 🏗️ ARQUITETURA

### **Multi-Service Pattern**

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ MyLibrary   │────▶│ Auth        │◀────│ Outros      │
│ Frontend    │     │ Service     │     │ Services    │
└─────────────┘     └─────────────┘     └─────────────┘
       │                   │
       ▼                   ▼
┌─────────────┐     ┌─────────────┐
│ MyLibrary   │     │ Auth DB     │
│ PostgreSQL  │     │ PostgreSQL  │
└─────────────┘     └─────────────┘
```

### **Fluxo de Requisição:**

```
1. User faz login → POST /api/v1/auth/login
2. Auth Service valida credenciais
3. Gera JWT (access + refresh tokens)
4. Retorna tokens para o client

5. Client faz requisição → GET /api/v1/books (com JWT)
6. MyLibrary valida JWT no Auth Service
7. Auth Service retorna { valid: true, permissions: [...] }
8. MyLibrary processa requisição com tenant_id isolado
```

---

## 📊 DATABASE SCHEMA

### **Tabelas:**

**`tenants`** - Organizações/clientes
- `id`, `name`, `slug`, `plan`, `settings`

**`users`** - Usuários do sistema
- `id`, `tenant_id`, `email`, `password_hash`, `full_name`, `role_id`

**`roles`** - Níveis de acesso
- `id`, `tenant_id`, `name`, `level`, `is_system`

**`permissions`** - Permissões granulares
- `id`, `name`, `resource`, `action`

**`role_permissions`** - Mapeamento Role → Permissions
- `role_id`, `permission_id`

**`user_roles`** - Usuário pode ter múltiplas roles
- `user_id`, `role_id`

**`refresh_tokens`** - Tokens para refresh
- `id`, `user_id`, `token_hash`, `expires_at`

**`sessions`** - Auditoria de sessões
- `id`, `user_id`, `ip_address`, `user_agent`

---

## 🔐 JWT STRUCTURE

### **Access Token (15 min):**
```json
{
  "sub": "user-uuid",
  "tenant_id": "tenant-uuid",
  "email": "user@example.com",
  "roles": ["admin", "user"],
  "permissions": ["books.create", "books.delete"],
  "iat": 1704067200,
  "exp": 1704070800
}
```

### **Refresh Token (7 days):**
- Armazenado no banco (hash)
- Usado para obter novo access token
- Rotacionado a cada uso

---

## 🎭 RBAC - PERMISSÕES

### **Roles:**
- **Super Admin** (level 10) - Acesso total
- **Admin** (level 8) - CRUD livros, ver usuários
- **Manager** (level 6) - Gerenciar empréstimos
- **User** (level 5) - Criar livros, empréstimos
- **Viewer** (level 3) - Apenas leitura

### **Permissions:**
```
books.*         # Todas permissões de livros
books.create    # Criar livro
books.read      # Ver livros
books.update    # Atualizar
books.delete    # Deletar
books.export    # Exportar

loans.*         # Todas permissões de empréstimos
users.manage    # Gerenciar usuários
tenant.manage   # Gerenciar tenant
```

---

## 🎯 DECISÕES DE DESIGN

### **Por que microserviço separado?**
- ✅ Reuso em múltiplos sistemas
- ✅ Single Source of Truth para autenticação
- ✅ Escalabilidade independente
- ✅ Segurança centralizada

### **Por que JWT?**
- ✅ Stateless (não precisa armazenar sessões)
- ✅ Performático (validação rápida)
- ✅ Padrão da indústria

### **Por que RBAC?**
- ✅ Granular (permissões específicas)
- ✅ Hierárquico (levels de acesso)
- ✅ Flexível (múltiplas roles por user)

---

## 🧪 TESTING

### **Mock Data incluído:**
- 1 Tenant: `apsferreira`
- 2 Users: admin@apsferreira.com, user@apsferreira.com
- 5 Roles: Super Admin, Admin, Manager, User, Viewer
- 20+ Permissions: books.*, loans.*, users.*, tenant.*

### **Como testar:**
```bash
# Login
curl -X POST http://localhost:3002/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@apsferreira.com","password":"admin123"}'

# Validar token
curl -X POST http://localhost:3002/api/v1/auth/validate \
  -H "Authorization: Bearer <token>"
```

---

## 🚀 ROADMAP

### **Fase 1: MVP** ✅
- [x] Login/Refresh/Validate
- [x] RBAC básico
- [x] Multi-tenancy
- [x] Mock data

### **Fase 2: Integração** 🔜
- [ ] Middleware para MyLibrary
- [ ] Permission checks em endpoints
- [ ] Frontend integration

### **Fase 3: Avançado** 📋
- [ ] OAuth2 providers (Google, GitHub)
- [ ] 2FA
- [ ] Rate limiting
- [ ] Audit logs

---

**Última atualização:** 2025-01-22  
**Mantido por:** Antonio Ferreira

