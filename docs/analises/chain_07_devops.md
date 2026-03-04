# Análise @devops — auth-service-v2

## Alinhamento com @data, @finance e @backend
A perspectiva de DevOps foca em traduzir os requisitos de negócio e de sistema em uma infraestrutura confiável, segura e escalável. A "situação de risco P0 (Crítico)" de @data é um alarme direto de falha de processo de DevOps. A ação imediata proposta por @backend ("criar remote GitHub hoje + push") é o passo zero e fundamental. O custo de downtime apontado por @finance reforça a necessidade de alta disponibilidade e deploys zero-downtime.

## Plano de Ação de Infraestrutura e CI/CD

1.  **CI/CD Pipeline (GitHub Actions):**
    *   **Trigger:** A cada push para `main` ou em Pull Requests.
    *   **Etapas:**
        1.  **Lint & Vet:** Análise estática do código (`golangci-lint`, `go vet`).
        2.  **Test:** Execução da suíte de testes unitários e de integração (`go test ./... -coverprofile=coverage.out`).
        3.  **Coverage Check:** Falhar o build se a cobertura for menor que o threshold definido (iniciar com 80%, mirar 95%+).
        4.  **Build:** Compilar o binário Go em um ambiente determinístico (mesma versão do Go).
        5.  **Build & Push Docker Image:** Construir a imagem Docker do serviço e publicá-la em um registro (ex: GitHub Container Registry, Docker Hub).

2.  **Deployment em K3s (Produção):**
    *   **Estratégia de Deploy:** `RollingUpdate` para garantir zero downtime. O Kubernetes irá gradualmente substituir os Pods antigos pelos novos, aguardando que os novos estejam prontos (readiness probe) antes de remover os antigos.
    *   **Health Checks:**
        -   **Liveness Probe:** `GET /health` - Se falhar, o Kubernetes reinicia o Pod.
        -   **Readiness Probe:** `GET /health` - Se falhar, o Kubernetes para de enviar tráfego para o Pod, mas não o reinicia. Essencial para o zero-downtime deploy.
    *   **ConfigMaps e Secrets:**
        -   **ConfigMaps:** Para configurações não sensíveis (ex: portas, URLs de serviços).
        -   **Secrets:** Para todas as credenciais (chaves de API, segredos de JWT, connection strings de banco de dados). **As chaves NUNCA devem estar no código-fonte.**

3.  **Gestão de Segredos:**
    *   **Secrets Rotation:** Implementar uma política e um processo para rotacionar todos os segredos críticos (chaves de JWT, senhas de banco de dados) em intervalos regulares (ex: a cada 90 dias). Ferramentas como Vault podem automatizar isso, mas um processo manual bem definido é o mínimo aceitável.
    *   **Acesso em CI/CD:** Usar os secrets do GitHub Actions para injetar credenciais de forma segura durante o build e deploy.

4.  **Segurança de Rede (API Gateway - Traefik):**
    *   **Rate Limiting:** A primeira linha de defesa contra brute force deve ser no API Gateway. Configurar o middleware do Traefik para aplicar rate limiting por IP nos endpoints críticos do `auth-service`, como recomendado por @backend. Isso protege o serviço de ser sobrecarregado.
    *   **TLS Termination:** O Traefik será responsável por terminar a conexão TLS, garantindo que o tráfego externo seja sempre criptografado (HTTPS). A comunicação interna na rede do cluster pode ser HTTP.
