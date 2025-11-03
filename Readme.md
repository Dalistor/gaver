# Gaver Framework

**Framework web para Go com CLI, geração de código e ORM**

[![Version](https://img.shields.io/badge/version-0.1.0--beta-orange.svg)](https://github.com/Dalistor/gaver/releases)
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

## Instalação

```bash
go install github.com/Dalistor/gaver/cmd/gaver@latest
```

Ou clone e compile:

```bash
git clone https://github.com/Dalistor/gaver.git
cd gaver
go build -o gaver cmd/gaver/main.go
```

## Quick Start

```bash
# Criar projeto
gaver init meu-projeto -d mysql
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
gaver init <nome> [-d database]       # Criar projeto
gaver serve [-p port]                  # Rodar servidor

# Modules
gaver module create <nome>             # Criar módulo
gaver module model <mod> <Model> [...] # Criar model
gaver module crud <mod> <Model>        # Gerar CRUD
  --only=list,get                      # Apenas métodos especificados
  --except=delete                      # Excluir métodos

# Migrations
gaver makemigrations [-n nome]         # Detectar mudanças
gaver migrate up                       # Aplicar migrations
gaver migrate down                     # Reverter migrations
gaver migrate status                   # Ver status
```

## Estrutura

```
meu-projeto/
├── cmd/server/         # Aplicação principal
├── config/            # Configurações
│   ├── routes/       # Registry de rotas
│   ├── modules/      # Registro de módulos
│   ├── database/     # Conexão com banco
│   └── ...
├── modules/          # Seus módulos
│   └── users/
│       ├── models/         # Models
│       ├── handlers/       # Controllers
│       ├── services/       # Lógica
│       ├── repositories/   # Dados
│       └── module.go       # Rotas
├── migrations/       # SQL migrations
└── .env
```

## Bancos Suportados

- MySQL
- PostgreSQL  
- SQLite

## Versão Atual

**v0.1.0-beta** - Primeira versão beta com core features

**Implementado:**
- Sistema de modules
- Geração de CRUD
- Annotations gaverModel  
- Migrations (makemigrations/migrate)
- Callbacks Before/After
- Registro automático de rotas

## Contribuindo

Contribuições são bem-vindas! Este projeto está em beta e feedback é essencial.

## Licença

MIT License - veja [LICENSE](LICENSE) para detalhes.

