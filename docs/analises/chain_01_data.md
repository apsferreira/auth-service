# Análise @data — auth-service-v2

## NSM (North Star Metric)
A métrica principal de segurança para o `auth-service` é **0 incidentes de segurança em produção**. Este é um objetivo binário e inegociável. Qualquer incidente, como vazamento de credenciais, acesso não autorizado ou comprometimento de tokens, representa uma falha total da nossa NSM.

## Métricas de Suporte
Para garantir a saúde e a performance do serviço, monitoramos as seguintes métricas:

1.  **Taxa de Autenticação Bem-Sucedida:** `% de logins/verificações de OTP com sucesso`. Uma queda súbita pode indicar um ataque ou uma falha no sistema.
2.  **Latência de Login (p95):** `p95 < 200ms`. O tempo de resposta para endpoints críticos como `/login` e `/verify-otp` deve ser rápido para garantir uma boa experiência ao usuário.
3.  **Tentativas de Brute Force Bloqueadas:** `Contagem de IPs bloqueados por rate limiting`. Um aumento indica tentativas de ataque ativas.
4.  **Tokens Revogados:** `Contagem de refresh tokens revogados (logout)`. Ajuda a entender o ciclo de vida da sessão do usuário.

## Coleta de Dados Essencial (Desde o Dia 1)
Para que a análise de segurança e performance seja viável, a seguinte estrutura de dados deve ser registrada para cada evento de autenticação:

-   `auth_event_id` (UUID): Identificador único do evento.
-   `account_id` (UUID): ID do usuário associado ao evento.
-   `event_type` (Enum): `login_success`, `login_failed`, `logout`, `token_refresh`, `token_revoked`, `password_reset_request`, `password_reset_success`.
-   `ip_address` (String): IP de origem da requisição.
-   `user_agent` (String): User agent do cliente.
-   `created_at` (TimestampZ): Horário exato do evento.

## Eventos Publicados (Message Broker)
Para integração com outros serviços (como detecção de anomalias), o `auth-service` deve publicar os seguintes eventos:

-   `auth.login.success`
-   `auth.login.failed`
-   `auth.token.revoked`
-   `auth.password.reset`

## Situação Crítica Atual
A análise revela uma situação de **risco P0 (Crítico)**:
-   **0% de cobertura de testes:** O serviço mais crítico do ecossistema não possui uma única asserção automatizada, tornando qualquer alteração imprevisível e perigosa.
-   **Sem repositório remoto no GitHub:** O código-fonte existe apenas localmente, sem versionamento, colaboração ou backup. Uma falha de hardware local resultaria na perda total do serviço.

Esta combinação é inaceitável para um serviço de autenticação e deve ser tratada como a maior prioridade técnica da empresa.
