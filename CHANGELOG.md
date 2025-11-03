# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/),
e este projeto adere ao [Semantic Versioning](https://semver.org/lang/pt-BR/).

## [Unreleased]

## [0.1.0-beta.1] - 2025-11-03

### Adicionado

#### Sistema de Modules
- Comando `gaver module create` para criar módulos
- Comando `gaver module:model` para criar models com annotations
- Comando `gaver module:crud` para gerar CRUD completo
- Estrutura modular com models, handlers, services, repositories

#### Sistema de Annotations
- Parser AST para ler annotations `gaverModel`
- Suporte a tags: writable, readable, required, unique
- Validações: email, url, min, max, minLength, maxLength, pattern, enum
- Relacionamentos: hasOne, hasMany, belongsTo, manyToMany

#### Sistema de CRUD
- Geração automática de handlers com callbacks
- Callbacks: Before/After para List, Get, Create, Update, Patch, Delete
- Filtragem automática de campos baseada em writable/readable
- Validação automática baseada em annotations
- Registro automático de rotas no module.go

#### Sistema de Migrations
- Comando `gaver makemigrations` para detectar mudanças
- Comando `gaver migrate up/down` para aplicar/reverter
- Comando `gaver migrate status` para ver status
- Geração automática de SQL UP/DOWN
- Suporte a MySQL, PostgreSQL, SQLite

#### Sistema de Projeto
- Comando `gaver init` para criar projeto inicial
- Templates para configuração (database, middlewares, cors, env)
- Sistema de rotinas (tarefas agendadas)
- Integração com Gin Framework
- Sistema de validações

#### Templates
- module_init.tmpl - Arquivo module.go inicial
- module_model.tmpl - Models com annotations
- module_handler.tmpl - Handlers CRUD com callbacks
- module_service.tmpl - Services
- module_repository.tmpl - Repositories
- config_*.tmpl - Arquivos de configuração
- main.tmpl - Arquivo main.go do projeto
- routines.tmpl - Sistema de rotinas

### Mudanças
- N/A (primeira versão)

### Removido
- N/A (primeira versão)

### Corrigido
- N/A (primeira versão)

### Segurança
- N/A (primeira versão)

---

## Notas de Versionamento

### Versões Beta (0.x.x-beta)
- API pode mudar sem aviso
- Use para testes e feedback
- Não recomendado para produção

### Versão 1.0.0
- API estável
- Retrocompatibilidade garantida
- Pronto para produção

[Unreleased]: https://github.com/seu-usuario/gaver/compare/v0.1.0-beta.1...HEAD
[0.1.0-beta.1]: https://github.com/seu-usuario/gaver/releases/tag/v0.1.0-beta.1

