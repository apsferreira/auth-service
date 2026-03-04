# Análise @backend — auth-service-v2

## Construindo sobre @data, @finance e @legal
A análise de @data é clara: o "risco P0 (Crítico)" devido à "0% de cobertura de testes" e "sem repositório remoto" é a dívida técnica mais urgente do ecossistema. O impacto financeiro de "R$50k-500k" por incidente, destacado por @finance, e as obrigações da LGPD, apontadas por @legal, tornam a ação corretiva uma emergência.

## Plano de Ação Técnico

1.  **JWT (JSON Web Tokens):**
    *   **Algoritmo:** Adotar **RS256** (assinatura assimétrica) em vez de HS256. A chave privada fica exclusivamente no `auth-service`, enquanto a chave pública é distribuída para os serviços consumidores. Isso garante que apenas o `auth-service` possa emitir tokens, e os outros serviços possam apenas validá-los, sem risco de comprometimento da chave de assinatura.
    *   **Vida útil:** Manter Access Tokens com vida curta (ex: 15 minutos) e Refresh Tokens com vida mais longa (ex: 7 dias).

2.  **Segurança de Tokens:**
    *   **Refresh Token Rotation:** Implementar a rotação de refresh tokens. A cada uso, um novo refresh token é emitido e o anterior é invalidado. Isso detecta e previne o reuso de tokens roubados.
    *   **Revogação:** Manter uma `deny-list` (em Redis, com TTL) para refresh tokens que foram explicitamente revogados via logout.

3.  **Mecanismos de Defesa:**
    *   **Rate Limiting:** Aplicar rate limiting por IP em endpoints sensíveis como `/login`, `/otp/request` e `/otp/verify` para mitigar ataques de brute force, conforme a necessidade identificada pela métrica de "Tentativas de Brute Force Bloqueadas" de @data.
    *   **M2M (Machine-to-Machine) Auth:** Para comunicação entre serviços (ex: `customer-service` provisionando um usuário), usar um token estático de serviço (`X-Service-Token`) ou, preferencialmente, o padrão Client Credentials do OAuth2.

4.  **Conformidade Legal:**
    *   **2FA:** Implementar o fluxo de setup e verificação de 2FA (TOTP - Time-based One-Time Password) para contas administrativas, como exigido por @legal.

## Ação Imediata (Hoje)
A situação descrita por @data é inaceitável. As seguintes ações devem ser executadas **imediatamente, sem adiamento**:
1.  **Criar o repositório remoto no GitHub.**
2.  **Fazer o push do código atual para a branch `main`.**
3.  **Proteger a branch `main` contra pushes diretos.**
4.  **Criar um pipeline de CI/CD básico (GitHub Actions) que rode `go build` e `go vet` a cada commit.**

O desenvolvimento de testes deve começar logo em seguida, mas a garantia de que o código-fonte está seguro e versionado é a prioridade zero.
