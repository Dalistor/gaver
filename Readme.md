# Gaver Framework

**Framework web para Go com CLI, geração de código e ORM**

[![Version](https://img.shields.io/badge/version-0.1.1--beta-orange.svg)](https://github.com/Dalistor/gaver/releases)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

> ⚠️ **Beta:** Este projeto está em desenvolvimento ativo. A API pode mudar.

## Funcionalidades

- CLI completo com geração de código
- Sistema modular para organização
- Annotations `gaverModel` para controle de campos e validações
- CRUD automático com callbacks Before/After
- Migrations inteligentes (makemigrations/migrate)
- Suporte a MySQL, PostgreSQL, SQLite via GORM
- Framework HTTP com Gin
- Sistema de rotinas agendadas
- **Multi-plataforma**: Suporte para projetos Server, Android e Desktop
- **Frontend integrado**: Quasar Framework pré-configurado para Android (Capacitor) e Desktop (Electron)
- **Build automatizado**: Geração de APK (Android) e .exe (Desktop)

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

### Projeto Desktop

```bash
# Criar projeto Desktop
gaver init meu-app -d mysql -t desktop
cd meu-app
go mod tidy

# Instalar dependências do frontend
cd frontend
npm install

# Rodar servidor Go + Quasar dev (simultaneamente)
cd ..
gaver serve

# Gerar .exe
gaver build
```

### Projeto Web (SPA)

```bash
# Criar projeto Web
gaver init meu-app -d mysql -t web
cd meu-app
go mod tidy

# Instalar dependências do frontend
cd frontend
npm install

# Rodar servidor Go + Quasar dev (simultaneamente)
cd ..
gaver serve

# Gerar build estático
gaver build
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
gaver serve [--android]                   # Rodar servidor
                                                  # --android: abre Android Studio (apenas Android)
gaver build                                # Compilar projeto
                                                  # Android: gera AAR do Go e APK com Capacitor
                                                  # Desktop: gera binário Go e .exe com Electron
                                                  # Web: gera pasta build/ com binário Go e SPA
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

### SQLite em Projetos Android e Desktop

O SQLite é totalmente suportado em projetos Android e Desktop, utilizando o driver **modernc.org/sqlite** (puro Go, sem CGO). Isso significa:

- ✅ **Sem dependências externas**: O SQLite é embutido no executável
- ✅ **Funciona com CGO desabilitado**: Builds Android funcionam perfeitamente
- ✅ **Armazenamento persistente**: O banco é armazenado no diretório de dados do app

**Localização do arquivo SQLite:**
- **Android**: `getFilesDir()/data/<nome-do-banco>.db` (diretório de dados do app)
- **Desktop (Electron)**: `app.getPath('userData')/data/<nome-do-banco>.db` (diretório de dados do usuário)
- **Server/Web**: Diretório atual do projeto ou configurável via `APP_DATA_DIR`

O caminho é configurado automaticamente via variável de ambiente `APP_DATA_DIR` quando o app inicia.

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

### Desktop
Projeto completo com backend Go + frontend Quasar com Electron. Gera executável (.exe no Windows).
- Frontend pré-configurado com Quasar
- Router em modo history
- Cliente API base configurado
- Estrutura organizada para facilitar trabalho de IA no frontend
- Build gera binário Go e inclui no instalador Electron

### Web
Projeto completo com backend Go + frontend Quasar em modo SPA (Single Page Application). Gera build completo para deploy web.
- Frontend pré-configurado com Quasar (sem Capacitor/Electron)
- Router em modo history
- Cliente API base configurado
- Estrutura organizada para facilitar trabalho de IA no frontend
- Build gera pasta `build/` com binário Go e SPA prontos para deploy

## Frontend com Quasar

Projetos Android, Desktop e Web incluem Quasar Framework pré-configurado:

- **Proxy automático**: Frontend configurado para apontar ao servidor Go
- **Router history mode**: URLs sem # (hash)
- **Estrutura organizada**: Pastas separadas para composables, api, components, pages, layouts
- **Cliente API base**: Arquivo `client.js` pré-configurado para comunicação com backend
- **Composable useApi**: Composable Vue reutilizável para facilitar chamadas à API
- **Pronto para IA**: Estrutura pensada para facilitar trabalho de IA no desenvolvimento do frontend

### Fluxo de Trabalho

1. Dev executa `gaver init projeto -t android` (ou desktop)
2. Estrutura é criada com Quasar pré-configurado
3. Dev cria scripts de conexão com API em `frontend/src/api/`
4. Dev/IA trabalha em `frontend/src/components/` e `frontend/src/pages/`
5. Dev executa `gaver serve` (ou `gaver serve --android` para debug)
6. Quando finalizado, executa `gaver build` para gerar dist

## Versão Atual

**v0.1.1-beta** - Versão de testes

**Implementado:**
- Sistema de modules
- Geração de CRUD
- Annotations gaverModel  
- Migrations (makemigrations/migrate)
- Callbacks Before/After
- Registro automático de rotas
- **Multi-plataforma**: Projetos Server, Android e Desktop
- **Frontend integrado**: Quasar Framework com Capacitor/Electron
- **Build automatizado**: Geração de APK e .exe

## Contribuindo

Contribuições são bem-vindas! Este projeto está em beta e feedback é essencial.

## Licença

MIT License - veja [LICENSE](LICENSE) para detalhes.

