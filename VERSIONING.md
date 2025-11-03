# 📌 Guia de Versionamento - Gaver Framework

## 🎯 Versão Atual: `v0.1.0-beta`

Este documento explica como funciona o versionamento do Gaver Framework.

**Fase:** Beta Testing (Long-Term)
**Duração Estimada:** 6-12 meses
**Versão Estável Prevista:** v1.0.0 em Q2 2027

## 📋 Semantic Versioning

Usamos **Semantic Versioning 2.0.0** (https://semver.org/)

### Formato: `MAJOR.MINOR.PATCH[-PRERELEASE]`

Exemplo: `v1.2.3-beta.1`
- `1` = MAJOR (mudanças incompatíveis na API)
- `2` = MINOR (novas funcionalidades compatíveis)
- `3` = PATCH (correções de bugs)
- `-beta.1` = PRE-RELEASE (versão de teste)

## 🏷️ Tags de Pré-Lançamento

### Beta (Fase Atual - Long-Term)
```
v0.1.0-beta      # Core framework (atual)
v0.2.0-beta      # QuerySet e validações
v0.3.0-beta      # Developer experience
v0.4.0-beta      # Features avançadas
v0.5.0-beta      # Produção-ready features
v0.9.0-beta      # Feature freeze
```

**Características:**
- ⚠️ API pode mudar entre versões
- 🧪 Para testes, desenvolvimento e feedback
- ❌ Não recomendado para produção
- 📊 Versionamento por features, não por patches
- 🔄 Breaking changes permitidos entre minor versions

**Política de Breaking Changes:**
- Permitido entre v0.x.0-beta e v0.y.0-beta
- Documentado no CHANGELOG
- Anunciado com antecedência quando possível

### Alpha (Desenvolvimento Inicial)
```
v0.1.0-alpha.1   # Muito instável
v0.1.0-alpha.2   # Ainda em desenvolvimento
```

### Release Candidate
```
v1.0.0-rc.1      # Candidato a lançamento
v1.0.0-rc.2      # Quase pronto para produção
```

### Stable (Produção)
```
v1.0.0           # Primeira versão estável
v1.1.0           # Nova feature
v1.1.1           # Bug fix
v2.0.0           # Breaking change
```

## 🚀 Como Criar uma Nova Versão Beta

### Para Nova Feature (Minor Version)

```bash
# 1. Implementar features

# 2. Atualizar VERSION
echo "0.2.0-beta" > VERSION

# 3. Atualizar CHANGELOG.md
# Adicionar nova seção:
## [0.2.0-beta] - 2026-02-XX
### Adicionado
- QuerySet API completo
- Validações cross-field
### Mudanças
- Breaking: Alterada estrutura de callbacks
### Corrigido
- 15 bugs diversos

# 4. Commit
git add .
git commit -m "feat: release v0.2.0-beta

- QuerySet API
- Validações avançadas
- Breaking changes documentados no CHANGELOG
"

# 5. Criar tag
git tag -a v0.2.0-beta -m "v0.2.0-beta - QuerySet e Validações"

# 6. Push
git push origin main
git push origin v0.2.0-beta
```

### Para Bug Fix Crítico (Patch - Raro)

```bash
# Apenas para bugs que impedem uso
echo "0.1.1-beta" > VERSION
git commit -m "fix: corrige bug crítico X"
git tag -a v0.1.1-beta -m "Hotfix crítico"
git push origin main v0.1.1-beta
```

**Nota:** Preferimos acumular fixes para próxima versão minor.

## 📦 Publicação

### Primeira Publicação (Beta)

```bash
# 1. Certifique-se que tudo está commitado
git status

# 2. Crie a tag beta
git tag -a v0.1.0-beta -m "Initial beta release

Core framework functionality:
- Module system
- CRUD generation
- Annotations gaverModel
- Smart migrations
- Gin integration
- Scheduled routines
"

# 3. Push tudo
git push origin main
git push origin v0.1.0-beta
```

### Usuários podem instalar com:

```bash
# Última versão (sempre beta durante desenvolvimento)
go install github.com/Dalistor/gaver/cmd/gaver@latest

# Versão específica
go install github.com/Dalistor/gaver/cmd/gaver@v0.1.0-beta

# Versão mais recente beta
go install github.com/Dalistor/gaver/cmd/gaver@v0.5.0-beta
```

### Política de @latest Durante Beta

Durante a fase beta, `@latest` sempre apontará para a versão beta mais recente:
- Agora: `@latest` = `v0.1.0-beta`
- Futuro: `@latest` = `v0.5.0-beta`
- Após v1.0.0: `@latest` = versão estável mais recente

## 🔄 Ciclo de Desenvolvimento

### Durante Beta (v0.x.0-beta) - Simplificado

**Versionamento por Features (não por patches):**

1. **Qualquer mudança**: Nova versão beta
   - `v0.1.0-beta` → `v0.2.0-beta` (nova feature)
   - `v0.2.0-beta` → `v0.3.0-beta` (mais features)

2. **Breaking changes**: Permitidos e documentados
   - `v0.3.0-beta` → `v0.4.0-beta` (pode ter breaking changes)

3. **Bug fixes críticos**: Podem gerar releases pontuais
   - `v0.1.0-beta` → `v0.1.1-beta` (apenas se crítico)
   - Mas preferimos acumular fixes para próxima versão

**Filosofia:**
- Menos releases, mais features por release
- Breaking changes bem documentados
- Feedback da comunidade guia desenvolvimento

### Timeline de Saída do Beta

```
v0.1.0-beta   Nov 2025  ← VOCÊ ESTÁ AQUI
v0.2.0-beta   Q1 2026   (3-4 meses)
v0.3.0-beta   Q2 2026   (3-4 meses)
v0.4.0-beta   Q3 2026   (3-4 meses)
v0.5.0-beta   Q4 2026   (3-4 meses)
v0.9.0-beta   Q1 2027   (feature freeze)
v1.0.0-rc.1   Q1 2027   (release candidate)
v1.0.0        Q2 2027   (ESTÁVEL!)
```

**Total:** ~12-18 meses em beta

### Após v1.0.0 (Produção)

1. **Bug fix**: Incrementa PATCH
   - `v1.2.3` → `v1.2.4`

2. **Nova feature (compatível)**: Incrementa MINOR
   - `v1.2.3` → `v1.3.0`

3. **Breaking change**: Incrementa MAJOR
   - `v1.2.3` → `v2.0.0`

## 📊 Branches Recomendadas

```
main            # Versão estável (v1.x.x)
develop         # Desenvolvimento (v0.x.x-beta)
feature/*       # Features específicas
hotfix/*        # Correções urgentes
release/*       # Preparação de releases
```

## ⚙️ GitHub Releases

### Criar Release no GitHub

1. Acesse: `https://github.com/Dalistor/gaver/releases/new`

2. Escolha a tag: `v0.1.0-beta.1`

3. Título: `v0.1.0-beta.1 - First Beta Release`

4. Descrição:
   ```markdown
   ## 🎉 Primeira versão beta do Gaver Framework!
   
   ### ✨ Funcionalidades
   - Sistema de modules
   - Annotations gaverModel
   - CRUD automático com callbacks
   - Migrations inteligentes
   
   ### 📦 Instalação
   ```bash
   go install github.com/Dalistor/gaver/cmd/gaver@v0.1.0-beta.1
   ```
   
   ### ⚠️ Aviso
   Esta é uma versão beta. A API pode mudar.
   ```

5. Marque: ✅ **This is a pre-release**

6. Clique em **Publish release**

## 📝 Convenções de Commit

Use commits semânticos:

```
feat: adiciona novo comando X
fix: corrige bug no parser
docs: atualiza README
refactor: reorganiza código do generator
test: adiciona testes para validator
chore: atualiza dependências
```

## 🎯 Roadmap de Versões Planejado

```
Nov 2025    v0.1.0-beta   ← VOCÊ ESTÁ AQUI (Core Framework)
Q1 2026     v0.2.0-beta   (QuerySet & Validations)
Q2 2026     v0.3.0-beta   (DX & Examples)
Q3 2026     v0.4.0-beta   (Advanced Features)
Q4 2026     v0.5.0-beta   (Production-Ready)
Q1 2027     v0.9.0-beta   (Feature Freeze)
Q1 2027     v1.0.0-rc.1   (Release Candidate)
Q2 2027     v1.0.0        (STABLE!)
```

**Observações:**
- Timeline é flexível baseado em feedback
- Breaking changes permitidos entre versões beta
- Cada versão beta pode levar 3-4 meses
- RC phase pode ter múltiplas versões (rc.1, rc.2, etc)
- v1.0.0 só será lançada quando realmente estável

## 🔗 Links Úteis

- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [Go Modules Reference](https://go.dev/ref/mod)

