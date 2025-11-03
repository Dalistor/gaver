# 🤝 Guia de Contribuição

Obrigado por considerar contribuir com o Gaver Framework! 

## 📋 Como Contribuir

### 1. Fork o Projeto

```bash
# No GitHub, clique em "Fork"
# Depois clone seu fork:
git clone https://github.com/seu-usuario/gaver.git
cd gaver
```

### 2. Crie uma Branch

```bash
# Para nova feature
git checkout -b feature/minha-feature

# Para bug fix
git checkout -b fix/corrigir-bug

# Para documentação
git checkout -b docs/melhorar-readme
```

### 3. Faça suas Mudanças

```bash
# Edite os arquivos necessários
# Teste suas mudanças
go test ./...
go build cmd/gaver/main.go
```

### 4. Commit com Mensagem Descritiva

Use commits semânticos:

```bash
git commit -m "feat: adiciona comando para deletar módulos"
git commit -m "fix: corrige parser de annotations com vírgulas"
git commit -m "docs: adiciona exemplo de relacionamentos"
```

**Tipos de commit:**
- `feat:` - Nova funcionalidade
- `fix:` - Correção de bug
- `docs:` - Documentação
- `refactor:` - Refatoração de código
- `test:` - Adicionar/modificar testes
- `chore:` - Manutenção (deps, build, etc)
- `perf:` - Melhoria de performance

### 5. Push para seu Fork

```bash
git push origin feature/minha-feature
```

### 6. Abra um Pull Request

1. Vá para `https://github.com/Dalistor/gaver`
2. Clique em **Pull Requests** → **New Pull Request**
3. Selecione seu fork e branch
4. Preencha a descrição do PR

## ✅ Checklist do Pull Request

Antes de abrir um PR, verifique:

- [ ] Código compila sem erros (`go build ./...`)
- [ ] Sem erros de linter
- [ ] Testes passam (se houver)
- [ ] Documentação atualizada (se necessário)
- [ ] CHANGELOG.md atualizado
- [ ] Commits seguem convenção semântica

## 🎨 Padrões de Código

### Formatação

```bash
# Formatar código
go fmt ./...

# Verificar erros
go vet ./...
```

### Nomenclatura

- **Pacotes**: lowercase, singular (`parser`, não `parsers`)
- **Arquivos**: snake_case (`module_generator.go`)
- **Funções exportadas**: PascalCase (`NewGenerator()`)
- **Funções internas**: camelCase (`parseField()`)
- **Constantes**: PascalCase ou UPPER_CASE

### Comentários

Toda função/tipo exportado deve ter comentário GoDoc:

```go
// NewGenerator cria uma nova instância do gerador de código.
// Recebe o caminho dos templates e o diretório de saída.
func NewGenerator(templatesPath, outputPath string) *Generator {
    // ...
}
```

## 🐛 Reportando Bugs

### Antes de reportar:

1. Verifique se já não existe uma issue
2. Use a versão mais recente
3. Teste se consegue reproduzir

### Template de Bug Report

```markdown
**Descrição do Bug**
Descrição clara do que aconteceu.

**Reproduzir**
Passos para reproduzir:
1. Execute `gaver init test`
2. Execute `gaver module create users`
3. Veja o erro

**Comportamento Esperado**
O que deveria acontecer.

**Screenshots/Logs**
Se aplicável.

**Ambiente:**
- OS: [Windows/Linux/Mac]
- Go Version: [1.21]
- Gaver Version: [v0.1.0-beta.1]
```

## 💡 Sugerindo Features

### Template de Feature Request

```markdown
**Sua feature resolve que problema?**
Descrição clara do problema.

**Solução Proposta**
Como você imagina que funcione.

**Alternativas Consideradas**
Outras soluções que você pensou.

**Contexto Adicional**
Qualquer outra informação relevante.
```

## 🧪 Desenvolvimento Local

### Setup do Ambiente

```bash
# 1. Clone
git clone https://github.com/Dalistor/gaver.git
cd gaver

# 2. Instale dependências
go mod download

# 3. Build
go build -o gaver cmd/gaver/main.go

# 4. Teste
./gaver --help
```

### Testando Mudanças

```bash
# Compilar e testar
go build -o gaver cmd/gaver/main.go

# Criar projeto de teste
./gaver init test-project
cd test-project

# Testar comandos
../gaver module create users
../gaver module:model users User name:string email:string
```

## 📚 Áreas que Precisam de Contribuições

### Alta Prioridade

- [ ] Testes unitários
- [ ] Testes de integração
- [ ] Documentação de annotations
- [ ] Exemplos práticos
- [ ] Validação completa do parser

### Média Prioridade

- [ ] QuerySet estilo Django
- [ ] Admin interface
- [ ] Autenticação JWT
- [ ] Suporte a mais bancos de dados

### Baixa Prioridade

- [ ] CLI colorido
- [ ] Progress bars
- [ ] Geração de documentação automática
- [ ] Docker support

## 💬 Comunicação

- **Issues**: Para bugs e sugestões
- **Discussions**: Para dúvidas gerais
- **Pull Requests**: Para contribuições de código

## 📄 Licença

Ao contribuir, você concorda que suas contribuições serão licenciadas sob a MIT License.

## 🙏 Agradecimentos

Toda contribuição é muito apreciada! Obrigado por ajudar a melhorar o Gaver Framework.

---

**Dúvidas?** Abra uma [Discussion](https://github.com/Dalistor/gaver/discussions) no GitHub!

