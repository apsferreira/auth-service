# Análise @legal — auth-service-v2

## Alinhamento com @data e @finance
A "situação de risco P0 (Crítico)" identificada por @data e a análise de custo de incidente de "R$50k-500k" de @finance são a materialização dos riscos legais e de conformidade que o `auth-service` atualmente representa. A ausência de testes e versionamento constitui negligência na proteção de dados pessoais, um fator agravante sob a ótica da LGPD.

## Implicações da LGPD (Lei Geral de Proteção de Dados)

1.  **Logs de Autenticação como Dados Pessoais:** Os dados coletados, como `ip_address`, `user_agent`, e a associação com `account_id`, são considerados dados pessoais. A sua coleta é legítima para fins de segurança, mas impõe responsabilidades.
    -   **Política de Retenção:** Logs operacionais para fins de segurança devem ser mantidos por um período justificado e limitado. Recomenda-se a retenção por no máximo **6 meses**, com expurgo automático após esse período para minimizar a superfície de exposição.

2.  **Responsabilidade em Caso de Incidente:** Em um vazamento de dados, a empresa é responsável por demonstrar que empregou as "medidas de segurança, técnicas e administrativas aptas a proteger os dados pessoais". A falta de uma suíte de testes e de um repositório versionado, como apontado por @data, invalida qualquer argumento de defesa de que as medidas adequadas foram tomadas.

3.  **Comunicação de Incidentes:** A LGPD exige a comunicação à ANPD e aos titulares dos dados em caso de incidente de segurança que possa acarretar risco ou dano relevante. Os custos financeiros citados por @finance estão diretamente ligados a esta obrigação.

## Requisitos de Conformidade

1.  **2FA (Autenticação de Dois Fatores) Obrigatório:** Para usuários com acesso a painéis administrativos ou que operam dados sensíveis (operadores de sistemas, administradores de tenant), a autenticação de dois fatores deve ser **obrigatória e não opcional**. Isso reduz drasticamente o risco de comprometimento de contas privilegiadas.

2.  **Prova de Consentimento e Acesso:** O sistema deve ser capaz de fornecer a um usuário, mediante solicitação, todos os logs de acesso (`auth_event`) associados à sua conta, como parte do direito de acesso aos seus dados.

**Conclusão Legal:** A estabilização do `auth-service` não é uma opção, mas uma obrigação legal. A situação atual expõe a empresa a sanções severas e a uma perda de reputação que pode comprometer a operação como um todo. A implementação das recomendações de @data é o primeiro passo para a mitigação do risco jurídico.
