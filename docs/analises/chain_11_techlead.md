# Análise @techlead — auth-service-v2

## Síntese das Análises Anteriores
A liderança técnica consolida as visões de todas as áreas para tomar decisões estratégicas.
-   **@data e @qa:** Deixaram claro o "risco P0 (Crítico)" e a necessidade de 100% de cobertura de testes. Isso não é negociável.
-   **@finance e @legal:** Quantificaram o risco em termos financeiros e de conformidade, fornecendo o "porquê" de negócio para a ação imediata.
-   **@backend e @devops:** Definiram o "como" técnico: RS256, token rotation, CI/CD, K3s, e a ação emergencial de versionar o código.
-   **@ai:** Apresentou uma visão de futuro, uma camada de segurança proativa que só pode ser construída sobre uma fundação estável.

## Decisões Arquiteturais (ADR)

### ADR 1: JWT Stateless vs. Session Stateful
-   **Contexto:** Precisamos de um mecanismo de sessão que seja seguro e escalável horizontalmente.
-   **Opções:**
    1.  **Stateful:** Armazenar a sessão no servidor (ex: em Redis). Requer uma consulta ao Redis a cada requisição, mas permite revogação instantânea.
    2.  **Stateless (JWT):** A sessão (claims) está contida no próprio token. Altamente escalável, pois não requer consulta a um banco de dados central. A revogação é mais complexa.
-   **Decisão:** **Adotar um modelo híbrido, primariamente stateless com JWTs, complementado por uma `deny-list` para revogação.**
-   **Justificativa:** Este modelo combina a escalabilidade dos JWTs (recomendado por @backend) com um mecanismo de segurança essencial (revogação de refresh tokens). O Access Token de vida curta (15 min) é validado de forma stateless, minimizando a latência. O Refresh Token, de vida longa, é verificado contra a deny-list no momento da renovação, garantindo que sessões revogadas não possam ser estendidas.

## Plano de Ação Imediato
A situação atual é uma emergência técnica. O plano a seguir deve ser executado com prioridade máxima, congelando, se necessário, o desenvolvimento de novas features em outros projetos.

**Fase 1: Contenção (Próximas 24 horas)**
1.  **Versionamento:** Criar o repositório remoto no GitHub e fazer o push do código. (Responsável: @backend)
2.  **CI Básico:** Configurar a CI para rodar build e lint a cada commit. (Responsável: @devops)
3.  **Proteção de Branch:** Proteger a branch `main`. Todo o trabalho futuro deve ser feito via Pull Requests.

**Fase 2: Estabilização (Próximas 2 semanas)**
1.  **Fundação de Testes:** Desenvolver a estrutura de testes de integração com um banco de dados de teste. (Responsável: @qa)
2.  **Cobertura Crítica:** Escrever testes para cobrir 100% dos fluxos de login, geração e validação de JWT, e revogação. Atingir no mínimo 80% de cobertura geral. (Responsável: @backend, @qa)
3.  **Pipeline Completo:** Adicionar o passo de testes e verificação de cobertura à CI. Deploys para staging/produção só serão permitidos se o pipeline passar. (Responsável: @devops)

**Fase 3: Auditoria e Melhoria Contínua (Próximo Mês)**
1.  **Auditoria de Segurança:** Contratar uma auditoria de segurança externa (pentest) para validar a robustez da implementação. O custo disso, como apontado por @finance, é marginal comparado ao de um incidente.
2.  **Implementar Recomendações:** Adotar as melhorias de segurança propostas (RS256, 2FA, etc.). (Responsável: @backend)
3.  **Desenvolver Playbooks:** Documentar os playbooks de suporte. (Responsável: @support)
4.  **Iniciar Prova de Conceito de AI:** Começar a coletar e analisar os eventos de login para a detecção de anomalias. (Responsável: @ai)

Este plano transforma o `auth-service` de maior passivo da empresa em um ativo estratégico e um exemplo de excelência técnica a ser seguido pelos outros serviços.
