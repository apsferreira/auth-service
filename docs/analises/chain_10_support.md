# Análise @support — auth-service-v2

## Alinhamento com @data, @finance e @qa
A equipe de Suporte está na linha de frente, lidando diretamente com as consequências das falhas do sistema. A falta de testes ("risco P0" por @data) se traduz em um volume maior de tickets e em problemas de difícil diagnóstico. O custo de um incidente (citado por @finance) também envolve o custo de horas do time de suporte gerenciando a crise. A clareza dos cenários de teste de @qa ajuda a antecipar os problemas que os usuários enfrentarão.

## Playbooks de Suporte Essenciais
Para garantir um atendimento rápido, consistente e seguro, os seguintes playbooks devem ser criados e disponibilizados para a equipe de suporte.

### Playbook 1: Conta Bloqueada por Suspeita de Brute Force
-   **Sintoma relatado pelo usuário:** "Não consigo fazer login, o sistema diz para eu esperar."
-   **Diagnóstico:**
    1.  Verificar nos logs (usando o `account_id` do usuário) se há múltiplos eventos de `auth.login.failed` de um mesmo IP em um curto período.
    2.  Confirmar se o rate limiting foi ativado para o IP ou a conta.
-   **Resolução:**
    1.  **NÃO desativar o bloqueio manualmente.**
    2.  Informar ao usuário que o bloqueio é uma medida de segurança temporária.
    3.  Instruí-lo a usar o fluxo de "esqueci minha senha" (se aplicável) para garantir que ele tem controle do email.
    4.  Aconselhar o usuário a verificar a segurança de sua conta de email e a não reutilizar senhas.

### Playbook 2: Reset de Senha / Recuperação de Acesso
-   **Sintoma relatado pelo usuário:** "Não recebo o email/mensagem para resetar minha senha/fazer login."
-   **Diagnóstico:**
    1.  Verificar se o `account_id` do usuário está ativo no sistema.
    2.  Confirmar no log do `notification-service` se a notificação (OTP/reset) foi enviada com sucesso para o email/telefone cadastrado.
-   **Resolução:**
    1.  Pedir ao usuário para verificar a caixa de spam/lixo eletrônico.
    2.  Pedir para confirmar se o email/telefone de destino `j****@e****.com` está correto.
    3.  Se tudo falhar, e após uma verificação de identidade rigorosa (ex: confirmar dados pessoais do `customer-service`), um administrador pode disparar o envio manualmente. **O suporte NUNCA deve definir uma senha para o usuário.**

### Playbook 3: Suspeita de Acesso Indevido (Account Takeover)
-   **Sintoma relatado pelo usuário:** "Recebi um email de login que não fui eu", "Minha senha/email foi alterado".
-   **Diagnóstico (ALERTA DE SEGURANÇA):**
    1.  Este é um incidente de segurança ativo. Escalar para o time de segurança/backend imediatamente.
    2.  Verificar os logs de `auth.login.success` da conta. Cruzar os IPs e User Agents com os que o usuário reconhece como seus.
-   **Resolução (Ação de Emergência):**
    1.  Um administrador deve **imediatamente revogar todas as sessões ativas** da conta (invalidar todos os refresh tokens).
    2.  Forçar um logout global do usuário.
    3.  Comunicar ao usuário para que ele inicie o fluxo de recuperação de acesso, preferencialmente por um canal seguro.
    4.  Aconselhar fortemente a ativação do 2FA, se disponível.

Estes playbooks, baseados em dados claros (`auth_event` de @data), permitem que o Suporte resolva problemas de forma eficaz sem comprometer a segurança.
