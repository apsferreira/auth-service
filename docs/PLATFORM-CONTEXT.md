# auth-service — Contexto de Plataforma IIT

> Documento complementar ao CONTEXT.md existente. O CONTEXT.md documenta a arquitetura técnica interna; este arquivo posiciona o auth-service no ecossistema IIT conforme o IIT-PLATFORM-PRD.md (Março/2026).

---

## Papel Estratégico no Ecossistema

O auth-service é a **fundação crítica** (P0) de toda a plataforma IIT. Nenhum produto pode existir sem ele.

**Princípio:** Um usuário não deve ter 5 senhas diferentes para 5 produtos IIT — deve ter uma identidade única. O auth-service implementa esse SSO centralizado.

**Produtos que dependem diretamente:**

| Produto | Status de Integração |
|---------|---------------------|
| my-library | Integrado (RBAC + JWT) |
| focus-hub | Pendente (integração bloqueante) |
| iit-agents | A definir |
| attend-agent | Planejado — obrigatório antes do lançamento |
| jiu-jitsu-academy | Planejado — obrigatório antes do lançamento |
| restaurant-qr | Planejado |

---

## Evolução Planejada (IIT-PLATFORM-PRD)

### v1.0 (atual — em produção)
- OTP por email (Resend)
- JWT + refresh token rotacionado
- Rate limiting (5 OTPs/hora/email)
- Painel admin React (local only)
- App Registry (client_id + client_secret por produto)

### v1.5 — Q2 2026
- **OTP por SMS** (Twilio) — handoff ao notification-service
- **OTP por WhatsApp** (Meta Cloud API) — handoff ao notification-service
- **Google OAuth** — reduz fricção de cadastro
- Nota: quando notification-service estiver disponível, o auth-service **para de enviar email direto** e publica evento `auth.otp.requested` que o notification-service consome

### v2.0 — Q3 2026
- **SAML/SSO enterprise** — para clientes corporativos do attend-agent
- **TOTP 2FA** — para admins IIT e tenants de alto risco
- **Audit log LGPD** — exportação de histórico de acesso por usuário (direito de acesso LGPD Art. 18)

---

## Integração com notification-service (Pendente)

O auth-service atualmente envia email OTP diretamente via Resend. Quando o **notification-service** for lançado (Q1 2026), o fluxo muda:

**Fluxo atual:**
```
auth-service → Resend API → email do usuário
```

**Fluxo futuro (eventos):**
```
auth-service → publica auth.otp.requested (RabbitMQ)
notification-service ← consome evento
notification-service → Resend (email) ou Twilio (SMS) ou Meta (WhatsApp)
```

**Evento a publicar:**
```json
{
  "type": "auth.otp.requested",
  "tenant_id": "iit-platform",
  "payload": {
    "user_email": "usuario@exemplo.com",
    "otp_code": "123456",
    "expires_in_minutes": 10,
    "channel": "email",
    "app_name": "FocusHub"
  }
}
```

Outros eventos que o auth-service deve publicar via events-service SDK:
```
auth.user.login_success
auth.user.login_failed
auth.otp.requested      → notification-service processa
auth.user.logout
```

---

## App Registry — Como Cada Produto se Registra

Todo produto IIT deve se registrar no auth-service para receber um `client_id` e `client_secret`:

```http
POST /v1/admin/apps
{
  "name": "FocusHub",
  "redirect_uris": ["https://focus-hub.institutoitinerante.com.br/auth/callback"],
  "allowed_origins": ["https://focus-hub.institutoitinerante.com.br"]
}
```

Resposta:
```json
{
  "client_id": "focushub-abc123",
  "client_secret": "secret-xyz789"
}
```

Tokens emitidos para o FocusHub terão `app_id: "focushub-abc123"` no payload JWT.

---

## Requisitos Não-Funcionais Críticos

O auth-service é **mais crítico que qualquer outro serviço** — sua indisponibilidade bloqueia 100% dos produtos.

| Requisito | Valor | Justificativa |
|-----------|-------|---------------|
| Disponibilidade | **99,9%** | Todos os produtos bloqueados se cair |
| Latência P95 | < 200ms | Validação de OTP em cada login |
| Envio OTP | < 5s | UX — usuário está esperando |
| Volume esperado | 50k req/dia → 500k/dia (24 meses) | Escala com adoção dos produtos |
| Segurança | RS256, OTP hashed, HTTPS obrigatório | Dados sensíveis de acesso |

---

## LGPD — Dados Pessoais Tratados

Conforme OKR 3 (documentar dados pessoais por serviço até set/2026):

| Dado | Finalidade | Retenção | Base Legal |
|------|-----------|---------|-----------|
| email | Identificação e envio de OTP | Enquanto conta ativa; deletar em 30 dias após solicitação | Legítimo interesse / Consentimento |
| IP de acesso (sessions) | Segurança e auditoria | 90 dias | Legítimo interesse |
| Timestamps de login | Auditoria | 90 dias | Legítimo interesse |
| OTP hash | Autenticação | Deletado após uso ou expiração (10 min) | Execução de contrato |

**Direitos do usuário (Art. 18 LGPD):**
- `DELETE /v1/auth/users/:id` — soft delete com anonimização em 30 dias
- `GET /v1/auth/me` — acesso a todos os dados armazenados

---

## Dependências de Infraestrutura

```
auth-service depende de:
├── PostgreSQL (users, sessions, otp_codes, apps)
├── Redis (OTP cache + token blacklist — TTL automático)
└── SMTP/Resend (envio de email OTP — migrar para notification-service)

auth-service NÃO depende de:
├── customer-service (não gerencia perfis de negócio)
├── notification-service (atualmente — migração planejada v1.5)
└── Nenhum outro serviço IIT (é a fundação)
```

---

## Custo Operacional (Análise Financeira IIT)

| Item | Custo/mês | Observação |
|------|-----------|-----------|
| Infra K3s (compartilhada) | R$ 20 | Namespace production |
| Resend (OTP emails) | R$ 50 | ~10k emails OTP/mês |
| **Total** | **R$ 70/mês** | |

**Break-even vs. Auth0:** Auth0 custa R$ 115/mês para 1k MAU. IIT economiza R$ 45/mês desde o início — e escala melhor (sem custo por MAU).

---

## Pendências Críticas (Status Março 2026)

1. **Cloudflare Tunnel** — rota para `auth.institutoitinerante.com.br` ainda pendente de configuração
2. **Resend domain verification** — `institutoitinerante.com.br` não verificado; usando `apsf88@gmail.com` para testes
3. **Integração focus-hub** — auth pendente (bloqueante para uso multi-usuário do FocusHub)
4. **Integração iit-agents** — não iniciada
5. **GOTOOLCHAIN=auto** — necessário no Dockerfile (go.mod tem toolchain go1.24.13)

---

*Referência: IIT-PLATFORM-PRD.md v1.0.0 — Março 2026*
