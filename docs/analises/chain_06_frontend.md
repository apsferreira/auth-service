# Análise @frontend — auth-service-v2

## Alinhamento com as demais áreas
A perspectiva do Frontend é construir uma experiência de usuário que seja segura, fluida e que materialize as estratégias definidas pelas outras áreas. A análise de @growth sobre **Magic Link vs. Senha** é o principal direcionador da UI/UX. A necessidade de 2FA apontada por @legal e as especificações de tokens de @backend definem os requisitos técnicos. O custo de um incidente (citado por @finance) reforça a necessidade de tratar os tokens com o máximo de cuidado no cliente.

## Fluxos de Usuário Essenciais a Serem Construídos/Revisados

1.  **Fluxo de Login/Registro (OTP-first):**
    *   **UI:** Um único campo de entrada para "Email ou Telefone".
    *   **UX:** O usuário insere o contato e clica em "Continuar". O backend determina o canal (email, WhatsApp, etc.) e envia o OTP. A tela seguinte é para inserir o código recebido. Simples e com baixa fricção.
    *   **Técnica:** Após a verificação do OTP, o JWT (access e refresh token) recebido deve ser armazenado de forma segura.

2.  **Gerenciamento de Sessão:**
    *   **Armazenamento de Token:**
        -   **Access Token:** Armazenar em memória (ex: estado de um componente React/Vue) para evitar exposição a ataques XSS. Não usar `localStorage`.
        -   **Refresh Token:** Armazenar em um cookie `HttpOnly`, `Secure`, com `SameSite=Strict`. Isso impede que o JavaScript do cliente o acesse, protegendo-o contra XSS, e o protege contra ataques CSRF.
    *   **Renovação de Sessão:** Implementar um interceptor de requisições (ex: com Axios) que, ao receber um erro 401 (Unauthorized), tente usar o refresh token para obter um novo access token de forma transparente para o usuário. Se a renovação falhar, o usuário é redirecionado para a tela de login.

3.  **Fluxo de "Esqueci Minha Senha" (Admin):**
    *   Interface padrão de solicitação de reset por email, com envio de um link com token de uso único e tempo de vida curto.

4.  **Setup de 2FA (Admin):**
    *   **UI:** Apresentar um QR code (gerado a partir do segredo TOTP fornecido pelo backend) para ser escaneado por aplicativos como Google Authenticator ou Authy.
    *   **UX:** Incluir um campo para o usuário inserir o primeiro código gerado para confirmar que o setup foi bem-sucedido. Apresentar códigos de recuperação (backup) para o usuário salvar em local seguro.

5.  **Logout:**
    *   A ação de logout deve invalidar os tokens no cliente e chamar o endpoint de `/logout` do backend para revogar o refresh token no servidor, garantindo que a sessão seja totalmente encerrada.
