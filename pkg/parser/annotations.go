package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"strings"
)

// FieldMetadata contém metadados de um campo do model
type FieldMetadata struct {
	Name        string
	Type        string
	Writable    []string          // POST, PUT, PATCH
	Readable    bool
	Required    bool
	Unique      bool
	PrimaryKey  bool
	AutoInc     bool
	Validations map[string]string
	Relation    *Relation
	JSONTag     string
	GORMTag     string
	Default     string
	Index       bool
	Ignore      bool
	IgnoreWrite bool
	IgnoreRead  bool
}

// Relation representa um relacionamento entre models
type Relation struct {
	Type       string // hasOne, hasMany, belongsTo, manyToMany
	ForeignKey string
	Through    string
	Model      string
}

// ModelMetadata contém metadados completos de um model
type ModelMetadata struct {
	Name       string
	Package    string
	TableName  string
	Fields     []FieldMetadata
	Imports    []string
}

// ParseModelFile lê um arquivo .go e extrai metadata
func ParseModelFile(filePath string) (*ModelMetadata, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear arquivo: %w", err)
	}

	metadata := &ModelMetadata{
		Package: node.Name.Name,
		Fields:  []FieldMetadata{},
		Imports: []string{},
	}

	// Extrair imports
	for _, imp := range node.Imports {
		if imp.Path != nil {
			metadata.Imports = append(metadata.Imports, strings.Trim(imp.Path.Value, `"`))
		}
	}

	// Encontrar struct type
	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		metadata.Name = typeSpec.Name.Name
		metadata.TableName = toSnakeCase(pluralize(typeSpec.Name.Name))

		// Parsear cada campo
		for _, field := range structType.Fields.List {
			if len(field.Names) > 0 {
				fieldMeta := parseField(field)
				metadata.Fields = append(metadata.Fields, fieldMeta)
			}
		}

		return false
	})

	if metadata.Name == "" {
		return nil, fmt.Errorf("nenhuma struct encontrada no arquivo")
	}

	return metadata, nil
}

func parseField(field *ast.Field) FieldMetadata {
	meta := FieldMetadata{
		Name:        field.Names[0].Name,
		Writable:    []string{},
		Validations: make(map[string]string),
		Readable:    true, // Padrão é readable
	}

	// Extrair tipo
	meta.Type = getTypeString(field.Type)

	// Parsear tags struct
	if field.Tag != nil {
		meta.JSONTag = extractTag(field.Tag.Value, "json")
		meta.GORMTag = extractTag(field.Tag.Value, "gorm")
	}

	// Parsear annotation gaverModel dos comentários
	if field.Doc != nil {
		for _, comment := range field.Doc.List {
			text := comment.Text
			if strings.HasPrefix(text, "// gaverModel:") || strings.HasPrefix(text, "//gaverModel:") {
				parseAnnotation(text, &meta)
			}
		}
	}

	return meta
}

func parseAnnotation(comment string, meta *FieldMetadata) {
	// Remove "// gaverModel:" ou "//gaverModel:"
	content := strings.TrimPrefix(comment, "// gaverModel:")
	content = strings.TrimPrefix(content, "//gaverModel:")
	content = strings.TrimSpace(content)

	// Split por ";" ou ","
	parts := splitAnnotation(content)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, ":") {
			// Tag com valor: "writable:post,put"
			kv := strings.SplitN(part, ":", 2)
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])

			switch key {
			case "writable":
				methods := strings.Split(value, ",")
				for i, m := range methods {
					methods[i] = strings.ToUpper(strings.TrimSpace(m))
				}
				meta.Writable = methods

			case "ignore":
				if value == "write" {
					meta.IgnoreWrite = true
				} else if value == "read" {
					meta.IgnoreRead = true
				} else {
					meta.Ignore = true
				}

			case "min", "max", "minLength", "maxLength", "pattern", "enum":
				meta.Validations[key] = value

			case "default":
				meta.Default = value

			case "relation":
				if meta.Relation == nil {
					meta.Relation = &Relation{}
				}
				meta.Relation.Type = value

			case "foreignKey":
				if meta.Relation == nil {
					meta.Relation = &Relation{}
				}
				meta.Relation.ForeignKey = value

			case "through":
				if meta.Relation == nil {
					meta.Relation = &Relation{}
				}
				meta.Relation.Through = value

			case "model":
				if meta.Relation == nil {
					meta.Relation = &Relation{}
				}
				meta.Relation.Model = value
			}
		} else {
			// Tag boolean: "readable", "required"
			switch part {
			case "readable":
				meta.Readable = true
			case "required":
				meta.Required = true
			case "unique":
				meta.Unique = true
			case "primaryKey":
				meta.PrimaryKey = true
			case "autoIncrement", "autoInc":
				meta.AutoInc = true
			case "index":
				meta.Index = true
			case "ignore":
				meta.Ignore = true
			case "email":
				meta.Validations["email"] = "true"
			case "url":
				meta.Validations["url"] = "true"
			}
		}
	}
}

func splitAnnotation(content string) []string {
	// Primeiro tenta split por ";"
	if strings.Contains(content, ";") {
		return strings.Split(content, ";")
	}
	// Senão, split por ","
	return strings.Split(content, ",")
}

func getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", getTypeString(t.X), t.Sel.Name)
	case *ast.ArrayType:
		return "[]" + getTypeString(t.Elt)
	case *ast.StarExpr:
		return "*" + getTypeString(t.X)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", getTypeString(t.Key), getTypeString(t.Value))
	default:
		return fmt.Sprintf("%T", t)
	}
}

func extractTag(tagValue, tagName string) string {
	// Remove os backticks
	tagValue = strings.Trim(tagValue, "`")

	// Parse das tags
	re := regexp.MustCompile(tagName + `:"([^"]*)"`)
	matches := re.FindStringSubmatch(tagValue)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

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

// ToSnakeCase exportada para uso externo
func ToSnakeCase(s string) string {
	return toSnakeCase(s)
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

// GetGoType retorna o tipo Go baseado no tipo SQL
func GetGoType(sqlType string) string {
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

// GetSQLType retorna o tipo SQL baseado no tipo Go
func GetSQLType(goType string, driver string) string {
	// Remove ponteiro
	goType = strings.TrimPrefix(goType, "*")

	switch driver {
	case "mysql":
		return getMySQLType(goType)
	case "postgres":
		return getPostgresType(goType)
	case "sqlite":
		return getSQLiteType(goType)
	default:
		return getMySQLType(goType)
	}
}

func getMySQLType(goType string) string {
	typeMap := map[string]string{
		"string":    "VARCHAR(255)",
		"int":       "INT",
		"int64":     "BIGINT",
		"uint":      "INT UNSIGNED",
		"uint64":    "BIGINT UNSIGNED",
		"float64":   "DOUBLE",
		"bool":      "BOOLEAN",
		"time.Time": "TIMESTAMP",
		"[]byte":    "BLOB",
	}

	if sqlType, ok := typeMap[goType]; ok {
		return sqlType
	}
	return "VARCHAR(255)"
}

func getPostgresType(goType string) string {
	typeMap := map[string]string{
		"string":    "VARCHAR(255)",
		"int":       "INTEGER",
		"int64":     "BIGINT",
		"uint":      "INTEGER",
		"uint64":    "BIGINT",
		"float64":   "DOUBLE PRECISION",
		"bool":      "BOOLEAN",
		"time.Time": "TIMESTAMP",
		"[]byte":    "BYTEA",
	}

	if sqlType, ok := typeMap[goType]; ok {
		return sqlType
	}
	return "VARCHAR(255)"
}

func getSQLiteType(goType string) string {
	typeMap := map[string]string{
		"string":    "TEXT",
		"int":       "INTEGER",
		"int64":     "INTEGER",
		"uint":      "INTEGER",
		"uint64":    "INTEGER",
		"float64":   "REAL",
		"bool":      "INTEGER",
		"time.Time": "DATETIME",
		"[]byte":    "BLOB",
	}

	if sqlType, ok := typeMap[goType]; ok {
		return sqlType
	}
	return "TEXT"
}

// IsWritableInMethod verifica se o campo pode ser escrito no método HTTP especificado
func (f *FieldMetadata) IsWritableInMethod(method string) bool {
	method = strings.ToUpper(method)

	// Se ignore ou ignoreWrite, não pode escrever
	if f.Ignore || f.IgnoreWrite {
		return false
	}

	// Se não tem writable definido, permite todos
	if len(f.Writable) == 0 {
		return true
	}

	// Verifica se o método está na lista
	for _, m := range f.Writable {
		if m == method {
			return true
		}
	}

	return false
}

// IsReadable verifica se o campo pode ser lido
func (f *FieldMetadata) IsReadable() bool {
	return f.Readable && !f.Ignore && !f.IgnoreRead
}

// ValidateValue valida um valor baseado nas validações do campo
func (f *FieldMetadata) ValidateValue(value interface{}) error {
	// Required
	if f.Required && isNilOrEmpty(value) {
		return fmt.Errorf("campo '%s' é obrigatório", f.Name)
	}

	// Email
	if emailVal, ok := f.Validations["email"]; ok && emailVal == "true" {
		if str, ok := value.(string); ok {
			if !isValidEmail(str) {
				return fmt.Errorf("campo '%s' deve ser um email válido", f.Name)
			}
		}
	}

	// Min/Max para números
	if minVal, ok := f.Validations["min"]; ok {
		if !validateMin(value, minVal) {
			return fmt.Errorf("campo '%s' deve ser maior ou igual a %s", f.Name, minVal)
		}
	}

	if maxVal, ok := f.Validations["max"]; ok {
		if !validateMax(value, maxVal) {
			return fmt.Errorf("campo '%s' deve ser menor ou igual a %s", f.Name, maxVal)
		}
	}

	// MinLength/MaxLength para strings
	if str, ok := value.(string); ok {
		if minLen, ok := f.Validations["minLength"]; ok {
			if !validateMinLength(str, minLen) {
				return fmt.Errorf("campo '%s' deve ter no mínimo %s caracteres", f.Name, minLen)
			}
		}

		if maxLen, ok := f.Validations["maxLength"]; ok {
			if !validateMaxLength(str, maxLen) {
				return fmt.Errorf("campo '%s' deve ter no máximo %s caracteres", f.Name, maxLen)
			}
		}
	}

	return nil
}

func isNilOrEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	}
	
	return false
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func validateMin(value interface{}, min string) bool {
	// TODO: Implementar validação de mínimo
	return true
}

func validateMax(value interface{}, max string) bool {
	// TODO: Implementar validação de máximo
	return true
}

func validateMinLength(str, minLen string) bool {
	// TODO: Implementar validação de tamanho mínimo
	return true
}

func validateMaxLength(str, maxLen string) bool {
	// TODO: Implementar validação de tamanho máximo
	return true
}

