package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	templates "github.com/Dalistor/gaver/internal/templates"
	"github.com/Dalistor/gaver/pkg/generator"
	"github.com/Dalistor/gaver/pkg/parser"
)

// CreateModule cria a estrutura de pastas de um módulo
func CreateModule(moduleName string) error {
	basePath := filepath.Join("modules", moduleName)

	// Verificar se módulo já existe
	if _, err := os.Stat(basePath); err == nil {
		return fmt.Errorf("módulo '%s' já existe", moduleName)
	}

	// Criar pastas do módulo
	dirs := []string{
		basePath,
		filepath.Join(basePath, "models"),
		filepath.Join(basePath, "handlers"),
		filepath.Join(basePath, "services"),
		filepath.Join(basePath, "repositories"),
		filepath.Join(basePath, "validators"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %w", dir, err)
		}
	}

	// Criar arquivo module.go
	if err := createModuleFile(basePath, moduleName); err != nil {
		return fmt.Errorf("erro ao criar module.go: %w", err)
	}

	// Criar .gitkeep nas pastas vazias
	emptyDirs := []string{"models", "handlers", "services", "repositories", "validators"}
	for _, dir := range emptyDirs {
		gitkeep := filepath.Join(basePath, dir, ".gitkeep")
		if err := os.WriteFile(gitkeep, []byte(""), 0644); err != nil {
			return fmt.Errorf("erro ao criar .gitkeep: %w", err)
		}
	}

	return nil
}

func createModuleFile(basePath, moduleName string) error {
	gen := templates.New(basePath)

	data := generator.ModuleInitData{
		ModuleName: moduleName,
	}

	return gen.Generate("module_init.tmpl", "module.go", data)
}

// CreateModel cria um arquivo de model com annotations
func CreateModel(moduleName, modelName string, fields []string) error {
	// Verificar se módulo existe
	if _, err := os.Stat(filepath.Join("modules", moduleName)); os.IsNotExist(err) {
		return fmt.Errorf("módulo '%s' não existe. Use 'gaver module create %s' primeiro", moduleName, moduleName)
	}

	// Parsear campos
	parsedFields := parseFields(fields)

	// Converter para ModelFieldData
	modelFields := make([]generator.ModelFieldData, len(parsedFields))
	for i, field := range parsedFields {
		modelFields[i] = generator.ModelFieldData{
			Name:       field.Name,
			Type:       field.Type,
			JSONTag:    toSnakeCase(field.Name),
			GORMTag:    generateGORMTag(field),
			Annotation: generateAnnotation(field),
		}
	}

	// Usar generator com template embarcado
	gen := templates.New(filepath.Join("modules", moduleName, "models"))

	data := generator.ModuleModelData{
		ModelName: modelName,
		Fields:    modelFields,
	}

	filename := toSnakeCase(modelName) + ".go"
	return gen.Generate("module_model.tmpl", filename, data)
}

// CreateCRUD gera handlers, services e repositories lendo o model existente
func CreateCRUD(moduleName, modelName string, only, except []string) error {
	// Verificar se módulo existe
	if _, err := os.Stat(filepath.Join("modules", moduleName)); os.IsNotExist(err) {
		return fmt.Errorf("módulo '%s' não existe", moduleName)
	}

	// Verificar se model existe
	modelFile := filepath.Join("modules", moduleName, "models", toSnakeCase(modelName)+".go")
	if _, err := os.Stat(modelFile); os.IsNotExist(err) {
		return fmt.Errorf("model '%s' não existe no módulo '%s'", modelName, moduleName)
	}

	// Parsear o model para obter metadata
	metadata, err := parser.ParseModelFile(modelFile)
	if err != nil {
		return fmt.Errorf("erro ao parsear model: %w", err)
	}

	// Determinar quais métodos gerar
	methods := determineMethods(only, except)

	// Gerar handler com metadata
	if err := generateHandlerWithMetadata(moduleName, modelName, metadata, methods); err != nil {
		return fmt.Errorf("erro ao gerar handler: %w", err)
	}

	// Gerar service
	if err := generateService(moduleName, modelName, methods); err != nil {
		return fmt.Errorf("erro ao gerar service: %w", err)
	}

	// Gerar repository
	if err := generateRepository(moduleName, modelName, methods); err != nil {
		return fmt.Errorf("erro ao gerar repository: %w", err)
	}

	// Atualizar module.go com as rotas
	if err := updateModuleRoutes(moduleName, modelName, methods); err != nil {
		return fmt.Errorf("erro ao atualizar rotas: %w", err)
	}

	return nil
}

func parseFields(fields []string) []FieldDef {
	result := []FieldDef{}

	for _, field := range fields {
		parts := strings.Split(field, ":")
		if len(parts) < 2 {
			continue
		}

		fieldDef := FieldDef{
			Name: capitalize(parts[0]),
			Type: getGoType(parts[1]),
			Tags: []string{},
		}

		// Tags adicionais (unique, index, etc)
		if len(parts) > 2 {
			fieldDef.Tags = parts[2:]
		}

		result = append(result, fieldDef)
	}

	return result
}

type FieldDef struct {
	Name string
	Type string
	Tags []string
}

// generateModelContent foi removida - agora usa templates

func generateAnnotation(field FieldDef) string {
	annotations := []string{"writable:post,put,patch", "readable"}

	for _, tag := range field.Tags {
		switch tag {
		case "unique":
			annotations = append(annotations, "unique")
		case "required":
			annotations = append(annotations, "required")
		case "index":
			annotations = append(annotations, "index")
		}
	}

	return strings.Join(annotations, "; ")
}

func generateGORMTag(field FieldDef) string {
	tags := []string{}

	for _, tag := range field.Tags {
		switch tag {
		case "unique":
			tags = append(tags, "uniqueIndex")
		case "index":
			tags = append(tags, "index")
		}
	}

	return strings.Join(tags, ";")
}

func getGoType(sqlType string) string {
	typeMap := map[string]string{
		"string":   "string",
		"int":      "int",
		"int64":    "int64",
		"uint":     "uint",
		"uint64":   "uint64",
		"float":    "float64",
		"float64":  "float64",
		"bool":     "bool",
		"time":     "time.Time",
		"date":     "time.Time",
		"datetime": "time.Time",
		"text":     "string",
	}

	if goType, ok := typeMap[strings.ToLower(sqlType)]; ok {
		return goType
	}
	return "string"
}

func determineMethods(only, except []string) map[string]bool {
	allMethods := map[string]bool{
		"list":   true,
		"get":    true,
		"create": true,
		"update": true,
		"patch":  true,
		"delete": true,
	}

	// Se only está definido, usar apenas esses
	if len(only) > 0 {
		allMethods = map[string]bool{}
		for _, method := range only {
			allMethods[strings.ToLower(method)] = true
		}
		return allMethods
	}

	// Se except está definido, remover esses
	if len(except) > 0 {
		for _, method := range except {
			delete(allMethods, strings.ToLower(method))
		}
	}

	return allMethods
}

func generateHandler(moduleName, modelName string, methods map[string]bool) error {
	// Obter nome do projeto atual
	projectName, err := getProjectName()
	if err != nil {
		projectName = "gaver-project" // fallback
	}

	gen := generator.NewModuleGenerator("templates", projectName)
	return gen.GenerateHandler(moduleName, modelName, methods)
}

func generateHandlerWithMetadata(moduleName, modelName string, metadata *parser.ModelMetadata, methods map[string]bool) error {
	projectName, err := getProjectName()
	if err != nil {
		projectName = "gaver-project"
	}

	gen := generator.NewModuleGenerator("templates", projectName)
	return gen.GenerateHandlerWithMetadata(moduleName, modelName, metadata, methods)
}

func generateService(moduleName, modelName string, methods map[string]bool) error {
	projectName, err := getProjectName()
	if err != nil {
		projectName = "gaver-project"
	}

	gen := generator.NewModuleGenerator("templates", projectName)
	return gen.GenerateService(moduleName, modelName, methods)
}

func generateRepository(moduleName, modelName string, methods map[string]bool) error {
	projectName, err := getProjectName()
	if err != nil {
		projectName = "gaver-project"
	}

	gen := generator.NewModuleGenerator("templates", projectName)
	return gen.GenerateRepository(moduleName, modelName, methods)
}

func getProjectName() (string, error) {
	// Tenta ler go.mod para pegar o nome do projeto
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	return "", fmt.Errorf("nome do projeto não encontrado")
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// updateModuleRoutes atualiza o arquivo module.go com as rotas do CRUD
func updateModuleRoutes(moduleName, modelName string, methods map[string]bool) error {
	moduleFile := filepath.Join("modules", moduleName, "module.go")

	// Ler arquivo existente
	content, err := os.ReadFile(moduleFile)
	if err != nil {
		return err
	}

	contentStr := string(content)

	// Preparar código das rotas
	routesCode := generateRoutesCode(moduleName, modelName, methods)

	// Verificar se já existe código de rotas
	if strings.Contains(contentStr, "RegisterRoutes") {
		// Substituir o conteúdo da função RegisterRoutes
		contentStr = replaceRoutesFunction(contentStr, routesCode)
	} else {
		// Adicionar função RegisterRoutes
		contentStr = addRoutesFunction(contentStr, routesCode)
	}

	// Salvar arquivo atualizado
	return os.WriteFile(moduleFile, []byte(contentStr), 0644)
}

func generateRoutesCode(moduleName, modelName string, methods map[string]bool) string {
	var code strings.Builder

	modelLower := toLower(modelName)
	handlerVar := modelLower + "Handler"

	code.WriteString("\t// Inicializar handler\n")
	code.WriteString(fmt.Sprintf("\t%sRepo := repositories.New%sRepository()\n", modelLower, modelName))
	code.WriteString(fmt.Sprintf("\t%sService := services.New%sService(%sRepo)\n", modelLower, modelName, modelLower))
	code.WriteString(fmt.Sprintf("\t%s := handlers.New%sHandler(%sService)\n\n", handlerVar, modelName, modelLower))

	resourcePath := "/" + toSnakeCase(pluralize(modelName))

	if methods["list"] {
		code.WriteString(fmt.Sprintf("\trouter.GET(\"%s\", %s.List)\n", resourcePath, handlerVar))
	}
	if methods["get"] {
		code.WriteString(fmt.Sprintf("\trouter.GET(\"%s/:id\", %s.Get)\n", resourcePath, handlerVar))
	}
	if methods["create"] {
		code.WriteString(fmt.Sprintf("\trouter.POST(\"%s\", %s.Create)\n", resourcePath, handlerVar))
	}
	if methods["update"] {
		code.WriteString(fmt.Sprintf("\trouter.PUT(\"%s/:id\", %s.Update)\n", resourcePath, handlerVar))
	}
	if methods["patch"] {
		code.WriteString(fmt.Sprintf("\trouter.PATCH(\"%s/:id\", %s.Patch)\n", resourcePath, handlerVar))
	}
	if methods["delete"] {
		code.WriteString(fmt.Sprintf("\trouter.DELETE(\"%s/:id\", %s.Delete)\n", resourcePath, handlerVar))
	}

	return code.String()
}

func replaceRoutesFunction(content, routesCode string) string {
	// Encontrar função RegisterRoutes e substituir conteúdo
	startMarker := "func (m *Module) RegisterRoutes(router *gin.RouterGroup) {"

	startIdx := strings.Index(content, startMarker)
	if startIdx == -1 {
		return content
	}

	// Encontrar o fechamento da função (procurar o } correspondente)
	braceCount := 0
	endIdx := startIdx + len(startMarker)

	for i := endIdx; i < len(content); i++ {
		if content[i] == '{' {
			braceCount++
		} else if content[i] == '}' {
			if braceCount == 0 {
				endIdx = i
				break
			}
			braceCount--
		}
	}

	// Reconstruir conteúdo
	newContent := content[:startIdx+len(startMarker)] + "\n" + routesCode + "\n" + content[endIdx:]
	return newContent
}

func addRoutesFunction(content, routesCode string) string {
	// Adicionar imports necessários se não existirem
	if !strings.Contains(content, "\"github.com/gin-gonic/gin\"") {
		// Adicionar import do gin
		importIdx := strings.Index(content, "import (")
		if importIdx != -1 {
			insertPoint := importIdx + len("import (") + 1
			content = content[:insertPoint] + "\t\"github.com/gin-gonic/gin\"\n" + content[insertPoint:]
		}
	}

	// Adicionar função RegisterRoutes antes do último }
	lastBrace := strings.LastIndex(content, "}")
	if lastBrace != -1 {
		routesFunc := fmt.Sprintf(`
// RegisterRoutes registra as rotas do módulo
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
%s
}
`, routesCode)
		content = content[:lastBrace] + routesFunc + "\n" + content[lastBrace:]
	}

	return content
}

func pluralize(s string) string {
	s = strings.ToLower(s)
	if strings.HasSuffix(s, "s") {
		return s
	}
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	if strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") || strings.HasSuffix(s, "x") {
		return s + "es"
	}
	return s + "s"
}

// toLower converte primeira letra para minúscula
func toLower(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// toSnakeCase converte CamelCase para snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
