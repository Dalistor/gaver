# Gaver Framework

> Framework web para Go com CLI, geração de código e ORM estilo Django

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-0.1.0--beta-orange.svg)](https://github.com/Dalistor/gaver/releases)

**Versão Atual:** `v0.1.0-beta` (Beta Testing)

⚠️ **Projeto em fase beta.** API pode sofrer alterações. Não recomendado para produção.

## Instalação

```bash
go install github.com/Dalistor/gaver/cmd/gaver@latest
```

## Início Rápido

```bash
# Criar projeto
gaver init meu-app -d mysql
cd meu-app
go mod tidy

# Criar módulo
gaver module create users

# Criar model
gaver module model users User name:string email:string:unique age:int

# Gerar CRUD completo
gaver module crud users User

# Migrations
gaver makemigrations
gaver migrate up

# Rodar servidor
gaver serve
```

Servidor em: `http://localhost:8080`

## Funcionalidades

- 🎯 **CLI completo** - Geração automática de código
- 📦 **Sistema de Modules** - Organize por domínios
- 🔖 **Annotations gaverModel** - Controle validações e permissões
- 🔄 **CRUD automático** - Handlers, services e repositories
- 📊 **Migrations inteligentes** - Detecta mudanças automaticamente
- 🗄️ **Multi-database** - MySQL, PostgreSQL, SQLite
- 🌐 **Gin Framework** - HTTP rápido e simples
- ⚙️ **Rotinas agendadas** - Tarefas em background

## Comandos

### Projeto
```bash
gaver init <nome> [-d mysql|postgres|sqlite]  # Criar projeto
gaver serve [-p porta]                        # Rodar servidor
```

### Modules
```bash
gaver module create <nome>                    # Criar módulo
gaver module model <module> <Model> [campos]  # Criar model
gaver module crud <module> <Model>            # Gerar CRUD
  --only=list,get                            # Apenas métodos específicos
  --except=delete                            # Excluir métodos
```

### Migrations
```bash
gaver makemigrations [-n nome]     # Detectar mudanças
gaver migrate up                   # Aplicar migrations
gaver migrate down                 # Reverter migrations
gaver migrate status               # Ver status
```

## Annotations gaverModel

Controle campos do model com annotations:

```go
type User struct {
    // gaverModel: primaryKey; autoIncrement
    ID uint `json:"id" gorm:"primaryKey"`
    
    // gaverModel: writable:post,put; readable; required; minLength:3
    Name string `json:"name"`
    
    // gaverModel: writable:post; readable; required; unique; email
    Email string `json:"email"`
    
    // gaverModel: writable:post,put,patch; readable; min:18; max:120
    Age int `json:"age"`
    
    // gaverModel: ignore:write; readable
    CreatedAt time.Time `json:"created_at"`
    
    // gaverModel: ignore
    Password string `json:"-"`
}
```

### Tags Disponíveis

**Controle de Acesso:**
- `writable:post,put,patch` - Métodos que podem escrever
- `readable` - Pode ser lido
- `ignore:write` ou `ignore:read` - Ignorar escrita/leitura
- `ignore` - Completamente ignorado

**Validações:**
- `required` - Obrigatório
- `unique` - Valor único
- `email`, `url` - Formato específico
- `min:N`, `max:N` - Valores numéricos
- `minLength:N`, `maxLength:N` - Tamanho strings
- `enum:val1,val2` - Valores permitidos

**Relacionamentos:**
- `relation:hasOne|hasMany|belongsTo|manyToMany`
- `foreignKey:field`
- `through:table` - Para M2M

## Callbacks

Personalize o CRUD editando o handler gerado:

```go
// modules/users/handlers/user_handler.go

func (h *UserHandler) BeforeCreate(c *gin.Context, data map[string]interface{}) error {
    // Hash de senha antes de salvar
    if password, ok := data["password"].(string); ok {
        hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
        data["password"] = string(hashed)
    }
    return nil
}

func (h *UserHandler) AfterGet(c *gin.Context, user models.User) models.User {
    // Remover senha antes de retornar
    user.Password = ""
    return user
}

func (h *UserHandler) OnValidate(data map[string]interface{}, operation string) error {
    // Validações customizadas
    if age, ok := data["age"].(float64); ok {
        if age < 18 {
            return fmt.Errorf("usuário deve ter 18+ anos")
        }
    }
    return nil
}
```

## Estrutura Gerada

```
meu-projeto/
├── cmd/server/           # Aplicação principal
├── config/               # Configurações
│   ├── database/        # Conexão DB
│   ├── routes/          # Registry de rotas
│   ├── modules/         # Registro de módulos
│   └── ...
├── modules/              # Seus módulos
│   └── users/
│       ├── models/      # Models com annotations
│       ├── handlers/    # HTTP handlers
│       ├── services/    # Lógica de negócio
│       ├── repositories/# Acesso a dados
│       └── module.go    # Rotas do módulo
└── migrations/          # SQL migrations
```

## Rotas Automáticas

Ao criar um CRUD, as rotas são registradas automaticamente:

```go
// modules/users/module.go - Gerado automaticamente
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
    userRepo := repositories.NewUserRepository()
    userService := services.NewUserService(userRepo)
    userHandler := handlers.NewUserHandler(userService)

    router.GET("/users", userHandler.List)
    router.GET("/users/:id", userHandler.Get)
    router.POST("/users", userHandler.Create)
    router.PUT("/users/:id", userHandler.Update)
    router.PATCH("/users/:id", userHandler.Patch)
    router.DELETE("/users/:id", userHandler.Delete)
}
```

Registrado em `config/modules/modules.go`:
```go
func RegisterModules(registry *routes.Registry) {
    registry.Register("users", users.NewModule())
}
```

Resultado: Rotas disponíveis em `/api/v1/users` automaticamente!

## Bancos de Dados Suportados

- MySQL
- PostgreSQL
- SQLite

## Licença

MIT License - veja [LICENSE](LICENSE)

## Links

- [CHANGELOG](CHANGELOG.md) - Histórico de versões
- [CONTRIBUTING](CONTRIBUTING.md) - Como contribuir
- [BETA-WARNING](BETA-WARNING.md) - Aviso sobre versão beta

---

**Desenvolvido com ❤️ usando Go e Gin**
