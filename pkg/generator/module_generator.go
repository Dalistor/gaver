package generator

import (
	"path/filepath"
)

// ModuleGenerator gera código para módulos
type ModuleGenerator struct {
	templatesPath string
	projectName   string
}

// NewModuleGenerator cria um novo gerador de módulos
func NewModuleGenerator(templatesPath, projectName string) *ModuleGenerator {
	return &ModuleGenerator{
		templatesPath: templatesPath,
		projectName:   projectName,
	}
}

// GenerateHandler gera um handler usando template
func (g *ModuleGenerator) GenerateHandler(moduleName, modelName string, methods map[string]bool) error {
	gen := NewGenerator(g.templatesPath, "modules")

	data := ModuleHandlerData{
		ProjectName:    g.projectName,
		ModuleName:     moduleName,
		ModelName:      modelName,
		ModelNameLower: ToLower(modelName),
		HasList:        methods["list"],
		HasGet:         methods["get"],
		HasCreate:      methods["create"],
		HasUpdate:      methods["update"],
		HasPatch:       methods["patch"],
		HasDelete:      methods["delete"],
	}

	outputPath := filepath.Join(moduleName, "handlers", ToSnakeCase(modelName)+"_handler.go")
	return gen.Generate("module_handler.tmpl", outputPath, data)
}

// GenerateService gera um service usando template
func (g *ModuleGenerator) GenerateService(moduleName, modelName string, methods map[string]bool) error {
	gen := NewGenerator(g.templatesPath, "modules")

	data := ModuleServiceData{
		ProjectName:    g.projectName,
		ModuleName:     moduleName,
		ModelName:      modelName,
		ModelNameLower: ToLower(modelName),
		HasList:        methods["list"],
		HasGet:         methods["get"],
		HasCreate:      methods["create"],
		HasUpdate:      methods["update"],
		HasDelete:      methods["delete"],
	}

	outputPath := filepath.Join(moduleName, "services", ToSnakeCase(modelName)+"_service.go")
	return gen.Generate("module_service.tmpl", outputPath, data)
}

// GenerateRepository gera um repository usando template
func (g *ModuleGenerator) GenerateRepository(moduleName, modelName string, methods map[string]bool) error {
	gen := NewGenerator(g.templatesPath, "modules")

	data := ModuleRepositoryData{
		ProjectName:    g.projectName,
		ModuleName:     moduleName,
		ModelName:      modelName,
		ModelNameLower: ToLower(modelName),
		HasList:        methods["list"],
		HasGet:         methods["get"],
		HasCreate:      methods["create"],
		HasUpdate:      methods["update"],
		HasDelete:      methods["delete"],
	}

	outputPath := filepath.Join(moduleName, "repositories", ToSnakeCase(modelName)+"_repository.go")
	return gen.Generate("module_repository.tmpl", outputPath, data)
}

// GenerateHandlerWithMetadata gera handler usando metadata do model parseado
func (g *ModuleGenerator) GenerateHandlerWithMetadata(moduleName, modelName string, metadata interface{}, methods map[string]bool) error {
	// Por enquanto, usa o gerador normal
	// TODO: Implementar geração mais inteligente baseada em metadata
	return g.GenerateHandler(moduleName, modelName, methods)
}

