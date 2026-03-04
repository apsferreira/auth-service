# Análise @ai — auth-service-v2

## Alinhamento com @data, @finance e @backend
A área de AI entra como uma camada proativa de segurança, construindo sobre a fundação de dados robusta proposta por @data. Enquanto @backend e @devops focam em defesas determinísticas (rate limiting, senhas fortes), a AI pode detectar ameaças que escapam a essas regras. O impacto financeiro de um incidente (citado por @finance) justifica o investimento em inteligência para prevenção. Os eventos publicados no message broker, como `auth.login.failed`, são a matéria-prima essencial para esta análise.

## Oportunidades de Aplicação de AI

1.  **Detecção de Anomalias em Tempo Real:**
    *   **Objetivo:** Identificar comportamentos de login que fogem do padrão de um usuário e que podem indicar um "Account Takeover" (ATO).
    *   **Features para Análise:**
        -   **Geolocalização do IP:** O usuário sempre loga do Brasil e de repente há um login do Leste Europeu?
        -   **Provedor de Internet (ASN):** O padrão é um provedor residencial (ex: Vivo Fibra) e surge um login de um provedor de data center (ex: AWS, DigitalOcean), que pode indicar uso de VPNs, proxies ou servidores comprometidos.
        -   **Frequência e Horário:** Tentativas de login fora do horário habitual do usuário (ex: madrugada).
        -   **User-Agent:** Mudança drástica no navegador ou sistema operacional.
    *   **Implementação:**
        -   Um serviço consumidor escuta os eventos `auth.login.success` e `auth.login.failed`.
        -   Para cada evento, enriquece o dado com geolocalização do IP (ex: usando uma base de dados local como a da MaxMind).
        -   Mantém um perfil histórico simples do usuário (em Redis ou um banco de chave-valor) com os últimos 5-10 países, ASNs e horários de login.
        -   **Modelo:** Para o MVP, um sistema de regras simples já é eficaz. Para a V2, um modelo de LLM (como **Gemini Flash Lite**) pode ser usado para uma análise mais contextual. Exemplo de prompt: `"Analise este evento de login e retorne um score de risco de 0 a 100. Histórico do usuário: [últimos logins]. Evento atual: [dados do evento atual]."`.

2.  **Score de Risco por Sessão:**
    *   **Objetivo:** Atribuir um "score de risco" a cada sessão de usuário, que pode ser usado por outros serviços para tomar decisões.
    *   **Funcionamento:**
        -   Sessão inicia com score 0.
        -   Login de um IP/país incomum: `score += 20`.
        -   Múltiplas tentativas falhas antes do sucesso: `score += 15`.
        -   Acesso a uma funcionalidade crítica (ex: alterar email) pela primeira vez: `score += 10`.
    *   **Ação:** Se o score ultrapassa um limiar (ex: `> 50`), o sistema pode exigir uma re-autenticação (verificação de OTP) antes de permitir a continuação, implementando um "step-up authentication".

## Próximos Passos
O primeiro passo é garantir a coleta de dados de alta qualidade definida por @data. Com os eventos `auth.login.success` e `auth.login.failed` fluindo via message broker, um protótipo de detecção de anomalias pode ser desenvolvido em poucas semanas, adicionando uma camada de segurança inteligente e adaptativa ao ecossistema.
