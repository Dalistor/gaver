# 📌 Guia de Versionamento - Gaver Framework

## 🎯 Versão Atual: `v0.1.0-beta.1`

Este documento explica como funciona o versionamento do Gaver Framework.

## 📋 Semantic Versioning

Usamos **Semantic Versioning 2.0.0** (https://semver.org/)

### Formato: `MAJOR.MINOR.PATCH[-PRERELEASE]`

Exemplo: `v1.2.3-beta.1`
- `1` = MAJOR (mudanças incompatíveis na API)
- `2` = MINOR (novas funcionalidades compatíveis)
- `3` = PATCH (correções de bugs)
- `-beta.1` = PRE-RELEASE (versão de teste)

## 🏷️ Tags de Pré-Lançamento

### Beta (Fase Atual)
```
v0.1.0-beta.1    # Primeira versão beta
v0.1.0-beta.2    # Segunda versão beta (correções)
v0.2.0-beta.1    # Nova feature em beta
```

**Características:**
- ⚠️ API pode mudar
- 🧪 Para testes e feedback
- ❌ Não use em produção

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

## 🚀 Como Criar uma Nova Versão

### 1. Atualizar VERSION
```bash
echo "0.1.0-beta.2" > VERSION
```

### 2. Atualizar CHANGELOG.md
```markdown
## [0.1.0-beta.2] - 2025-11-04

### Adicionado
- Nova funcionalidade X

### Corrigido
- Bug Y corrigido
```

### 3. Commit das mudanças
```bash
git add VERSION CHANGELOG.md
git commit -m "chore: bump version to v0.1.0-beta.2"
```

### 4. Criar tag Git
```bash
# Criar tag anotada (recomendado)
git tag -a v0.1.0-beta.2 -m "Release v0.1.0-beta.2

- Nova funcionalidade X
- Bug Y corrigido
"

# Push da tag
git push origin v0.1.0-beta.2
```

### 5. Push do código
```bash
git push origin main
```

## 📦 Publicação

### Primeira Publicação (Beta)

```bash
# 1. Certifique-se que tudo está commitado
git status

# 2. Crie a tag beta
git tag -a v0.1.0-beta.1 -m "Initial beta release"

# 3. Push tudo
git push origin main
git push origin v0.1.0-beta.1
```

### Usuários podem instalar com:

```bash
# Última versão beta
go install github.com/Dalistor/gaver/cmd/gaver@latest

# Versão específica
go install github.com/Dalistor/gaver/cmd/gaver@v0.1.0-beta.1
```

## 🔄 Ciclo de Desenvolvimento

### Durante Beta (v0.x.x-beta)

1. **Bug fix**: Incrementa último número
   - `v0.1.0-beta.1` → `v0.1.0-beta.2`

2. **Nova feature**: Incrementa MINOR
   - `v0.1.0-beta.1` → `v0.2.0-beta.1`

3. **Breaking change**: OK durante beta
   - `v0.1.0-beta.1` → `v0.2.0-beta.1`

### Quando sair do Beta

```bash
# Remover sufixo -beta
v0.1.0-beta.5 → v0.1.0 (primeira versão estável)

# Ou ir direto para v1.0.0
v0.5.0-beta.3 → v1.0.0 (lançamento oficial)
```

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

## 🎯 Roadmap de Versões

```
v0.1.0-beta.1  ← VOCÊ ESTÁ AQUI
v0.1.0-beta.2  (correções de bugs)
v0.2.0-beta.1  (novas features)
v0.3.0-beta.1  (mais features)
v1.0.0-rc.1    (release candidate)
v1.0.0         (primeira versão estável!)
```

## 🔗 Links Úteis

- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [Go Modules Reference](https://go.dev/ref/mod)

