# ⚠️ AVISO IMPORTANTE - VERSÃO BETA

## 🚧 Este projeto está em Beta de Longa Duração

### O que você PRECISA saber antes de usar:

## 📅 Timeline

- **Início Beta:** Novembro 2025
- **Duração Estimada:** 6-12 meses
- **Versão Estável:** v1.0.0 prevista para Q2 2027
- **Versão Atual:** v0.1.0-beta

## ⚠️ Mudanças Esperadas

### API NÃO é estável

```go
// v0.1.0-beta
gaver module:crud users User

// v0.2.0-beta (pode mudar!)
gaver generate crud users User --with-auth

// v0.3.0-beta (pode mudar novamente!)
gaver scaffold users User --full
```

**Consequências:**
- Comandos podem ser renomeados
- Estrutura de arquivos gerados pode mudar
- Templates podem ter alterações significativas
- Breaking changes entre versões beta

## ❌ NÃO USE Para:

- ✖️ Aplicações em produção
- ✖️ Projetos comerciais críticos
- ✖️ Sistemas com dados sensíveis
- ✖️ Projetos que precisam de estabilidade

## ✅ USE Para:

- ✔️ Aprendizado e experimentação
- ✔️ Protótipos e MVPs
- ✔️ Projetos pessoais
- ✔️ Testes e desenvolvimento
- ✔️ Dar feedback e contribuir

## 🐛 Bugs Esperados

Durante o beta, você **vai** encontrar bugs:
- Parser de annotations pode falhar em casos edge
- Migrations podem não detectar todas mudanças
- Templates podem gerar código inválido
- Performance não otimizada

**Isso é normal e esperado!** Por favor, reporte no GitHub.

## 🔄 Breaking Changes

### Vão acontecer!

Entre v0.1.0-beta e v0.5.0-beta, esperamos:
- 10-20 breaking changes
- Mudanças na estrutura de comandos
- Alterações nos templates
- Refatoração de APIs

### Como lidar:

1. **Fixe a versão** no seu projeto:
   ```bash
   go install github.com/Dalistor/gaver/cmd/gaver@v0.1.0-beta
   ```

2. **Leia CHANGELOG** antes de atualizar:
   ```bash
   # Ver o que mudou
   git diff v0.1.0-beta v0.2.0-beta CHANGELOG.md
   ```

3. **Teste em branch separada**:
   ```bash
   git checkout -b test-gaver-v0.2
   go install github.com/Dalistor/gaver/cmd/gaver@v0.2.0-beta
   # Testar mudanças
   ```

## 📊 Progresso até v1.0.0

```
Fase Beta Atual: ▓░░░░░░░░░ 10%

v0.1.0-beta ✅ (você está aqui)
v0.2.0-beta ⏳ (Q1 2026)
v0.3.0-beta ⏳ (Q2 2026)
v0.4.0-beta ⏳ (Q3 2026)
v0.5.0-beta ⏳ (Q4 2026)
v0.9.0-beta ⏳ (Q1 2027)
v1.0.0     ⏳ (Q2 2027)
```

## 💬 Seu Feedback é Essencial!

A versão beta existe para:
1. **Testar** ideias e abordagens
2. **Coletar feedback** da comunidade
3. **Identificar** problemas cedo
4. **Iterar** rapidamente

### Como ajudar:

- 🐛 **Reporte bugs**: Abra uma issue
- 💡 **Sugira features**: Discussions no GitHub
- 📝 **Melhore docs**: Pull requests são bem-vindos
- ⭐ **Dê uma estrela**: Ajuda o projeto a crescer
- 🗣️ **Compartilhe**: Fale sobre o projeto

## 🎯 Quando Usar em Produção?

**Espere até v1.0.0** se você precisa de:
- ✔️ API estável
- ✔️ Sem breaking changes
- ✔️ Documentação completa
- ✔️ Suporte long-term
- ✔️ Testes extensivos
- ✔️ Performance otimizada

**Pode experimentar agora** se:
- ✔️ Aceita riscos de breaking changes
- ✔️ Quer contribuir com desenvolvimento
- ✔️ Projeto não é crítico
- ✔️ Pode atualizar código quando necessário

## 📞 Suporte

Durante o beta:
- GitHub Issues para bugs
- GitHub Discussions para dúvidas
- CHANGELOG para breaking changes
- Sem SLA ou garantias

## ⏰ Quando v1.0.0 Será Lançado?

**Resposta curta:** Quando estiver pronto.

**Resposta longa:** 
- Estimativa: Q2 2027 (12-18 meses)
- Dependente de feedback e qualidade
- Não apressaremos o lançamento
- Preferimos estável e tarde do que cedo e bugado

---

## ✅ Você Foi Avisado!

Ao usar Gaver Framework v0.1.0-beta, você concorda que:
1. Entende que é uma versão beta instável
2. Aceita que a API pode mudar
3. Não usará em produção crítica
4. Reportará bugs encontrados
5. Terá paciência com o desenvolvimento

**Se concordar, aproveite e bem-vindo ao beta! 🎉**

---

**Última atualização:** Novembro 2025
**Versão do aviso:** v0.1.0-beta

