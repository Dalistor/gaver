# Gaver Framework

<div align="center">

**Framework multi-plataforma para Go com CLI, geração de código e ORM**

[![Version](https://img.shields.io/badge/version-1.1.1-FF6B35?style=flat-square)](https://github.com/Dalistor/gaver/releases)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-000000?style=flat-square)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Web%20%7C%20Desktop%20Windows-FF6B35?style=flat-square)](https://github.com/Dalistor/gaver)

</div>

---

## 🎯 Objetivo do Gaver

O **Gaver Framework** foi criado para encontrar um **meio termo ideal entre desenvolvimento humano e IA**:

- **Backend Go**: Desenvolvido principalmente pelo desenvolvedor
  - Estrutura modular
  - CLI para geração de código
  - Sistema de modules, CRUD automático, migrations
  - Annotations para controle de campos

- **Frontend**: Apenas composables de conexão com API feitos pelo dev
  - Estrutura pré-configurada com Quasar Framework
  - Cliente API base configurado
  - Composables prontos para uso

- **IA faz o resto**: Componentes, páginas, layouts e toda a interface
  - Estrutura organizada para facilitar trabalho de IA
  - Recomendado usar **Cursor** ou outra IDE com IA integrada
  - A IA trabalha em `frontend/src/components/` e `frontend/src/pages/`

**Fluxo de Trabalho:**
1. Dev cria backend Go com `gaver module` e `gaver crud`
2. Dev cria composables de API em `frontend/src/api/`
3. IA desenvolve componentes e páginas usando os composables
4. Resultado: Aplicação completa com backend robusto e frontend moderno

---

## 🚀 Quick Start

```bash
# Instalar
go install github.com/Dalistor/gaver/cmd/gaver@latest

# Criar projeto Web
gaver init meu-app -d sqlite -t web
cd meu-app
gaver serve

# Criar projeto Desktop (Windows)
gaver init meu-app -d sqlite -t desktop
cd meu-app
gaver serve
gaver build  # Gera .exe com servidor Go embutido
```

## ✨ Funcionalidades

### Backend (Desenvolvido pelo Dev)
- 🎯 **CLI completo** com geração de código
- 📦 **Sistema modular** para organização
- 🏷️ **Annotations `gaverModel`** para controle de campos e validações
- 🔄 **CRUD automático** com callbacks Before/After
- 📊 **Migrations inteligentes** (makemigrations/migrate)
- 🗄️ **Suporte a MySQL, PostgreSQL, SQLite** via GORM
- 🌐 **Framework HTTP** com Gin
- ⏰ **Sistema de rotinas** agendadas

### Frontend (Estrutura para IA)
- 🎨 **Quasar Framework** pré-configurado
- 🔌 **Composables de API** prontos para uso
- 📁 **Estrutura organizada** para facilitar trabalho de IA
- 🖥️ **Multi-plataforma**: Web e Desktop Windows
- 📦 **Build automatizado**: Gera executáveis prontos para distribuição

## 📚 Documentação

Para documentação completa, veja [DOCUMENTATION.md](DOCUMENTATION.md)

## 🎨 Plataformas Suportadas

| Plataforma | Status | Descrição |
|------------|--------|-----------|
| **Web** | ✅ Completo | SPA com Quasar Framework |
| **Desktop Windows** | ✅ Completo | Electron com servidor Go embutido |
| **Android** | ⚠️ Em desenvolvimento | Capacitor com servidor Go nativo |

## 📦 Instalação

```bash
go install github.com/Dalistor/gaver/cmd/gaver@latest
```

Ou clone e compile:

```bash
git clone https://github.com/Dalistor/gaver.git
cd gaver
go build -o gaver ./cmd/gaver
```

## 🛠️ Comandos Principais

### Criar Projeto

```bash
gaver init <nome-do-projeto> [flags]

Flags:
  -d, --database string    Tipo de banco de dados (padrão: "mysql")
                          Opções: mysql, postgres, sqlite
  
  -t, --type string       Tipo de projeto (padrão: "server")
                          Opções: server, web, desktop

Exemplos:
  # Projeto server com MySQL
  gaver init meu-api -d mysql -t server
  
  # Projeto web com SQLite
  gaver init meu-app -d sqlite -t web
  
  # Projeto desktop com PostgreSQL
  gaver init meu-app -d postgres -t desktop
```

### Desenvolvimento

```bash
gaver serve [flags]

Flags:
  --android    Abre Android Studio (apenas projetos Android)
  --cgo        Habilita CGO para SQLite (requer compilador C)

# Inicia servidor Go + frontend (se aplicável)
gaver serve
```

### Build

```bash
gaver build

# O comportamento varia conforme o tipo de projeto:
# - Server: Gera binário Go em bin/
# - Web: Gera pasta build/ com binário Go + SPA compilada
# - Desktop: Gera instalador .exe em frontend/dist/electron/
# - Android: Gera APK em frontend/android/app/build/outputs/apk/
```

### Modules

```bash
gaver module create <nome>
gaver module model <mod> <Model>
gaver module crud <mod> <Model> [flags]
  --only=list,get      # Apenas métodos especificados
  --except=delete     # Excluir métodos
```

### Migrations

```bash
gaver makemigrations [-n nome]
gaver migrate up
gaver migrate down
gaver migrate status
```

## 📖 Exemplo Rápido

```bash
# 1. Criar projeto
gaver init blog -d sqlite -t web

# 2. Criar módulo
gaver module create posts

# 3. Criar model
gaver module model posts Post

# 4. Editar modules/posts/models/post.go
# Adicionar campos com annotations:
# // gaverModel: writable:post,put; readable; required
# Title string `json:"title"`

# 5. Gerar CRUD
gaver module crud posts Post

# 6. Migrations
gaver makemigrations
gaver migrate up

# 7. Rodar
gaver serve
```

## 🤝 Contribuindo

Contribuições são bem-vindas! Este projeto está em desenvolvimento ativo.

## 📄 Licença

MIT License - veja [LICENSE](LICENSE) para detalhes.

---

<div align="center">

**Desenvolvido com ❤️ usando Go e Quasar**

[Documentação Completa](DOCUMENTATION.md) • [Issues](https://github.com/Dalistor/gaver/issues) • [Releases](https://github.com/Dalistor/gaver/releases)

</div>
