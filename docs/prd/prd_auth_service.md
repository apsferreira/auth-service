# PRD — auth-service

**Versão:** 1.0  
**Data:** 2026-02-26  
**Status:** Implementado (~80%) — Sem testes  
**Stack:** Go + Fiber  
**Porta:** :3010

---

## 1. Visão

Serviço centralizado de autenticação e autorização para todo o ecossistema. Fornece JWT, OTP multi-canal, RBAC e multi-tenant para todos os produtos verticais.

---

## 2. Funcionalidades Implementadas

| Funcionalidade | Status |
|---------------|--------|
| Login com credenciais (email/senha) | ✅ |
| OTP por email (Resend) | ✅ |
| OTP por Telegram | ✅ |
| JWT access + refresh token | ✅ |
| RBAC: roles e permissões | ✅ |
| Multi-tenant por tenant_id | ✅ |
| Admin login (migrations 011-012) | ✅ |
| Validação de token (M2M) | ✅ |

---

## 3. Endpoints

| Método | Rota | Descrição |
|--------|------|-----------|
| POST | /auth/login | Login email+senha → JWT |
| POST | /auth/otp/request | Solicita OTP (email ou Telegram) |
| POST | /auth/otp/verify | Verifica OTP → JWT |
| POST | /auth/refresh | Renova access token via refresh |
| POST | /auth/logout | Revoga refresh token |
| GET | /auth/validate | Valida JWT (M2M) → retorna claims |
| POST | /auth/provision-user | Provisiona usuário para customer-service |
| GET | /health | Health check |

---

## 4. Schema JWT Claims

```json
{
  "sub": "user_uuid",
  "tenant_id": "academia-fitlife",
  "email": "joao@email.com",
  "roles": ["student", "admin"],
  "exp": 1234567890
}
```

---

## 5. Pendências

| Item | Prioridade |
|------|------------|
| Suite de testes ≥80% | 🔴 Crítico |
| Repositório remoto GitHub | 🔴 Crítico |
| OpenAPI/Swagger spec | 🟠 Alta |
| OTP via WhatsApp | 🟡 Média |
