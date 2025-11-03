# 🚀 Gaver Framework

> **Framework web completo para Go com CLI, geração de código e ORM estilo Django**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-0.1.0--beta-orange.svg)](https://github.com/seu-usuario/gaver/releases)

## 📋 Status: Beta Testing

Este projeto está em **fase beta**. A API pode sofrer alterações até a versão 1.0.0.

## ✨ Funcionalidades

- 🎯 **CLI completo** com comandos intuitivos
- 📦 **Sistema de Modules** organizados e reutilizáveis
- 🔖 **Annotations gaverModel** para validações e controle de campos
- 🔄 **CRUD automático** com callbacks personalizáveis (Before/After)
- 📊 **Migrations inteligentes** - detecta mudanças automaticamente
- 🗄️ **ORM sobre GORM** - suporta MySQL, PostgreSQL, SQLite
- 🌐 **Framework HTTP** com Gin
- ⚙️ **Sistema de Rotinas** para tarefas agendadas
- 🔐 **Middlewares** prontos (CORS, Auth, Logger)

## 🚀 Instalação

### Opção 1: Via `go install` (Recomendado quando publicado)

```bash
go install github.com/Dalistor/gaver/cmd/gaver@latest
```

### Opção 2: Build Manual (Beta Testing)

```bash
git clone https://github.com/Dalistor/gaver.git
cd gaver
go build -o gaver cmd/gaver/main.go
```

## 📚 Guia Rápido

### 1. Criar Projeto

```bash
gaver init meu-projeto -d mysql
cd meu-projeto
go mod tidy
```

### 2. Criar Módulo

```bash
gaver module create users
```

### 3. Criar Model com Annotations

```bash
gaver module:model users User name:string email:string:unique age:int
```

Isso gera `modules/users/models/user.go`:

```go
type User struct {
    // gaverModel: primaryKey; autoIncrement
    ID uint `json:"id" gorm:"primaryKey"`

    // gaverModel: writable:post,put,patch; readable
    Name string `json:"name"`

    // gaverModel: writable:post,put,patch; readable; unique
    Email string `json:"email" gorm:"uniqueIndex"`

    // gaverModel: ignore:write; readable
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
```

### 4. Gerar CRUD Completo

```bash
gaver module:crud users User
```

Isso gera:
- ✅ Handler com callbacks Before/After
- ✅ Service com lógica de negócio
- ✅ Repository para acesso a dados
- ✅ Rotas registradas automaticamente

### 5. Migrations

```bash
# Detectar mudanças e gerar migration
gaver makemigrations --name create_users

# Aplicar migrations
gaver migrate up

# Ver status
gaver migrate status
```

### 6. Rodar Servidor

```bash
go run cmd/server/main.go
```

Servidor rodando em `http://localhost:8080`

## 🎯 Annotations gaverModel

Controle total sobre seus models com annotations:

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

## 🔄 Sistema de Callbacks

Personalize comportamento do CRUD:

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

## 🔄 Rotinas Agendadas

Sistema de tarefas em background:

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

## 🛠️ Comandos CLI

### Projeto
```bash
gaver init <nome> [-d database]    # Cria novo projeto
```

### Modules
```bash
gaver module create <nome>                     # Cria módulo
gaver module:model <module> <Model> [campos]   # Cria model
gaver module:crud <module> <Model>             # Gera CRUD completo
  --only=list,get                             # Apenas métodos especificados
  --except=delete                             # Tudo exceto delete
```

### Migrations
```bash
gaver makemigrations [-n nome] [-d]   # Detecta mudanças e gera SQL
gaver migrate up [-s steps]           # Aplica migrations
gaver migrate down [-s steps]         # Reverte migrations
gaver migrate status                  # Status das migrations
```

## 📁 Estrutura Gerada

```
meu-projeto/
├── cmd/server/              # Aplicação principal
├── config/                  # Configurações
│   ├── database/           # Conexão com banco
│   ├── middlewares/        # Middlewares HTTP
│   ├── cors/               # Config CORS
│   ├── env/                # Variáveis ambiente
│   └── routines/           # Tarefas agendadas
├── modules/                # Seus módulos
│   └── users/
│       ├── models/         # Models com annotations
│       ├── handlers/       # Controllers REST
│       ├── services/       # Lógica de negócio
│       ├── repositories/   # Camada de dados
│       └── module.go       # Registro de rotas
├── migrations/             # Migrations SQL
├── .env                    # Variáveis de ambiente
└── go.mod
```

## 🗄️ Suporte a Bancos de Dados

- ✅ MySQL
- ✅ PostgreSQL
- ✅ SQLite

## 🤝 Contribuindo

Contribuições são muito bem-vindas! Este projeto está em beta e qualquer feedback é valioso.

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/NovaFeature`)
3. Commit suas mudanças (`git commit -m 'Adiciona NovaFeature'`)
4. Push para a branch (`git push origin feature/NovaFeature`)
5. Abra um Pull Request

## 📝 Roadmap

### v0.1.0-beta (Atual)
- [x] CLI básico
- [x] Sistema de modules
- [x] Geração de CRUD
- [x] Annotations gaverModel
- [x] Migrations

### v0.2.0-beta (Próximo)
- [ ] QuerySet estilo Django
- [ ] Admin interface
- [ ] Autenticação JWT
- [ ] WebSockets
- [ ] Testes automatizados

### v1.0.0 (Estável)
- [ ] Documentação completa
- [ ] Cobertura de testes 80%+
- [ ] API estável
- [ ] Exemplos completos

## 📄 Licença

Este projeto está sob a licença MIT - veja [LICENSE](LICENSE) para detalhes.

## ⭐ Apoie o Projeto

Se você achou útil, considere dar uma estrela no GitHub! ⭐

---

**Desenvolvido com ❤️ usando Go e Gin**

