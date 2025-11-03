package generator

// ModelData contém dados para gerar um model
type ModelData struct {
	PackageName string
	ModelName   string
	TableName   string
	Fields      []Field
	Imports     []string
	HasTime     bool
}

// Field representa um campo do model
type Field struct {
	Name     string
	Type     string
	JSONTag  string
	GORMTag  string
	Validate string
}

// HandlerData contém dados para gerar um handler
type HandlerData struct {
	PackageName string
	ModelName   string
	Fields      []Field
}

// ConfigData contém dados para gerar configuração
type ConfigData struct {
	PackageName          string
	DatabaseDriver       string
	DatabaseDriverImport string
	DatabasePort         string
	DatabaseUser         string
	ProjectName          string
}

// ModuleHandlerData contém dados para gerar um handler de módulo
type ModuleHandlerData struct {
	ProjectName   string
	ModuleName    string
	ModelName     string
	ModelNameLower string
	HasList       bool
	HasGet        bool
	HasCreate     bool
	HasUpdate     bool
	HasPatch      bool
	HasDelete     bool
}

// ModuleServiceData contém dados para gerar um service de módulo
type ModuleServiceData struct {
	ProjectName    string
	ModuleName     string
	ModelName      string
	ModelNameLower string
	HasList        bool
	HasGet         bool
	HasCreate      bool
	HasUpdate      bool
	HasDelete      bool
}

// ModuleRepositoryData contém dados para gerar um repository de módulo
type ModuleRepositoryData struct {
	ProjectName    string
	ModuleName     string
	ModelName      string
	ModelNameLower string
	HasList        bool
	HasGet         bool
	HasCreate      bool
	HasUpdate      bool
	HasDelete      bool
}

// ModuleInitData contém dados para gerar module.go inicial
type ModuleInitData struct {
	ModuleName string
}

// ModuleModelData contém dados para gerar um model
type ModuleModelData struct {
	ModelName string
	Fields    []ModelFieldData
}

// ModelFieldData representa um campo do model
type ModelFieldData struct {
	Name       string
	Type       string
	JSONTag    string
	GORMTag    string
	Annotation string
}
