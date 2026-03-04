# PRD Enriquecido — auth-service-v2

**Versão:** 2.0  
**Data:** 2026-03-03  
**Status:** **Emergência Técnica - Ação Imediata Requerida**

> **Veredito do PM:** A análise cross-funcional revelou um consenso unânime: o estado atual do `auth-service` representa o **risco P0 mais crítico para a empresa**. A falta de testes e de versionamento é inaceitável para o serviço que guarda a porta de entrada de todo o nosso ecossistema. O custo financeiro e de reputação de um incidente, como detalhado pela análise, seria catastrófico. Portanto, o plano de ação delineado pela liderança técnica não é uma sugestão, mas a **prioridade máxima da companhia**, com efeito imediato.

---

## 1. Visão Estratégica

O `auth-service` não é apenas um utilitário; é a **fundação da confiança** do nosso ecossistema. Ele deve garantir segurança de nível militar, uma experiência de usuário sem fricção e ser um habilitador para o crescimento, permitindo a expansão segura para novos produtos e modelos de negócio através de SSO.

-   **North Star Metric:** 0 incidentes de segurança em produção.

---

## 2. Situação Atual (Resumo da Análise)

-   **Risco Identificado:** **P0 (Crítico)**. O serviço opera com **0% de cobertura de testes** e **sem repositório de código remoto**, tornando-o uma "bomba-relógio" técnica.
-   **Impacto Financeiro:** Um incidente de segurança tem um custo estimado de **R$50k a R$500k**, sem contar o custo de um downtime total do ecossistema.
-   **Risco Legal:** A situação atual constitui **negligência** na proteção de dados sob a LGPD, invalidando qualquer defesa em caso de vazamento.
-   **Dívida Técnica:** A decisão de não implementar boas práticas desde o início resultou em uma dívida que agora precisa ser paga com juros, em regime de urgência.

---

## 3. Plano de Ação (Resumo do Tech Lead)

Este plano é mandatório e tem precedência sobre outras iniciativas.

### Fase 1: Contenção (Deadline: 24 Horas)
1.  **Versionamento Imediato:** Criar repositório no GitHub, fazer push do código e proteger a branch `main`.
2.  **CI Básica:** Configurar pipeline que rode `build` e `lint` a cada commit.

### Fase 2: Estabilização (Deadline: 2 Semanas)
1.  **Cobertura de Testes:** Atingir **100% de cobertura nos fluxos críticos** e no mínimo 80% de cobertura geral.
2.  **CI/CD Completa:** O pipeline deve rodar todos os testes e verificar a cobertura. **Nenhum deploy será permitido se o pipeline falhar.**

### Fase 3: Melhoria e Auditoria (Deadline: 1 Mês)
1.  **Auditoria de Segurança:** Contratar pentest externo.
2.  **Endurecimento Técnico:** Implementar as melhorias de segurança definidas (JWT RS256, 2FA, etc.).
3.  **Processos de Suporte:** Documentar e treinar a equipe nos playbooks de incidentes.

---

## 4. Requisitos Funcionais e Técnicos (V2)

### 4.1. Core de Autenticação
-   [ ] **JWT com RS256:** Usar assinatura assimétrica para segurança máxima.
-   [ ] **Refresh Token Rotation:** Mitigar risco de reuso de tokens roubados.
-   [ ] **Revogação de Sessão:** Manter uma `deny-list` para revogação efetiva de tokens.
-   [ ] **Fluxo OTP-first:** Priorizar Magic Link (OTP) para a melhor experiência do usuário e segurança. Login por senha apenas para painéis de administração.
-   [ ] **2FA Obrigatório (TOTP):** Para todas as contas administrativas.

### 4.2. Infraestrutura e DevOps
-   [ ] **CI/CD Robusta:** Pipeline completo com testes, coverage check e build de imagem Docker.
-   [ ] **Deploy Zero-Downtime:** Usar estratégia `RollingUpdate` no K3s.
-   [ ] **Health Checks:** Implementar probes `liveness` e `readiness`.
-   [ ] **Gestão de Segredos:** Nenhuma credencial no código. Usar Kubernetes Secrets com política de rotação.
-   [ ] **Rate Limiting no Gateway:** Primeira linha de defesa contra brute force via Traefik.

### 4.3. Testes e Qualidade
-   [ ] **Cobertura de 100%:** Para lógica de negócio crítica.
-   [ ] **Testes de Segurança:** Testar exaustivamente cenários de brute force, token replay e abuso de sessão.

### 4.4. Dados e Monitoramento
-   [ ] **Log Estruturado:** Coletar os `auth_events` com todos os campos definidos (`ip_address`, `user_agent`, etc.).
-   [ ] **Métricas:** Dashboards para monitorar taxa de sucesso, latência p95, tentativas de brute force e tokens revogados.
-   [ ] **Event Sourcing:** Publicar eventos (`auth.login.success`, `auth.login.failed`) no message broker.

### 4.5. Segurança Proativa (AI)
-   [ ] **MVP de Detecção de Anomalias:** Desenvolver serviço que consome eventos de login e identifica padrões suspeitos (geolocalização, ASN, horário).
-   [ ] **Score de Risco por Sessão:** Implementar um score que possa acionar "step-up authentication" em caso de atividades de risco.

### 4.6. Conformidade Legal e Suporte
-   [ ] **Retenção de Logs:** Implementar política de expurgo automático de logs de autenticação após 6 meses.
-   [ ] **Playbooks de Suporte:** Documentar e treinar os cenários de conta bloqueada, recuperação de acesso e suspeita de invasão.

---

## 5. Endpoints (Atualizado)

Nenhuma mudança nos endpoints existentes, mas o comportamento e a segurança subjacente serão drasticamente melhorados. A especificação OpenAPI/Swagger torna-se um entregável da Fase 3.

---

Este PRD revisado reflete a maturidade que o `auth-service` deve atingir. Ele passa de um componente funcional, porém frágil, para o pilar de segurança e confiança de todo o nosso ecossistema de produtos. A execução do plano de ação é a tarefa mais importante para a engenharia no momento.
