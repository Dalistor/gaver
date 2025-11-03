package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type Generator struct {
	TemplatesPath string
	OutputPath    string
}

// NewGenerator cria uma nova instância do gerador
func NewGenerator(templatesPath, outputPath string) *Generator {
	return &Generator{
		TemplatesPath: templatesPath,
		OutputPath:    outputPath,
	}
}

// Generate processa um template e gera o arquivo
func (g *Generator) Generate(templateName string, outputFile string, data interface{}) error {
	// 1. Carregar o template
	tmplPath := filepath.Join(g.TemplatesPath, templateName)
	tmpl, err := template.New(filepath.Base(templateName)).Funcs(g.getFuncMap()).ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar template: %w", err)
	}

	// 2. Criar diretório de saída se não existir
	outputPath := filepath.Join(g.OutputPath, outputFile)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	// 3. Criar arquivo de saída
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer file.Close()

	// 4. Executar template e escrever no arquivo
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("erro ao executar template: %w", err)
	}

	return nil
}

// getFuncMap retorna funções customizadas para usar nos templates
func (g *Generator) getFuncMap() template.FuncMap {
	return template.FuncMap{
		"toLower":     toLower,
		"toUpper":     toUpper,
		"capitalize":  capitalize,
		"pluralize":   pluralize,
		"toSnakeCase": toSnakeCase,
		"toCamelCase": toCamelCase,
	}
}
