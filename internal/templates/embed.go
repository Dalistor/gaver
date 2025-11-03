package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed all:*
var TemplatesFS embed.FS

// Generator gera código usando templates embarcados
type Generator struct {
	OutputPath string
}

// New cria um gerador com templates embarcados
func New(outputPath string) *Generator {
	return &Generator{
		OutputPath: outputPath,
	}
}

// Generate processa um template embarcado e gera o arquivo
func (g *Generator) Generate(templateName string, outputFile string, data interface{}) error {
	// 1. Ler template do embed
	templateContent, err := TemplatesFS.ReadFile(templateName)
	if err != nil {
		return fmt.Errorf("erro ao ler template embarcado %s: %w", templateName, err)
	}

	// 2. Parsear template
	tmpl, err := template.New(templateName).Funcs(getFuncMap()).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("erro ao parsear template: %w", err)
	}

	// 3. Criar diretório de saída se não existir
	outputPath := filepath.Join(g.OutputPath, outputFile)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	// 4. Criar arquivo de saída
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer file.Close()

	// 5. Executar template e escrever no arquivo
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("erro ao executar template: %w", err)
	}

	return nil
}

// getFuncMap retorna funções customizadas para usar nos templates
func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"toLower":     toLower,
		"toUpper":     toUpper,
		"capitalize":  capitalize,
		"pluralize":   pluralize,
		"toSnakeCase": toSnakeCase,
		"toCamelCase": toCamelCase,
	}
}

// Funções auxiliares para templates

func toLower(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]|0x20) + s[1:]
}

func toUpper(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]&^0x20) + s[1:]
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return toUpper(s[:1]) + s[1:]
}

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r|0x20) // lowercase
	}
	return string(result)
}

func toCamelCase(s string) string {
	parts := splitBy(s, '_')
	for i := range parts {
		parts[i] = capitalize(parts[i])
	}
	return join(parts, "")
}

func pluralize(s string) string {
	s = toLower(s)
	if hasSuffix(s, "s") {
		return s
	}
	if hasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	if hasSuffix(s, "ch") || hasSuffix(s, "sh") || hasSuffix(s, "x") {
		return s + "es"
	}
	return s + "s"
}

// Helpers simples para evitar imports
func splitBy(s string, sep rune) []string {
	var result []string
	var current []rune
	
	for _, r := range s {
		if r == sep {
			if len(current) > 0 {
				result = append(result, string(current))
				current = []rune{}
			}
		} else {
			current = append(current, r)
		}
	}
	
	if len(current) > 0 {
		result = append(result, string(current))
	}
	
	return result
}

func join(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += sep + parts[i]
	}
	
	return result
}

func hasSuffix(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}

// ListTemplates lista todos os templates disponíveis
func ListTemplates() ([]string, error) {
	var templates []string

	err := fs.WalkDir(TemplatesFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".tmpl" {
			templates = append(templates, filepath.Base(path))
		}
		return nil
	})

	return templates, err
}

