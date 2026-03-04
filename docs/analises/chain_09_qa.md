# Análise @qa — auth-service-v2

## Alinhamento com as demais áreas
A perspectiva de QA (Quality Assurance) é a materialização das preocupações de todas as outras áreas em um plano de testes concreto. O "risco P0" de @data é o chamado à ação. O impacto financeiro de @finance justifica o rigor. Os requisitos de @backend, as defesas de @devops e a detecção de anomalias de @ai definem *o que* testar.

## Cobertura de Testes: 100% Obrigatório
Para a maioria dos serviços, uma cobertura de 80-90% é um objetivo razoável. Para um serviço de segurança como o `auth-service`, a meta deve ser **100% de cobertura de testes unitários para a lógica de negócio crítica**. Qualquer linha de código que decida se um usuário pode ou não acessar o sistema deve ser coberta por um teste automatizado. Não há espaço para "achismos" em autenticação.

## Plano de Testes por Camada

1.  **Testes Unitários (Go `testing` + `testify`):**
    *   **Foco:** Testar cada função isoladamente.
    *   **Exemplos:**
        -   Geração e validação de senhas com hash.
        -   Lógica de expiração de tokens.
        -   Validação de formato de email.
        -   Geração de códigos OTP.
        -   Lógica de rotação de refresh tokens.

2.  **Testes de Integração (Go `httptest`):**
    *   **Foco:** Testar o serviço como uma caixa-preta, interagindo com seus endpoints HTTP e uma base de dados de teste (em-memória ou um container Docker dedicado).
    *   **Exemplos:**
        -   Fluxo completo: `/request-otp` -> recebe código -> `/verify-otp` -> recebe JWT válido.
        -   Tentar usar um refresh token revogado e garantir que o acesso seja negado.
        -   Tentar acessar um recurso protegido com um token expirado e esperar um erro 401.
        -   Simular múltiplas falhas de login e verificar se o rate limiting é ativado.

## Casos de Teste de Segurança Críticos

Estes cenários devem ser testados de forma exaustiva:

1.  **Ataques de Brute Force:**
    *   **Cenário:** Múltiplas requisições para `/login` ou `/verify-otp` com credenciais inválidas a partir de um mesmo IP.
    *   **Resultado Esperado:** O serviço deve retornar um erro `429 Too Many Requests` após o N-ésimo pedido.

2.  **Token Replay:**
    *   **Cenário:** Capturar um refresh token, usá-lo para obter um novo par de tokens e, em seguida, tentar usar o mesmo refresh token novamente.
    *   **Resultado Esperado:** A segunda tentativa deve falhar, pois a rotação de tokens já invalidou o token original.

3.  **Controle de Sessões Concorrentes:**
    *   **Cenário:** Usuário loga no Dispositivo A. Depois, loga no Dispositivo B.
    *   **Resultado Esperado:** A sessão no Dispositivo A deve ser invalidada (se essa for a política de negócio) e seu refresh token revogado.

4.  **Validação de Claims do JWT:**
    *   **Cenário:** Tentar usar um token de um `tenant_id` para acessar recursos de outro.
    *   **Resultado Esperado:** Acesso deve ser negado.

A implementação deste plano de testes, integrada ao pipeline de CI/CD proposto por @devops, é a única forma de transformar a situação atual de "risco P0" em um estado de confiança e estabilidade.
