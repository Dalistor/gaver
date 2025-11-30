# Gaver Framework

**Framework multi-plataforma para Go com CLI, geração de código e ORM**

[![Version](https://img.shields.io/badge/version-1.1.0-blue.svg)](https://github.com/Dalistor/gaver/releases)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Web%20%7C%20Desktop%20Windows-lightgrey)](https://github.com/Dalistor/gaver)

> 🚀 **v1.1.0** - Suporte completo para Web e Desktop Windows com Electron

## Funcionalidades

- CLI completo com geração de código
- Sistema modular para organização
- Annotations `gaverModel` para controle de campos e validações
- CRUD automático com callbacks Before/After
- Migrations inteligentes (makemigrations/migrate)
- Suporte a MySQL, PostgreSQL, SQLite via GORM
- Framework HTTP com Gin
- Sistema de rotinas agendadas
- **Multi-plataforma**: Suporte para projetos Server, Web, Desktop (Windows) e Android
- **Frontend integrado**: Quasar Framework pré-configurado
  - **Web**: SPA (Single Page Application)
  - **Desktop**: Electron com servidor Go embutido
  - **Android**: Capacitor com servidor Go nativo
- **Build automatizado**: 
  - Web: Build estático + binário Go
  - Desktop: Executável .exe com servidor Go embutido
  - Android: APK com servidor Go nativo
- **SQLite sem CGO**: Driver puro Go para builds cross-platform

## Instalação

```bash
go install github.com/Dalistor/gaver/cmd/gaver@latest
```

Ou clone e compile:

```bash
git clone https://github.com/Dalistor/gaver.git
cd gaver
go build -o gaver ./cmd/gaver
# ou
go install ./cmd/gaver
```

## Quick Start

### Projeto Server (Backend apenas)

```bash
# Criar projeto server
gaver init meu-projeto -d mysql -t server
cd meu-projeto
go mod tidy

# Criar módulo
gaver module create users

# Criar model template
gaver module model users User

# Editar modules/users/models/user.go e adicionar seus campos

# Gerar CRUD (handlers, services, repositories + rotas)
gaver module crud users User

# Migrations
gaver makemigrations
gaver migrate up

# Rodar servidor
gaver serve
```

Servidor disponível em `http://localhost:8080`

### Projeto Android

```bash
# Criar projeto Android
gaver init meu-app -d mysql -t android
cd meu-app
go mod tidy

# Instalar dependências do frontend
cd frontend
npm install

# Rodar servidor Go + Quasar dev (simultaneamente)
cd ..
gaver serve

# Para abrir Android Studio para debug
gaver serve --android

# Gerar APK
gaver build
```

### Projeto Desktop (Windows)

```bash
# Criar projeto Desktop
gaver init meu-app -d sqlite -t desktop
cd meu-app
go mod tidy

# Instalar dependências do frontend (automático no init)
# npm install já é executado automaticamente

# Desenvolvimento: Rodar servidor Go + Quasar dev + Electron
gaver serve
# O comando inicia:
# 1. Servidor Go na porta 8080
# 2. Quasar dev server (faz proxy de /api para o Go)
# 3. Electron abre e carrega o frontend do Quasar dev server

# Build: Gerar executável .exe com servidor Go embutido
gaver build
# Gera:
# - frontend/dist/electron/ com executável .exe
# - Servidor Go compilado embutido no app
# - Ao abrir o app, o servidor Go inicia automaticamente
```

**Fluxo de Desenvolvimento Desktop:**
1. `gaver serve` inicia o servidor Go primeiro
2. Aguarda o servidor Go estar pronto
3. Inicia o Quasar dev server (porta padrão: 9000)
4. Electron abre e carrega o frontend do Quasar dev server
5. O Quasar dev server faz proxy de `/api` para o servidor Go

**Fluxo de Build Desktop:**
1. Compila o servidor Go para binário
2. Copia o binário para `frontend/src-electron/`
3. Build do Quasar Electron (gera arquivos estáticos)
4. Electron empacota tudo em um .exe
5. Ao abrir o app, o Electron inicia o servidor Go automaticamente

### Projeto Web (SPA)

```bash
# Criar projeto Web
gaver init meu-app -d mysql -t web
cd meu-app
go mod tidy

# Instalar dependências do frontend (automático no init)
# npm install já é executado automaticamente

# Desenvolvimento: Rodar servidor Go + Quasar dev
gaver serve
# O comando inicia:
# 1. Servidor Go na porta 8080
# 2. Quasar dev server (faz proxy de /api para o Go)
# 3. Abre navegador automaticamente

# Build: Gerar build estático para deploy
gaver build
# Gera:
# - build/ com binário Go e SPA compilada
# - Pronto para deploy em servidor web
```

### Rotas geradas automaticamente:

```
GET    /api/v1/users
GET    /api/v1/users/:id
POST   /api/v1/users
PUT    /api/v1/users/:id
PATCH  /api/v1/users/:id
DELETE /api/v1/users/:id
```

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

## Rotinas Agendadas

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

## Comandos

```bash
# Projeto
gaver init <nome> [-d database] [-t type]  # Criar projeto
                                                  # -t: server (padrão), android, desktop, web
gaver serve [--android] [--cgo]          # Rodar servidor
                                                  # --android: abre Android Studio (apenas Android)
                                                  # --cgo: habilita CGO para SQLite (requer compilador C)
gaver build                                # Compilar projeto
                                                  # Web: gera pasta build/ com binário Go e SPA
                                                  # Desktop: gera .exe com servidor Go embutido
                                                  # Android: gera APK com servidor Go nativo
                                                  # Server: build Go normal

# Modules
gaver module create <nome>                 # Criar módulo
gaver module model <mod> <Model> [...]    # Criar model
gaver module crud <mod> <Model>            # Gerar CRUD
  --only=list,get                          # Apenas métodos especificados
  --except=delete                          # Excluir métodos

# Migrations
gaver makemigrations [-n nome]             # Detectar mudanças
gaver migrate up                           # Aplicar migrations
gaver migrate down                         # Reverter migrations
gaver migrate status                       # Ver status
```

## Estrutura

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

### Projeto Android/Desktop/Web

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
│   │   ├── router/       # Configuração de rotas (history mode)
│   │   └── boot/         # Boot files do Quasar
│   ├── quasar.config.js  # Configuração Quasar (proxy para servidor Go)
│   ├── package.json
│   └── [capacitor.config.js (Android) ou electron/ (Desktop) ou nenhum (Web)]
├── android/               # Projeto Android nativo (apenas Android)
├── config/               # Configurações backend
├── modules/              # Módulos backend
└── migrations/           # SQL migrations
```

## Bancos Suportados

- MySQL
- PostgreSQL  
- SQLite

### SQLite em Projetos Web, Desktop e Android

O SQLite é totalmente suportado em todos os tipos de projeto, utilizando o driver **github.com/glebarez/sqlite** (puro Go, sem CGO). Isso significa:

- ✅ **Sem dependências externas**: O SQLite é embutido no executável
- ✅ **Funciona com CGO desabilitado**: Builds cross-platform funcionam perfeitamente
- ✅ **Sem compilador C necessário**: Funciona em qualquer ambiente Go
- ✅ **Armazenamento persistente**: O banco é armazenado no diretório apropriado para cada plataforma

**Localização do arquivo SQLite:**
- **Web/Server**: Diretório `data/` do projeto ou configurável via `APP_DATA_DIR`
- **Desktop (Electron)**: `app.getPath('userData')/data/<nome-do-banco>.db` (diretório de dados do usuário)
- **Android**: `getFilesDir()/data/<nome-do-banco>.db` (diretório de dados do app)

O caminho é configurado automaticamente via variável de ambiente `APP_DATA_DIR` quando o app inicia.

**Banco SQLite embutido:**
- No build Desktop/Android, se existir um arquivo `.db` no projeto, ele será copiado para o app
- Na primeira execução, o banco embutido é copiado para o diretório de dados do usuário
- Isso permite distribuir apps com banco pré-populado

## Tipos de Projeto

O Gaver suporta três tipos de projetos:

### Server
Projeto backend apenas, ideal para APIs REST. Estrutura mínima com servidor Go.

### Android
Projeto completo com backend Go + frontend Quasar com Capacitor. Gera APK para Android.
- Frontend pré-configurado com Quasar
- Router em modo history (sem # nas URLs)
- Cliente API base configurado para comunicação com backend
- Estrutura organizada para facilitar trabalho de IA no frontend
- Suporte a filesystem para armazenamento local
- Build gera AAR do Go e inclui no APK via Capacitor

### Desktop (Windows)
Projeto completo com backend Go + frontend Quasar com Electron. Gera executável (.exe) com servidor Go embutido.

**Características:**
- Frontend pré-configurado com Quasar Framework
- Router em modo history (sem # nas URLs)
- Cliente API base configurado para comunicação com backend
- Estrutura organizada para facilitar trabalho de IA no frontend
- **Servidor Go embutido**: O binário do servidor é incluído no .exe
- **Inicialização automática**: Ao abrir o app, o servidor Go inicia automaticamente
- **Modo dev**: No desenvolvimento, o Electron se conecta ao servidor Go já rodando via `gaver serve`
- **Modo produção**: O servidor Go é iniciado automaticamente pelo Electron

**Fluxo de Desenvolvimento:**
1. `gaver serve` inicia servidor Go (porta 8080)
2. Quasar dev server inicia (faz proxy de `/api` para o Go)
3. Electron abre e carrega frontend do Quasar dev server
4. Desenvolvimento com hot-reload

**Fluxo de Build:**
1. Compila servidor Go para binário
2. Copia binário para `frontend/src-electron/`
3. Build do Quasar Electron
4. Gera executável .exe com tudo embutido
5. Ao abrir o app, servidor Go inicia automaticamente

### Web (SPA)
Projeto completo com backend Go + frontend Quasar em modo SPA (Single Page Application). Gera build completo para deploy web.

**Características:**
- Frontend pré-configurado com Quasar Framework (sem Capacitor/Electron)
- Router em modo history (sem # nas URLs)
- Cliente API base configurado para comunicação com backend
- Estrutura organizada para facilitar trabalho de IA no frontend
- **Proxy automático**: No dev, Quasar faz proxy de `/api` para o servidor Go
- Build gera pasta `build/` com binário Go e SPA compilada prontos para deploy

**Fluxo de Desenvolvimento:**
1. `gaver serve` inicia servidor Go (porta 8080)
2. Quasar dev server inicia (faz proxy de `/api` para o Go)
3. Navegador abre automaticamente
4. Desenvolvimento com hot-reload

**Fluxo de Build:**
1. Compila servidor Go para binário
2. Build do Quasar SPA
3. Copia tudo para pasta `build/`
4. Pronto para deploy em servidor web

## Frontend com Quasar

Projetos Web, Desktop e Android incluem Quasar Framework pré-configurado:

- **Proxy automático**: No dev, Quasar faz proxy de `/api` para o servidor Go
- **Router history mode**: URLs sem # (hash)
- **Estrutura organizada**: Pastas separadas para composables, api, components, pages, layouts
- **Cliente API base**: Arquivo `client.js` pré-configurado para comunicação com backend
- **Composable useApi**: Composable Vue reutilizável para facilitar chamadas à API
- **Pronto para IA**: Estrutura pensada para facilitar trabalho de IA no desenvolvimento do frontend

### Fluxo de Trabalho

1. Dev executa `gaver init projeto -t web` (ou desktop)
2. Estrutura é criada com Quasar pré-configurado
3. `npm install` é executado automaticamente
4. Dev cria scripts de conexão com API em `frontend/src/api/`
5. Dev/IA trabalha em `frontend/src/components/` e `frontend/src/pages/`
6. Dev executa `gaver serve` para desenvolvimento
7. Quando finalizado, executa `gaver build` para gerar distribuição

### Desenvolvimento Desktop (Windows)

**Modo Dev (`gaver serve`):**
- Servidor Go inicia primeiro na porta 8080
- Quasar dev server inicia e faz proxy de `/api` para o Go
- Electron abre e carrega frontend do Quasar dev server
- Hot-reload funciona normalmente
- Electron **não** inicia outro servidor Go (usa o que já está rodando)

**Modo Build (`gaver build`):**
- Compila servidor Go para binário
- Copia binário para `frontend/src-electron/`
- Build do Quasar Electron (gera arquivos estáticos)
- Electron empacota tudo em executável .exe
- Ao abrir o app, Electron inicia o servidor Go automaticamente
- Frontend carrega dos arquivos estáticos (file://)

### Desenvolvimento Web (SPA)

**Modo Dev (`gaver serve`):**
- Servidor Go inicia primeiro na porta 8080
- Quasar dev server inicia e faz proxy de `/api` para o Go
- Navegador abre automaticamente
- Hot-reload funciona normalmente

**Modo Build (`gaver build`):**
- Compila servidor Go para binário
- Build do Quasar SPA (gera arquivos estáticos)
- Copia tudo para pasta `build/`
- Pronto para deploy em servidor web (Nginx, Apache, etc.)

## Versão Atual

**v1.1.0** - Suporte completo para Web e Desktop Windows

**Compatibilidade:**
- ✅ **Web**: SPA completa com Quasar Framework
- ✅ **Desktop Windows**: Electron com servidor Go embutido
- ✅ **Android**: Capacitor com servidor Go nativo (em desenvolvimento)

**Implementado:**
- Sistema de modules
- Geração de CRUD
- Annotations gaverModel  
- Migrations (makemigrations/migrate)
- Callbacks Before/After
- Registro automático de rotas
- **Multi-plataforma**: Projetos Server, Web, Desktop (Windows) e Android
- **Frontend integrado**: Quasar Framework
  - Web: SPA com proxy automático
  - Desktop: Electron com servidor Go embutido
  - Android: Capacitor com servidor Go nativo
- **Build automatizado**: 
  - Web: Build estático + binário Go
  - Desktop: Executável .exe com servidor embutido
  - Android: APK com servidor nativo
- **SQLite sem CGO**: Driver puro Go (github.com/glebarez/sqlite)
- **Inicialização automática**: Servidor Go inicia automaticamente em apps Desktop/Android
- **Modo dev otimizado**: Electron se conecta ao servidor já rodando

## Contribuindo

Contribuições são bem-vindas! Este projeto está em beta e feedback é essencial.

## Licença

MIT License - veja [LICENSE](LICENSE) para detalhes.

