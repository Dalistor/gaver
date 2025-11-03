# 📤 Guia de Publicação - Gaver Framework

## ✅ Checklist Pré-Publicação

Antes de publicar no GitHub, verifique:

- [x] `go.mod` com path correto: `github.com/Dalistor/gaver`
- [x] Versão do Go corrigida: `go 1.21`
- [x] LICENSE criada (MIT)
- [x] README.md completo
- [x] CHANGELOG.md atualizado
- [x] .gitignore configurado
- [x] VERSION definida: `0.1.0-beta`
- [ ] Código compila sem erros
- [ ] Sem erros de linter

## 🚀 Passo a Passo para Publicar

### 1. Verificar se Compila

```bash
# Compilar
go build -o gaver.exe cmd/gaver/main.go

# Testar
./gaver.exe --help
```

### 2. Inicializar Git (se ainda não fez)

```bash
# Inicializar repositório
git init

# Adicionar todos os arquivos
git add .

# Primeiro commit
git commit -m "feat: initial release v0.1.0-beta"
```

### 3. Criar Repositório no GitHub

1. Acesse: https://github.com/new
2. Nome do repositório: `gaver`
3. Descrição: `Framework web completo para Go com CLI, geração de código e ORM estilo Django`
4. Público
5. **NÃO** adicione README, LICENSE ou .gitignore (já temos)
6. Criar repositório

### 4. Conectar ao GitHub

```bash
# Adicionar remote
git remote add origin https://github.com/Dalistor/gaver.git

# Renomear branch para main (se necessário)
git branch -M main

# Push inicial
git push -u origin main
```

### 5. Criar Tag de Versão Beta

```bash
# Criar tag anotada
git tag -a v0.1.0-beta -m "v0.1.0-beta - Initial Beta Release

🎉 Primeira versão beta do Gaver Framework!

Core Functionality:
- Sistema de modules completo
- Annotations gaverModel para controle de campos
- CRUD automático com callbacks Before/After
- Migrations inteligentes (makemigrations/migrate)
- Suporte a MySQL, PostgreSQL, SQLite
- Framework HTTP com Gin
- Sistema de rotinas agendadas
- Validações automáticas

⚠️ Beta de longa duração (6-12 meses)
API pode ter breaking changes entre versões.
Não use em produção ainda.

Timeline previsto: v1.0.0 em Q2 2027
"

# Push da tag
git push origin v0.1.0-beta
```

### 6. Criar Release no GitHub

1. Vá para: https://github.com/Dalistor/gaver/releases/new
2. Choose a tag: `v0.1.0-beta`
3. Release title: `v0.1.0-beta - Initial Beta Release`
4. Descrição:

```markdown
## 🎉 Gaver Framework - Primeira Versão Beta!

**Framework web completo para Go** inspirado no Django, com CLI poderoso e geração de código inteligente.

### ⚠️ Beta de Longa Duração

Este projeto ficará em **beta por 6-12 meses**. A API pode sofrer mudanças significativas entre versões.
- 📅 Versão estável prevista: **v1.0.0 em Q2 2027**
- 🔄 Breaking changes permitidos entre versões beta
- 🧪 Use para desenvolvimento e testes, **não para produção**

### ✨ Funcionalidades v0.1.0-beta

- 🎯 **Sistema de Modules** - Organize código em módulos independentes
- 🔖 **Annotations gaverModel** - Controle validações e permissões por annotations
- 🔄 **CRUD Automático** - Gere handlers, services e repositories completos
- 📊 **Migrations Inteligentes** - `makemigrations` detecta mudanças automaticamente
- 🗄️ **Multi-Database** - MySQL, PostgreSQL, SQLite via GORM
- 🌐 **Gin Framework** - Performance e simplicidade HTTP
- ⚙️ **Rotinas Agendadas** - Sistema de cron jobs integrado
- 🎨 **Callbacks** - Before/After em todas operações CRUD
- ✅ **Validações Automáticas** - Baseadas em annotations

### 📦 Instalação

```bash
go install github.com/Dalistor/gaver/cmd/gaver@v0.1.0-beta
```

### 🚀 Quick Start

```bash
# Criar projeto
gaver init meu-projeto
cd meu-projeto

# Criar módulo
gaver module create users

# Criar model
gaver module:model users User name:string email:string

# Gerar CRUD
gaver module:crud users User

# Migrations
gaver makemigrations
gaver migrate up

# Rodar
go run cmd/server/main.go
```

### ⚠️ Importante - Beta de Longa Duração

Esta é uma **versão beta de longa duração** (~12-18 meses até v1.0.0).

**O que isso significa:**
- API pode ter breaking changes entre versões
- Novas features sendo adicionadas constantemente
- Bugs esperados e bem-vindos
- Feedback da comunidade molda o framework
- **NÃO use em produção ainda**

**Ideal para:**
- ✅ Projetos pessoais e aprendizado
- ✅ Protótipos e MVPs
- ✅ Desenvolvimento e experimentação
- ❌ Aplicações em produção

### 📖 Documentação

[Ver README completo](https://github.com/Dalistor/gaver#readme)

### 🐛 Encontrou um Bug?

[Reporte aqui](https://github.com/Dalistor/gaver/issues/new)
```

5. Marque: ✅ **Set as a pre-release**
6. Clique em **Publish release**

## 🎉 Pronto! Agora os usuários podem instalar:

```bash
go install github.com/Dalistor/gaver/cmd/gaver@v0.1.0-beta.1
```

## 📈 Próximas Versões

### v0.2.0-beta (QuerySet API) - Q1 2026

```bash
# 1. Implementar features do roadmap v0.2.0
git commit -m "feat: adiciona QuerySet API completo"
git commit -m "feat: validações cross-field"

# 2. Atualizar VERSION
echo "0.2.0-beta" > VERSION

# 3. Atualizar CHANGELOG.md
# Adicionar seção completa [0.2.0-beta]

# 4. Commit e tag
git add VERSION CHANGELOG.md
git commit -m "chore: release v0.2.0-beta

Major features:
- QuerySet API estilo Django
- Validações avançadas
- Breaking: Nova estrutura de validators
"
git tag -a v0.2.0-beta -m "v0.2.0-beta - QuerySet API"

# 5. Push
git push origin main v0.2.0-beta
```

### v0.3.0-beta (Developer Experience) - Q2 2026

```bash
# Implementar DX improvements
echo "0.3.0-beta" > VERSION
git commit -m "chore: release v0.3.0-beta"
git tag -a v0.3.0-beta -m "v0.3.0-beta - DX Improvements"
git push origin main v0.3.0-beta
```

### Hotfix Crítico (Raro)

Se houver bug crítico entre versões:
```bash
echo "0.1.1-beta" > VERSION
git tag -a v0.1.1-beta -m "Critical hotfix"
git push origin main v0.1.1-beta
```

## 🔐 Segurança

Para reportar vulnerabilidades de segurança, envie email para:
security@example.com (ou crie uma GitHub Security Advisory)

## 📊 Estatísticas

Depois de publicado, você pode ver:
- Downloads via `go install`
- Stars no GitHub
- Forks
- Issues/PRs

## 🎯 Metas para v1.0.0 (Versão Estável)

### Critérios Obrigatórios

**Qualidade:**
- [ ] 85%+ cobertura de testes
- [ ] Zero bugs críticos conhecidos
- [ ] Performance benchmarks publicados
- [ ] Security audit completo

**Documentação:**
- [ ] Guias completos para todas features
- [ ] Exemplos de projetos reais
- [ ] API reference completa
- [ ] Tutoriais em vídeo

**Estabilidade:**
- [ ] API estável por 3+ meses sem breaking changes
- [ ] 100+ projetos usando em desenvolvimento
- [ ] 50+ issues resolvidas
- [ ] Feedback positivo da comunidade

**Features Completas:**
- [ ] QuerySet API completo
- [ ] Sistema de auth integrado
- [ ] Admin interface
- [ ] CLI com todas features planejadas
- [ ] Migrations 100% funcionais
- [ ] Validações robustas

### Timeline Realista

- **Nov 2025 - Mar 2026:** Desenvolvimento ativo, breaking changes frequentes
- **Abr 2026 - Set 2026:** Estabilização, menos breaking changes
- **Out 2026 - Dez 2026:** Feature complete, apenas refinamentos
- **Jan 2027 - Mar 2027:** Feature freeze, bug fixes e docs
- **Abr 2027 - Jun 2027:** Release candidates
- **Jul 2027:** v1.0.0 (se tudo correr bem)

**Nota:** Preferimos lançar tarde e estável do que cedo e bugado!

---

**Qualquer dúvida?** Abra uma issue ou discussion no GitHub!

