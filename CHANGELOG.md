# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/),
e este projeto adere ao [Semantic Versioning](https://semver.org/lang/pt-BR/).

## [Unreleased]

### Planejado para próximas versões beta
- QuerySet API estilo Django
- Testes automatizados completos
- Documentação expandida
- Exemplos de projetos completos
- Melhorias de performance

## [0.1.0-beta] - 2025-11-03

### Instalação

```bash
go install github.com/Dalistor/gaver/cmd/gaver@latest
```

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

### Fase Beta (0.x.x-beta) - Esperado: 6-12 meses
Durante a fase beta:
- ⚠️ API pode mudar sem aviso
- 🧪 Use para testes e desenvolvimento
- ❌ Não recomendado para produção
- 💬 Feedback é essencial
- 🐛 Bugs esperados

### Versões Planejadas

**Beta Phases:**
- `v0.1.0-beta` - Core framework (atual)
- `v0.2.0-beta` - QuerySet e validações avançadas
- `v0.3.0-beta` - Testes e exemplos
- `v0.4.0-beta` - Performance e otimizações
- `v0.5.0-beta` - Features avançadas
- `v0.9.0-beta` - Feature freeze
- `v1.0.0-rc.1` - Release candidate

**Stable:**
- `v1.0.0` - Primeira versão estável (quando API estiver madura)

### Critérios para v1.0.0
- [ ] API estável sem breaking changes por 2+ meses
- [ ] Cobertura de testes 80%+
- [ ] Documentação completa
- [ ] 50+ usuários usando em desenvolvimento
- [ ] Performance validada
- [ ] Zero bugs críticos conhecidos

[Unreleased]: https://github.com/Dalistor/gaver/compare/v0.1.0-beta...HEAD
[0.1.0-beta]: https://github.com/Dalistor/gaver/releases/tag/v0.1.0-beta

