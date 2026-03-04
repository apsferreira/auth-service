# PRD — auth-service (ENRIQUECIDO)
**Versão:** 1.0-enriched | **Data:** 2026-03-03 | **Status:** Implementado (~80%) — Sem testes | Sem remote GitHub
**Stack:** Go 1.23 + Fiber v2 | **Porta API:** :3010 | **Porta UI:** :3003
**Perspectivas:** @pm · @dados · @financeiro · @techlead · @qa

---

## Executive Summary + Veredito

O `auth-service` é o único ponto de autenticação e autorização de todo o ecossistema IIT. Já está implementado (~80%), mas possui zero cobertura de testes e não tem repositório remoto no GitHub — o que significa que existe apenas localmente, sem backup, sem CI/CD e sem segurança operacional alguma.

**Veredito: ESTABILIZAR URGENTE — serviço crítico de segurança operando sem rede de proteção.**

Situação atual:
- ✅ Funcionalidades core implementadas (OTP email, Telegram, JWT, RBAC, multi-tenant)
- ❌ 0% de cobertura de testes (zero arquivos `*_test.go`)
- ❌ Sem repositório remoto (risco de perda total do código)
- ❌ Sem CI/CD
- ❌ Sem OpenAPI spec para outros serviços consultarem

A prioridade imediata não é adicionar features, mas **proteger o que já existe** com testes, remote GitHub e pipeline de CI.

---

## Papel no ecossistema IIT

O auth-service é a fundação de segurança de todos os demais serviços. Nenhum produto vertical consegue operar sem ele.

```
┌─────────────────────────────────────────────────────────────┐
│                     auth-service :3010                      │
│                                                              │
│  OTP (email/TG/WA) → JWT (access + refresh) → RBAC         │
│  Multi-tenant → Admin panel → Provision-user M2M            │
└──────────────────────────┬──────────────────────────────────┘
                           │ valida JWT (GET /auth/validate)
         ┌─────────────────┼─────────────────────────────┐
         │                 │                             │
   customer-service   catalog-service               todos os
   scheduling-service  cart-service                 produtos
   checkout-service    events-service               verticais
```

**Responsabilidades:**
- Autenticação OTP multi-canal (email via Resend, Telegram, WhatsApp)
- Login por senha para admin panels
- Emissão e validação de JWT (access + refresh tokens)
- RBAC: roles (student, instructor, admin, owner) + permissões granulares
- Multi-tenant: isolamento completo por `tenant_id`
- Provisionamento de usuários M2M (customer-service → auth-service)
- Auditoria de acessos

---

## Quem usa e como (dependências)

### Usuários diretos

**Administradores de tenants:**
- Login por senha no admin panel (React UI :3003)
- Criam e gerenciam usuários, roles e tenants

**Clientes/alunos dos produtos verticais:**
- Autenticação OTP (email ou Telegram) sem senha
- Fluxo: solicitar OTP → receber código → verificar → JWT emitido

**Serviços do ecossistema (M2M):**
- Todos os serviços validam JWTs via `GET /auth/validate`
- customer-service provisiona usuários via `POST /auth/provision-user`

### Mapeamento de roles por produto

| Produto | Roles |
|---------|-------|
| jiu-jitsu-academy | owner, admin, instructor, student |
| food-marketplace | owner, admin, customer |
| restaurant-qr | owner, admin, waiter (garçom anônimo sem auth) |
| my-library | owner (produto pessoal) |
| focus-hub | owner (produto pessoal) |

### Dependências técnicas

| Componente | Uso |
|-----------|-----|
| PostgreSQL auth_db | Usuários, tenants, roles, refresh tokens |
| Redis DB2 | OTPs (TTL 5min), rate limiting, sessões |
| Resend | Envio de OTP por email |
| Evolution API | OTP via WhatsApp |
| Telegram Bot API | OTP via Telegram |

---

## Funcionalidades MVP (P0/P1/P2)

### P0 — Já implementado (verificar + cobrir com testes)

| ID | Funcionalidade | Status |
|----|---------------|--------|
| A-01 | Login email+senha (admin) | ✅ Implementado |
| A-02 | OTP por email (Resend) | ✅ Implementado |
| A-03 | OTP por Telegram | ✅ Implementado |
| A-04 | JWT access + refresh token | ✅ Implementado |
| A-05 | Validar JWT (M2M) | ✅ Implementado |
| A-06 | Refresh token | ✅ Implementado |
| A-07 | RBAC roles e permissões | ✅ Implementado |
| A-08 | Multi-tenant por tenant_id | ✅ Implementado |
| A-09 | Provisionar usuário (M2M) | ✅ Implementado |
| A-10 | Admin panel React :3003 | ✅ Implementado |

### P1 — Dívida técnica crítica (fazer antes de qualquer nova feature)

| ID | Tarefa | Descrição |
|----|--------|-----------|
| A-11 | Criar remote GitHub | `apsferreira/auth-service` + push do código atual |
| A-12 | Testes handlers | TestLoginHandler, TestOTPRequestHandler, TestOTPVerifyHandler |
| A-13 | Testes JWT service | geração, validação, expiração, claims, adulteração |
| A-14 | Testes RBAC service | herança de roles, check de permissões |
| A-15 | Testes refresh + logout | token válido, expirado, revogado |
| A-16 | Testes validate M2M | retorno correto de claims, token adulterado |
| A-17 | CI/CD GitHub Actions | go test + go vet + go build + coverage ≥80% |
| A-18 | OpenAPI/Swagger spec | spec para todos os endpoints (swaggo/swag) |

### P2 — Roadmap futuro

| ID | Funcionalidade | Descrição |
|----|---------------|-----------|
| A-19 | OTP via WhatsApp | Canal adicional via Evolution API |
| A-20 | Revogação de tokens em massa | Quando tenant é desativado |
| A-21 | Log de auditoria | Tabela de eventos de acesso por usuário/tenant |
| A-22 | SSO / OAuth2 | Login com Google para produtos B2C |
| A-23 | MFA opcional | TOTP (Google Authenticator) para admins |

---

## Schema principal

```sql
-- auth_db

CREATE TABLE tenants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(200) NOT NULL,
  slug VARCHAR(100) UNIQUE NOT NULL,
  active BOOLEAN DEFAULT true,
  config JSONB,                   -- timezone, logo, custom settings
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id),
  email VARCHAR(254) UNIQUE NOT NULL,
  full_name VARCHAR(200),
  phone VARCHAR(20),
  telegram_id BIGINT,
  password_hash TEXT,             -- apenas para admins (login por senha)
  active BOOLEAN DEFAULT true,
  last_login_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE roles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id),
  name VARCHAR(50) NOT NULL,
  permissions JSONB NOT NULL      -- ["schedule:read", "billing:write", ...]
);

CREATE TABLE user_roles (
  user_id UUID REFERENCES users(id),
  role_id UUID REFERENCES roles(id),
  PRIMARY KEY (user_id, role_id)
);

CREATE TABLE refresh_tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  token_hash TEXT UNIQUE NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Endpoints

```
POST /api/v1/auth/request-otp     → solicita OTP (email ou Telegram)
POST /api/v1/auth/verify-otp      → verifica OTP → emite JWT
POST /api/v1/auth/admin-login     → login email+senha → JWT
POST /api/v1/auth/refresh         → renova access token via refresh
POST /api/v1/auth/logout          → revoga refresh token
GET  /api/v1/auth/validate        → valida JWT (M2M) → retorna claims
POST /api/v1/auth/provision-user  → M2M: cria usuário com role (customer-service)

-- Admin panel
GET/POST/PUT/DELETE /api/v1/admin/tenants
GET/POST/PUT/DELETE /api/v1/admin/users
GET/POST/PUT/DELETE /api/v1/admin/roles

GET /health
```

---

## Unit economics / impacto financeiro

### Custo de desenvolvimento (dívida técnica restante)

| Atividade | Horas | Custo @R$200/h |
|-----------|-------|----------------|
| Criar remote GitHub + configurar CI | 4h | R$800 |
| Suite de testes handlers (6 handlers) | 24h | R$4.800 |
| Testes JWT service + RBAC service | 16h | R$3.200 |
| Testes integração PostgreSQL (testcontainers) | 12h | R$2.400 |
| E2E: fluxo OTP completo | 8h | R$1.600 |
| OpenAPI spec (swaggo/swag) | 8h | R$1.600 |
| **Total dívida técnica** | **72h** | **R$14.400** |

### Custo de um incidente de segurança sem testes
- Vulnerabilidade em JWT (bypass de autenticação): acesso indevido a todos os produtos → dano reputacional + legal estimado em R$50K-500K
- Bug em RBAC: admin de tenant A acessa dados do tenant B → vazamento de dados de clientes
- **ROI de testes: qualquer incidente custa 3x-30x mais do que os R$14.400 de testes**

### Custo operacional

| Serviço | Custo mensal |
|---------|-------------|
| Resend (email OTP) | ~R$0 (free tier 3.000 emails/mês) → ~R$75/mês em escala |
| Telegram Bot | Gratuito |
| WhatsApp OTP (futuro) | Custo Evolution API (infra própria) |
| Redis DB2 | Compartilhado (shared-redis) |

---

## KPIs e North Star Metric

**North Star Metric:** `auth_success_rate` — % de tentativas de autenticação que resultam em JWT válido emitido

| KPI | Meta M1 | Meta M3 |
|-----|---------|---------|
| auth_success_rate (OTP) | > 95% | > 98% |
| OTP delivery latency (email) | < 5s P95 | < 3s P95 |
| JWT validation latency | < 10ms P99 | < 5ms P99 |
| Cobertura de testes | 80% | 90% |
| auth_db uptime | > 99.9% | > 99.9% |
| Falsos positivos de rate limit | < 0.1% | < 0.05% |
| OTP brute force bloqueados | > 99% | > 99.9% |

**Métricas de segurança:**
- `otp_failed_attempts_rate` — picos indicam tentativas de força bruta
- `jwt_validation_errors_rate` — tokens adulterados ou expirados
- `invalid_tenant_access_attempts` — possíveis ataques cross-tenant

---

## Riscos top 3

### Risco 1: Perda de código (ALTA PROBABILIDADE, CRÍTICA)
**Descrição:** O repositório existe apenas localmente. Qualquer falha no disco do macOS de desenvolvimento resulta na perda total do código implementado.
**Mitigação:** Criar remote GitHub `apsferreira/auth-service` e fazer push IMEDIATAMENTE. Esta é a ação de maior urgência do ecossistema inteiro.
**Severidade:** Crítica — perda irreversível de semanas de desenvolvimento

### Risco 2: Vulnerabilidade em produção sem testes (ALTA PROBABILIDADE)
**Descrição:** Com 0% de cobertura de testes, qualquer refactor ou adição de feature pode introduzir vulnerabilidades de segurança (bypass de JWT, privilege escalation, tenant leakage). Em um serviço de autenticação, isso é catastrófico.
**Mitigação:** Implementar suite de testes antes de qualquer deploy em produção. Nenhum serviço dependente deve usar o auth-service em produção sem 80% de cobertura.
**Severidade:** Alta — risco de segurança sistêmico

### Risco 3: Acoplamento M2M sem contratos formais (MÉDIA PROBABILIDADE)
**Descrição:** customer-service usa `POST /auth/provision-user` sem spec formal. Qualquer mudança de contrato quebra o fluxo de criação de clientes silenciosamente.
**Mitigação:** Gerar OpenAPI spec com swaggo/swag e versionar a API. Contratos de API formais entre serviços M2M.
**Severidade:** Média — causa bugs silenciosos difíceis de debugar

---

## Roadmap simplificado (3 fases)

### Fase 1 — Estabilização (URGENTE — Sprint atual, ~2 semanas)
- Criar repositório remoto `apsferreira/auth-service` + push imediato
- Suite completa de testes unitários (handlers + services + repository)
- GitHub Actions: go test + go vet + cobertura ≥80%
- OpenAPI spec (swaggo/swag) para todos os endpoints
- Testcontainers para testes de integração com PostgreSQL

### Fase 2 — Expansão de canais (Sprint 3-4)
- OTP via WhatsApp (Evolution API) — terceiro canal
- Auditoria de acessos (tabela de eventos)
- ADRs documentados: JWT vs Session, OTP vs senha, multi-tenant
- Rate limiting robusto por IP + por tenant

### Fase 3 — Enterprise features (Sprint 6+)
- SSO / OAuth2 (login com Google)
- MFA opcional (TOTP) para administradores
- Revogação em massa de tokens por tenant
- Self-service de recuperação de conta
- Logs de auditoria exportáveis (compliance)

---

## Checklist de estabilização (agir agora)

- [ ] **HOJE:** `git remote add origin https://github.com/apsferreira/auth-service && git push`
- [ ] Criar GitHub Actions workflow com go test + coverage
- [ ] Escrever TestLoginHandler (credencial válida, senha errada, tenant inválido)
- [ ] Escrever TestOTPRequestHandler + TestOTPVerifyHandler
- [ ] Escrever TestJWTService (geração, validação, expiração)
- [ ] Escrever TestRBACService (roles, permissões, herança)
- [ ] Escrever TestValidateTokenHandler (M2M)
- [ ] Atingir ≥80% de cobertura antes de qualquer deploy em produção
- [ ] Gerar OpenAPI spec com swaggo/swag
