# Documentação Completa - Gaver Framework

<div align="center">

**v1.1.0** - Framework multi-plataforma para Go

[![Version](https://img.shields.io/badge/version-1.1.0-FF6B35?style=flat-square)](https://github.com/Dalistor/gaver/releases)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-000000?style=flat-square)](LICENSE)

</div>

---

## 📋 Índice

- [Visão Geral](#visão-geral)
- [Instalação](#instalação)
- [Tipos de Projeto](#tipos-de-projeto)
- [Quick Start](#quick-start)
- [Sistema de Modules](#sistema-de-modules)
- [Annotations gaverModel](#annotations-gavermodel)
- [Callbacks](#callbacks)
- [Migrations](#migrations)
- [Rotinas Agendadas](#rotinas-agendadas)
- [Frontend com Quasar](#frontend-com-quasar)
- [Bancos de Dados](#bancos-de-dados)
- [Comandos CLI](#comandos-cli)
- [Estrutura de Projetos](#estrutura-de-projetos)

---

## Visão Geral

O **Gaver Framework** é um framework completo para desenvolvimento de aplicações Go com suporte multi-plataforma. Ele foi criado para encontrar um **meio termo ideal entre desenvolvimento humano e IA**.

### 🎯 Objetivo do Gaver

O Gaver estrutura o desenvolvimento de forma que:

**👨‍💻 Desenvolvedor trabalha em:**
- **Backend Go**: Desenvolvimento completo do backend
  - Criação de modules, models, handlers
  - Configuração de rotas e lógica de negócio
  - Definição de annotations `gaverModel`
  - Implementação de callbacks e validações
  - Criação de migrations

- **Frontend - Apenas Composables de API**: Apenas os composables de conexão com a API
  - Arquivos em `frontend/src/api/` para comunicação com backend
  - Composables em `frontend/src/composables/useApi.ts`
  - Configuração do cliente API base

**🤖 IA trabalha em:**
- **Frontend - Todo o resto**: Componentes, páginas, layouts e interface
  - Componentes Vue em `frontend/src/components/`
  - Páginas em `frontend/src/pages/`
  - Layouts em `frontend/src/layouts/`
  - Toda a interface do usuário

**💡 Recomendação:**
- Use **Cursor** ou outra IDE com IA integrada para desenvolvimento do frontend
- A estrutura do Gaver foi pensada para facilitar o trabalho de IA
- O backend robusto em Go fornece uma API estável para a IA trabalhar

### Características Principais

**Backend (Desenvolvido pelo Dev):**
- ✅ CLI completo com geração de código
- ✅ Sistema modular para organização
- ✅ Annotations `gaverModel` para controle de campos
- ✅ CRUD automático com callbacks
- ✅ Migrations inteligentes
- ✅ Suporte a MySQL, PostgreSQL, SQLite
- ✅ Framework HTTP com Gin
- ✅ Sistema de rotinas agendadas

**Frontend (Estrutura para IA):**
- ✅ Quasar Framework pré-configurado
- ✅ Composables de API prontos
- ✅ Estrutura organizada para facilitar trabalho de IA
- ✅ Multi-plataforma: Web e Desktop Windows
- ✅ Build automatizado

---

## Instalação

### Pré-requisitos

- Go 1.24+ instalado
- Node.js 18+ (para projetos Web/Desktop)
- npm ou yarn (para projetos Web/Desktop)

### Instalar Gaver CLI

```bash
go install github.com/Dalistor/gaver/cmd/gaver@latest
```

Verificar instalação:

```bash
gaver --version
```

### Instalar do código-fonte

```bash
git clone https://github.com/Dalistor/gaver.git
cd gaver
go build -o gaver ./cmd/gaver
# ou
go install ./cmd/gaver
```

---

## Tipos de Projeto

O Gaver suporta diferentes tipos de projetos:

### 🖥️ Server

Projeto backend apenas, ideal para APIs REST.

```bash
gaver init meu-api -d mysql -t server
```

**Características:**
- Estrutura mínima com servidor Go
- Sem frontend
- Ideal para microserviços e APIs

### 🌐 Web (SPA)

Projeto completo com backend Go + frontend Quasar em modo SPA.

```bash
gaver init meu-app -d sqlite -t web
```

**Características:**
- Frontend Quasar Framework pré-configurado
- Router em modo hash (compatível com file://)
- Cliente API base configurado
- Build gera pasta `build/` com binário Go e SPA compilada
- Pronto para deploy em servidor web

**Fluxo de Desenvolvimento:**
1. `gaver serve` inicia servidor Go (porta 8080)
2. Quasar dev server inicia (faz proxy de `/api` para o Go)
3. Navegador abre automaticamente
4. Hot-reload funciona normalmente

**Fluxo de Build:**
1. Compila servidor Go para binário
2. Build do Quasar SPA
3. Copia tudo para pasta `build/`
4. Pronto para deploy

### 🖥️ Desktop (Windows)

Projeto completo com backend Go + frontend Quasar com Electron.

```bash
gaver init meu-app -d sqlite -t desktop
```

**Características:**
- Frontend Quasar Framework pré-configurado
- Router em modo hash (compatível com file://)
- Cliente API base configurado
- **Servidor Go embutido**: Binário incluído no .exe
- **Inicialização automática**: Servidor Go inicia ao abrir o app
- **Instalador NSIS**: Gera instalador .exe profissional
- **Sem barra de menu**: Interface limpa em produção

**Fluxo de Desenvolvimento:**
1. `gaver serve` inicia servidor Go (porta 8080)
2. Quasar dev server inicia (faz proxy de `/api` para o Go)
3. Electron abre e carrega frontend do Quasar dev server
4. Hot-reload funciona normalmente
5. Electron **não** inicia outro servidor Go (usa o que já está rodando)

**Fluxo de Build:**
1. Compila servidor Go para binário
2. Copia binário, database.db e .env para `frontend/src-electron/`
3. Build do Quasar Electron (gera arquivos estáticos)
4. Electron-builder empacota tudo em instalador .exe
5. Arquivos (server.exe, database.db, .env) são incluídos via `extraResources`
6. Ao abrir o app, Electron inicia servidor Go automaticamente
7. Frontend carrega dos arquivos estáticos (file://)

**Ícones:**
- Logo padrão copiado automaticamente de `assets/logo.png`
- Ícones gerados automaticamente usando `@quasar/icongenie`
- Ícones ficam em `src-electron/icons/icon.ico`

### 📱 Android

Projeto completo com backend Go + frontend Quasar com Capacitor.

```bash
gaver init meu-app -d sqlite -t android
```

**Status:** ⚠️ Em desenvolvimento

**Características:**
- Frontend Quasar Framework pré-configurado
- Capacitor para desenvolvimento Android
- Servidor Go nativo via AAR
- Build gera APK

---

## Quick Start

### Projeto Server

```bash
# Criar projeto
gaver init meu-projeto -d mysql -t server
cd meu-projeto
go mod tidy

# Criar módulo
gaver module create users

# Criar model
gaver module model users User

# Editar modules/users/models/user.go e adicionar campos

# Gerar CRUD
gaver module crud users User

# Migrations
gaver makemigrations
gaver migrate up

# Rodar servidor
gaver serve
```

Servidor disponível em `http://localhost:8080`

### Projeto Web

```bash
# Criar projeto
gaver init meu-app -d sqlite -t web
cd meu-app
go mod tidy

# Desenvolvimento
gaver serve
# Inicia servidor Go + Quasar dev server

# Build
gaver build
# Gera pasta build/ com binário Go e SPA compilada
```

### Projeto Desktop (Windows)

```bash
# Criar projeto
gaver init meu-app -d sqlite -t desktop
cd meu-app
go mod tidy

# Desenvolvimento
gaver serve
# Inicia servidor Go + Quasar dev + Electron

# Build
gaver build
# Gera instalador .exe em frontend/dist/electron/
```

---

## Sistema de Modules

Modules são unidades organizacionais que agrupam funcionalidades relacionadas.

### Criar Module

```bash
gaver module create users
```

Isso cria a estrutura:

```
modules/users/
├── models/
├── handlers/
├── services/
├── repositories/
└── module.go
```

### Criar Model

```bash
gaver module model users User
```

Isso cria `modules/users/models/user.go` com template básico.

### Gerar CRUD

```bash
gaver module crud users User
```

Isso gera:
- Handlers (controllers)
- Services (lógica de negócio)
- Repositories (acesso a dados)
- Rotas automáticas

**Opções:**
```bash
# Apenas métodos específicos
gaver module crud users User --only=list,get

# Excluir métodos
gaver module crud users User --except=delete
```

### Rotas Geradas

```
GET    /api/v1/users
GET    /api/v1/users/:id
POST   /api/v1/users
PUT    /api/v1/users/:id
PATCH  /api/v1/users/:id
DELETE /api/v1/users/:id
```

---

## Annotations gaverModel

Controle de campos via annotations em comentários:

```go
type Product struct {
    // Controle de acesso
    // gaverModel: writable:post,put; readable; required
    Title string `json:"title"`
    
    // Validações
    // gaverModel: writable:post,put; readable; required; min:0; max:99999
    Price float64 `json:"price"`
    
    // Campos apenas leitura
    // gaverModel: ignore:write; readable
    ViewCount int `json:"view_count"`
    
    // Campos internos (não expostos na API)
    // gaverModel: ignore
    InternalCode string `json:"-"`
    
    // Relacionamentos
    // gaverModel: relation:belongsTo; foreignKey:category_id
    CategoryID uint     `json:"category_id"`
    Category   Category `json:"category" gorm:"foreignKey:CategoryID"`
}
```

### Tags Disponíveis

| Tag | Descrição | Exemplo |
|-----|-----------|---------|
| `writable:methods` | Métodos HTTP que podem escrever | `writable:post,put,patch` |
| `readable` | Pode ser lido em GET | `readable` |
| `required` | Campo obrigatório | `required` |
| `unique` | Valor único no banco | `unique` |
| `email` | Valida formato email | `email` |
| `min:N` / `max:N` | Valores numéricos | `min:18; max:120` |
| `minLength:N` / `maxLength:N` | Tamanho strings | `minLength:3; maxLength:100` |
| `enum:vals` | Valores permitidos | `enum:active,inactive,pending` |
| `relation:type` | Tipo de relacionamento | `relation:hasMany` |
| `ignore` | Ignorar campo completamente | `ignore` |
| `ignore:write` | Ignorar apenas em escrita | `ignore:write` |

---

## Callbacks

Personalize o comportamento do CRUD:

```go
// modules/users/handlers/user_handler.go

// Hash de senha antes de criar
func (h *UserHandler) BeforeCreate(c *gin.Context, data map[string]interface{}) error {
    if password, ok := data["password"].(string); ok {
        hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
        data["password"] = string(hashed)
    }
    return nil
}

// Remover senha antes de retornar
func (h *UserHandler) AfterGet(c *gin.Context, user models.User) models.User {
    user.Password = ""
    return user
}

// Validações customizadas
func (h *UserHandler) OnValidate(data map[string]interface{}, operation string) error {
    if age, ok := data["age"].(float64); ok {
        if age < 18 {
            return fmt.Errorf("usuário deve ter 18+ anos")
        }
    }
    return nil
}
```

### Callbacks Disponíveis

- `BeforeCreate` - Antes de criar
- `AfterCreate` - Depois de criar
- `BeforeUpdate` - Antes de atualizar
- `AfterUpdate` - Depois de atualizar
- `BeforeDelete` - Antes de deletar
- `AfterDelete` - Depois de deletar
- `BeforeGet` - Antes de buscar
- `AfterGet` - Depois de buscar
- `BeforeList` - Antes de listar
- `AfterList` - Depois de listar
- `OnValidate` - Validação customizada

---

## Migrations

Sistema de migrations inteligente que detecta mudanças automaticamente.

### Criar Migration

```bash
gaver makemigrations [-n nome_da_migration]
```

Detecta mudanças nos models e cria arquivo de migration.

### Aplicar Migrations

```bash
gaver migrate up
```

Aplica todas as migrations pendentes.

### Reverter Migration

```bash
gaver migrate down
```

Reverte a última migration aplicada.

### Status

```bash
gaver migrate status
```

Mostra status de todas as migrations.

---

## Rotinas Agendadas

Sistema de tarefas em background que rodam automaticamente.

```go
// config/routines/routines.go

func (m *Manager) RegisterDefaultRoutines() {
    // Limpar dados antigos diariamente
    m.Register("cleanup", 24*time.Hour, func() {
        log.Println("Limpando dados antigos...")
        // Seu código aqui
    })
    
    // Enviar emails a cada 5 minutos
    m.Register("emails", 5*time.Minute, func() {
        log.Println("Enviando emails pendentes...")
        // Seu código aqui
    })
}
```

### Casos de Uso

- ✅ Limpeza de dados antigos
- ✅ Sincronização com APIs externas
- ✅ Envio de emails/notificações agendadas
- ✅ Geração de relatórios
- ✅ Processamento de filas
- ✅ Health checks do sistema
- ✅ Backup de dados

---

## Frontend com Quasar

Projetos Web e Desktop incluem Quasar Framework pré-configurado com estrutura otimizada para trabalho com IA.

### 🎯 Divisão de Trabalho

**👨‍💻 Desenvolvedor cria:**
- **Composables de API**: Arquivos em `frontend/src/api/` para comunicação com backend
- **Cliente API base**: Configuração do axios em `frontend/src/api/client.js`
- **Composable useApi**: Composable reutilizável em `frontend/src/composables/useApi.ts`

**🤖 IA desenvolve:**
- **Componentes**: Tudo em `frontend/src/components/`
- **Páginas**: Tudo em `frontend/src/pages/`
- **Layouts**: Tudo em `frontend/src/layouts/`
- **Interface completa**: Usando os composables criados pelo dev

### Estrutura

```
frontend/src/
├── composables/     # Composables Vue reutilizáveis (DEV cria)
│   └── useApi.ts    # Composable base para API (DEV cria)
├── api/             # Arquivos JS/TS para comunicação com API (DEV cria)
│   └── client.js    # Cliente API base (DEV cria)
├── components/      # Componentes Vue (IA desenvolve)
├── pages/           # Páginas/views (IA desenvolve)
├── layouts/         # Layouts (IA desenvolve)
├── router/          # Configuração de rotas (pré-configurado)
└── boot/            # Boot files do Quasar (pré-configurado)
```

### Cliente API (Criado pelo Dev)

Arquivo `frontend/src/api/client.js` pré-configurado:

```javascript
import axios from 'axios'

const client = axios.create({
  baseURL: '/api',  // Proxy para servidor Go
  timeout: 10000
})

export default client
```

**O dev pode criar arquivos específicos por módulo:**
```javascript
// frontend/src/api/users.js
import client from './client'

export const usersApi = {
  list: () => client.get('/v1/users'),
  get: (id) => client.get(`/v1/users/${id}`),
  create: (data) => client.post('/v1/users', data),
  update: (id, data) => client.put(`/v1/users/${id}`, data),
  delete: (id) => client.delete(`/v1/users/${id}`)
}
```

### Composable useApi (Criado pelo Dev)

Composable reutilizável em `frontend/src/composables/useApi.ts`:

```typescript
import { ref } from 'vue'
import client from 'src/api/client'

export function useApi() {
  const loading = ref(false)
  const error = ref(null)

  const request = async (method: string, url: string, data?: any) => {
    loading.value = true
    error.value = null
    try {
      const response = await client[method](url, data)
      return response.data
    } catch (err: any) {
      error.value = err.response?.data || err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    get: (url: string) => request('get', url),
    post: (url: string, data: any) => request('post', url, data),
    put: (url: string, data: any) => request('put', url, data),
    del: (url: string) => request('delete', url),
    loading,
    error
  }
}
```

**Exemplo de uso (IA desenvolve):**
```vue
<script setup>
import { useApi } from 'src/composables/useApi'

const { get, post, loading, error } = useApi()

const fetchUsers = async () => {
  try {
    const users = await get('/v1/users')
    console.log(users)
  } catch (err) {
    console.error(err)
  }
}

const createUser = async (userData) => {
  try {
    const user = await post('/v1/users', userData)
    console.log('Usuário criado:', user)
  } catch (err) {
    console.error(err)
  }
}
</script>
```

### Router Mode

- **Web**: Modo `hash` (compatível com file:// e servidores web)
- **Desktop**: Modo `hash` (compatível com file:// no Electron)

O modo `hash` usa URLs com `#` (ex: `file:///path/index.html#/users`), que funcionam tanto com `file://` quanto com servidores web.

### Trabalho com IA

**Recomendado:**
- Use **Cursor** ou outra IDE com IA integrada
- A estrutura do Gaver facilita o trabalho de IA
- A IA trabalha apenas em componentes, páginas e layouts
- O backend Go fornece API estável e documentada

---

## Bancos de Dados

### Suportados

O Gaver suporta três bancos de dados:

| Banco | Tipo | Uso Recomendado |
|-------|------|-----------------|
| **MySQL** | Servidor externo | Produção, aplicações web com múltiplos usuários |
| **PostgreSQL** | Servidor externo | Produção, aplicações complexas |
| **SQLite** | Arquivo local | Desktop, desenvolvimento, aplicações single-user |

### MySQL e PostgreSQL

**Características:**
- Conexão com servidor de banco de dados externo
- Requer servidor MySQL/PostgreSQL rodando
- Ideal para produção e aplicações multi-usuário
- Suporta transações complexas e relacionamentos avançados

**Configuração (.env):**
```env
# MySQL
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=sua_senha
DB_NAME=meu_banco

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=sua_senha
DB_NAME=meu_banco
```

**Uso:**
- Configure o servidor de banco de dados separadamente
- O Gaver se conecta ao servidor via GORM
- Funciona igual em todos os tipos de projeto (Server, Web, Desktop)

### SQLite (Especial)

**Características:**
- ✅ **Driver puro Go**: Utiliza `github.com/glebarez/sqlite` (sem CGO)
- ✅ **Sem dependências externas**: SQLite embutido no executável
- ✅ **Funciona com CGO desabilitado**: Builds cross-platform funcionam perfeitamente
- ✅ **Sem compilador C necessário**: Funciona em qualquer ambiente Go
- ✅ **Integração na compilação**: Banco pode ser embutido no app (Desktop/Android)

**Localização do Banco:**
- **Web/Server**: Diretório `data/` do projeto ou configurável via `APP_DATA_DIR`
- **Desktop (Electron)**: `app.getPath('userData')/data/<nome-do-banco>.db`
- **Android**: `getFilesDir()/data/<nome-do-banco>.db`

**Banco SQLite Embutido (Desktop/Android):**
- Se existir um arquivo `.db` no projeto, ele será copiado para o app durante o build
- Na primeira execução, o banco embutido é copiado para o diretório de dados do usuário
- Permite distribuir apps com banco pré-populado
- O arquivo `.db` é incluído no instalador .exe via `extraResources`

**Configuração (.env):**
```env
# SQLite
DB_PATH=data/meu_banco.db
# ou use APP_DATA_DIR para Desktop/Android (configurado automaticamente)
```

**Diferença Principal:**
- **MySQL/PostgreSQL**: Requer servidor externo, conexão via rede
- **SQLite**: Arquivo local, pode ser embutido no executável (Desktop/Android)

---

## Comandos CLI

### Criar Projeto

```bash
gaver init <nome-do-projeto> [flags]
```

**Parâmetros:**
- `<nome-do-projeto>`: Nome do projeto (obrigatório)

**Flags:**
- `-d, --database string`: Tipo de banco de dados
  - **Padrão**: `mysql`
  - **Opções**: `mysql`, `postgres`, `sqlite`
  
- `-t, --type string`: Tipo de projeto
  - **Padrão**: `server`
  - **Opções**: `server`, `web`, `desktop`, `android`

**Exemplos:**
```bash
# Projeto server com MySQL (padrão)
gaver init meu-api

# Projeto server com PostgreSQL
gaver init meu-api -d postgres -t server

# Projeto web com SQLite
gaver init meu-app -d sqlite -t web

# Projeto desktop com SQLite
gaver init meu-app -d sqlite -t desktop

# Projeto desktop com MySQL
gaver init meu-app -d mysql -t desktop
```

### Desenvolvimento

```bash
gaver serve [flags]
```

**Flags:**
- `--android`: Abre Android Studio automaticamente (apenas projetos Android)
- `--cgo`: Habilita CGO para SQLite (requer compilador C instalado)

**Comportamento por tipo de projeto:**
- **Server**: Inicia apenas servidor Go na porta 8080
- **Web**: Inicia servidor Go + Quasar dev server (proxy automático)
- **Desktop**: Inicia servidor Go + Quasar dev server + Electron
- **Android**: Inicia servidor Go + Quasar dev server (use `--android` para abrir Android Studio)

**Exemplos:**
```bash
# Desenvolvimento normal
gaver serve

# Android com Android Studio
gaver serve --android

# SQLite com CGO (se necessário)
gaver serve --cgo
```

### Build

```bash
gaver build
```

**Comportamento varia conforme o tipo de projeto:**

#### Server
- Compila servidor Go
- Gera binário em `bin/` ou diretório raiz
- Apenas backend, sem frontend

#### Web
1. Compila servidor Go para binário
2. Build do Quasar SPA (gera arquivos estáticos)
3. Copia binário Go e SPA compilada para pasta `build/`
4. Resultado: Pasta `build/` pronta para deploy em servidor web

**Arquivos gerados:**
- `build/server` (ou `build/server.exe` no Windows)
- `build/dist/` (SPA compilada)
- `build/.env` (se existir)

#### Desktop (Windows)
1. Compila servidor Go para binário (`server.exe`)
2. Copia `server.exe`, `database.db` (se existir) e `.env` para `frontend/src-electron/`
3. Build do Quasar Electron (gera arquivos estáticos)
4. Electron-builder empacota tudo em instalador .exe
5. Arquivos (`server.exe`, `database.db`, `.env`) são incluídos via `extraResources`
6. Gera instalador NSIS em `frontend/dist/electron/Packaged/`

**Arquivos gerados:**
- `frontend/dist/electron/Packaged/<nome> Setup <versão>.exe` (instalador)
- `frontend/dist/electron/Packaged/win-unpacked/` (app não empacotado)
- Arquivos em `resources/`: `server.exe`, `database.db`, `.env`

**Características do build Desktop:**
- Instalador NSIS profissional
- Servidor Go embutido e inicia automaticamente
- Banco SQLite embutido (se existir)
- Arquivo .env incluído
- Menu removido em produção
- Ícones gerados automaticamente

#### Android
1. Compila servidor Go para AAR (biblioteca Android)
2. Copia AAR para `android/app/libs/`
3. Copia `database.db` (se existir) para `android/app/src/main/assets/`
4. Build do Quasar Capacitor
5. Build do Android (via Gradle)
6. Gera APK em `android/app/build/outputs/apk/`

**Arquivos gerados:**
- `android/app/build/outputs/apk/debug/app-debug.apk`
- `android/app/build/outputs/apk/release/app-release.apk`

### Modules

```bash
gaver module create <nome>
# Criar módulo

gaver module model <mod> <Model> [...]
# Criar model

gaver module crud <mod> <Model>
# Gerar CRUD
  --only=list,get      # Apenas métodos especificados
  --except=delete     # Excluir métodos
```

### Migrations

```bash
gaver makemigrations [-n nome]
# Detectar mudanças

gaver migrate up
# Aplicar migrations

gaver migrate down
# Reverter migrations

gaver migrate status
# Ver status
```

---

## Estrutura de Projetos

### Projeto Server

```
meu-projeto/
├── GaverProject.json      # Configuração do projeto
├── cmd/server/            # Aplicação principal
├── config/                # Configurações
│   ├── routes/           # Registry de rotas
│   ├── modules/          # Registro de módulos
│   ├── database/         # Conexão com banco
│   └── ...
├── modules/              # Seus módulos
│   └── users/
│       ├── models/       # Models
│       ├── handlers/     # Controllers
│       ├── services/     # Lógica
│       ├── repositories/ # Dados
│       └── module.go     # Rotas
├── migrations/           # SQL migrations
└── .env
```

### Projeto Web/Desktop

```
meu-projeto/
├── GaverProject.json      # Configuração do projeto
├── cmd/server/            # Servidor Go (backend)
├── frontend/              # Aplicação Quasar
│   ├── src/
│   │   ├── composables/  # Composables Vue reutilizáveis
│   │   │   └── useApi.ts # Composable base para API
│   │   ├── api/          # Arquivos JS/TS para comunicação com API
│   │   │   └── client.js # Cliente API base
│   │   ├── components/   # Componentes Vue
│   │   ├── pages/        # Páginas/views
│   │   ├── layouts/      # Layouts
│   │   ├── router/       # Configuração de rotas (hash mode)
│   │   └── boot/         # Boot files do Quasar
│   ├── quasar.config.js  # Configuração Quasar (proxy para servidor Go)
│   ├── package.json
│   └── [src-electron/ (Desktop) ou nenhum (Web)]
├── config/               # Configurações backend
├── modules/              # Módulos backend
└── migrations/           # SQL migrations
```

---

## Versão Atual

**v1.1.0** - Suporte completo para Web e Desktop Windows

### Compatibilidade

- ✅ **Web**: SPA completa com Quasar Framework
- ✅ **Desktop Windows**: Electron com servidor Go embutido
- ⚠️ **Android**: Capacitor com servidor Go nativo (em desenvolvimento)

### Implementado

- Sistema de modules
- Geração de CRUD
- Annotations gaverModel
- Migrations (makemigrations/migrate)
- Callbacks Before/After
- Registro automático de rotas
- Multi-plataforma: Projetos Server, Web, Desktop (Windows)
- Frontend integrado: Quasar Framework
  - Web: SPA com proxy automático
  - Desktop: Electron com servidor Go embutido
- Build automatizado:
  - Web: Build estático + binário Go
  - Desktop: Instalador .exe com servidor embutido
- SQLite sem CGO: Driver puro Go (github.com/glebarez/sqlite)
- Inicialização automática: Servidor Go inicia automaticamente em apps Desktop
- Modo dev otimizado: Electron se conecta ao servidor já rodando
- Router hash mode: Compatível com file:// no Electron
- Menu removido: Interface limpa em produção

---

## Contribuindo

Contribuições são bem-vindas! Este projeto está em desenvolvimento ativo.

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/MinhaFeature`)
3. Commit suas mudanças (`git commit -m 'Adiciona MinhaFeature'`)
4. Push para a branch (`git push origin feature/MinhaFeature`)
5. Abra um Pull Request

---

## Licença

MIT License - veja [LICENSE](LICENSE) para detalhes.

---

<div align="center">

**Desenvolvido com ❤️ usando Go e Quasar**

[Voltar ao README](Readme.md)

</div>

