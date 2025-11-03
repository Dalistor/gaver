# 📤 Guia de Publicação - Gaver Framework

## ✅ Checklist Pré-Publicação

Antes de publicar no GitHub, verifique:

- [x] `go.mod` com path correto: `github.com/Dalistor/gaver`
- [x] Versão do Go corrigida: `go 1.21`
- [x] LICENSE criada (MIT)
- [x] README.md completo
- [x] CHANGELOG.md atualizado
- [x] .gitignore configurado
- [x] VERSION definida: `0.1.0-beta.1`
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
git commit -m "feat: initial commit - gaver framework v0.1.0-beta.1"
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
git tag -a v0.1.0-beta.1 -m "First beta release

🎉 Primeira versão beta do Gaver Framework!

Funcionalidades:
- Sistema de modules
- Annotations gaverModel
- CRUD automático com callbacks
- Migrations inteligentes
- Suporte a MySQL, PostgreSQL, SQLite
- Framework HTTP com Gin
- Sistema de rotinas agendadas

⚠️ API pode mudar durante o beta testing.
"

# Push da tag
git push origin v0.1.0-beta.1
```

### 6. Criar Release no GitHub

1. Vá para: https://github.com/Dalistor/gaver/releases/new
2. Choose a tag: `v0.1.0-beta.1`
3. Release title: `v0.1.0-beta.1 - First Beta Release`
4. Descrição:

```markdown
## 🎉 Primeira versão beta do Gaver Framework!

### ✨ Funcionalidades

- 🎯 **Sistema de Modules** - Organize código em módulos reutilizáveis
- 🔖 **Annotations gaverModel** - Controle de validações e permissões
- 🔄 **CRUD automático** - Gere handlers, services e repositories
- 📊 **Migrations inteligentes** - Detecta mudanças automaticamente
- 🗄️ **Multi-database** - MySQL, PostgreSQL, SQLite
- 🌐 **Gin Framework** - Performance e simplicidade
- ⚙️ **Rotinas agendadas** - Tarefas em background

### 📦 Instalação

```bash
go install github.com/Dalistor/gaver/cmd/gaver@v0.1.0-beta.1
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

### ⚠️ Aviso Importante

Esta é uma **versão beta**. A API pode sofrer mudanças até a versão 1.0.0.
Use para testes e envie feedback!

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

### v0.1.0-beta.2 (Correções)

```bash
# Fazer correções
git commit -m "fix: corrige bug X"

# Atualizar VERSION
echo "0.1.0-beta.2" > VERSION

# Atualizar CHANGELOG.md
# Adicionar seção [0.1.0-beta.2]

# Commit e tag
git add VERSION CHANGELOG.md
git commit -m "chore: bump version to v0.1.0-beta.2"
git tag -a v0.1.0-beta.2 -m "Bug fixes"
git push origin main v0.1.0-beta.2
```

### v0.2.0-beta.1 (Nova Feature)

```bash
# Implementar feature
git commit -m "feat: adiciona suporte a WebSockets"

# Atualizar VERSION
echo "0.2.0-beta.1" > VERSION

# Commit e tag
git commit -m "chore: bump version to v0.2.0-beta.1"
git tag -a v0.2.0-beta.1 -m "Add WebSocket support"
git push origin main v0.2.0-beta.1
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

## 🎯 Metas para v1.0.0

- [ ] 80%+ cobertura de testes
- [ ] Documentação completa
- [ ] 10+ usuários beta testando
- [ ] API estável (sem breaking changes)
- [ ] Performance otimizada
- [ ] Exemplos completos

---

**Qualquer dúvida?** Abra uma issue ou discussion no GitHub!

