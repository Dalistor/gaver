package migrations

import (
	"fmt"
	"strings"

	"github.com/Dalistor/gaver/pkg/parser"
)

// SQLGenerator gera SQL para migrations
type SQLGenerator struct{}

// NewSQLGenerator cria um novo gerador de SQL
func NewSQLGenerator() *SQLGenerator {
	return &SQLGenerator{}
}

// Generate gera SQL UP e DOWN para as mudanças
func (g *SQLGenerator) Generate(changes []SchemaChange, driver string) (up string, down string) {
	var upSQL []string
	var downSQL []string

	for _, change := range changes {
		u, d := g.generateForChange(change, driver)
		if u != "" {
			upSQL = append(upSQL, u)
		}
		if d != "" {
			downSQL = append(downSQL, d)
		}
	}

	return strings.Join(upSQL, "\n\n"), strings.Join(downSQL, "\n\n")
}

func (g *SQLGenerator) generateForChange(change SchemaChange, driver string) (up string, down string) {
	switch change.Type {
	case "CREATE_TABLE":
		return g.generateCreateTable(change, driver)
	case "DROP_TABLE":
		return g.generateDropTable(change, driver)
	case "ADD_COLUMN":
		return g.generateAddColumn(change, driver)
	case "DROP_COLUMN":
		return g.generateDropColumn(change, driver)
	case "ALTER_COLUMN":
		return g.generateAlterColumn(change, driver)
	default:
		return "", ""
	}
}

func (g *SQLGenerator) generateCreateTable(change SchemaChange, driver string) (up string, down string) {
	if change.Model == nil {
		// Fallback se não tiver metadata
		up = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`, change.TableName)
		down = fmt.Sprintf("DROP TABLE IF EXISTS %s;", change.TableName)
		return up, down
	}

	// Gerar DDL baseado nos campos do model
	var columns []string

	for _, field := range change.Model.Fields {
		// Ignorar campos que não devem virar colunas
		if g.shouldSkipField(field) {
			continue
		}

		columnName := field.JSONTag
		if columnName == "" || columnName == "-" {
			columnName = toSnakeCase(field.Name)
		}

		columnDef := g.generateColumnDefinition(field, driver)
		if columnDef != "" {
			columns = append(columns, fmt.Sprintf("    %s %s", columnName, columnDef))
		}
	}

	tableDef := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n%s\n)", 
		change.TableName, 
		strings.Join(columns, ",\n"))

	if driver == "mysql" {
		tableDef += " ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci"
	}
	tableDef += ";"

	down = fmt.Sprintf("DROP TABLE IF EXISTS %s;", change.TableName)

	return tableDef, down
}

// generateColumnDefinition gera a definição SQL de uma coluna
func (g *SQLGenerator) generateColumnDefinition(field parser.FieldMetadata, driver string) string {
	sqlType := g.goTypeToSQL(field.Type, field.GORMTag, driver)
	if sqlType == "" {
		return ""
	}

	// Adicionar constraints
	constraints := []string{sqlType}

	// PRIMARY KEY
	if field.PrimaryKey {
		constraints = append(constraints, "PRIMARY KEY")
	}

	// NOT NULL (required)
	if field.Required && !field.PrimaryKey {
		constraints = append(constraints, "NOT NULL")
	}

	// UNIQUE
	if field.Unique {
		constraints = append(constraints, "UNIQUE")
	}

	// DEFAULT (extrair do GORM tag se existir)
	if field.Default != "" {
		constraints = append(constraints, fmt.Sprintf("DEFAULT %s", field.Default))
	}

	// AUTO_INCREMENT para MySQL
	if field.AutoInc && driver == "mysql" {
		constraints = append(constraints, "AUTO_INCREMENT")
	}

	return strings.Join(constraints, " ")
}

// goTypeToSQL converte tipo Go para tipo SQL
func (g *SQLGenerator) goTypeToSQL(goType string, gormTag string, driver string) string {
	// Se tem tipo específico no GORM tag, usar ele
	if strings.Contains(gormTag, "type:") {
		// Extrair type: do gorm tag
		parts := strings.Split(gormTag, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "type:") {
				return strings.TrimPrefix(part, "type:")
			}
		}
	}

	// Mapeamento de tipos Go para SQL
	switch goType {
	case "string":
		return "VARCHAR(255)"
	case "int", "int32":
		return "INT"
	case "int64":
		return "BIGINT"
	case "uint", "uint32":
		return "INT UNSIGNED"
	case "uint64":
		return "BIGINT UNSIGNED"
	case "float32":
		return "FLOAT"
	case "float64":
		return "DOUBLE"
	case "bool":
		if driver == "postgres" {
			return "BOOLEAN"
		}
		return "TINYINT(1)"
	case "time.Time":
		if driver == "postgres" {
			return "TIMESTAMP"
		}
		return "TIMESTAMP DEFAULT CURRENT_TIMESTAMP"
	case "uuid.UUID":
		if driver == "postgres" {
			return "UUID"
		}
		return "CHAR(36)"
	case "[]byte":
		return "BLOB"
	default:
		return "VARCHAR(255)"
	}
}

// shouldSkipField verifica se um campo deve ser ignorado
func (g *SQLGenerator) shouldSkipField(field parser.FieldMetadata) bool {
	// Ignorar campos marcados explicitamente
	if field.Ignore || (field.IgnoreRead && field.IgnoreWrite) {
		return true
	}

	// Ignorar se JSON tag é "-"
	if field.JSONTag == "-" {
		return true
	}

	// Ignorar relacionamentos que não têm coluna física
	// Relacionamentos são fields que:
	// 1. Têm Relation definida E
	// 2. Não terminam com ID/Id (foreign keys terminam com ID)
	if field.Relation != nil && !strings.HasSuffix(field.Name, "ID") && !strings.HasSuffix(field.Name, "Id") {
		return true
	}

	// Ignorar tipos complexos (structs customizados)
	// Mas permitir time.Time e uuid.UUID
	if strings.Contains(field.Type, ".") && !strings.HasPrefix(field.Type, "time.") && !strings.HasPrefix(field.Type, "uuid.") {
		return true
	}

	return false
}

// toSnakeCase converte CamelCase para snake_case
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

func (g *SQLGenerator) generateDropTable(change SchemaChange, driver string) (up string, down string) {
	up = fmt.Sprintf("DROP TABLE IF EXISTS %s;", change.TableName)
	// Down seria recriar a tabela, mas isso é complexo
	down = fmt.Sprintf("-- TODO: Recreate table %s", change.TableName)
	return up, down
}

func (g *SQLGenerator) generateAddColumn(change SchemaChange, driver string) (up string, down string) {
	// TODO: Determinar tipo SQL correto baseado no field
	up = fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s VARCHAR(255);",
		change.TableName, change.Field)

	down = fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;",
		change.TableName, change.Field)

	return up, down
}

func (g *SQLGenerator) generateDropColumn(change SchemaChange, driver string) (up string, down string) {
	up = fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;",
		change.TableName, change.Field)

	// Down seria adicionar a coluna de volta
	down = fmt.Sprintf("-- TODO: Add column %s back to %s",
		change.Field, change.TableName)

	return up, down
}

func (g *SQLGenerator) generateAlterColumn(change SchemaChange, driver string) (up string, down string) {
	up = fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s %v;",
		change.TableName, change.Field, change.NewValue)

	down = fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s %v;",
		change.TableName, change.Field, change.OldValue)

	return up, down
}

// GenerateFullTableDDL gera DDL completo para uma tabela baseado no model
func (g *SQLGenerator) GenerateFullTableDDL(metadata *parser.ModelMetadata, driver string) string {
	var ddl strings.Builder

	ddl.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", metadata.TableName))

	// Campos
	var columns []string
	var indexes []string

	for _, field := range metadata.Fields {
		col := g.generateColumnDDL(field, driver)
		columns = append(columns, "    "+col)

		// Índices
		if field.Index {
			idx := fmt.Sprintf("INDEX idx_%s_%s (%s)",
				metadata.TableName, field.Name, field.Name)
			indexes = append(indexes, "    "+idx)
		}

		if field.Unique {
			idx := fmt.Sprintf("UNIQUE INDEX uniq_%s_%s (%s)",
				metadata.TableName, field.Name, field.Name)
			indexes = append(indexes, "    "+idx)
		}
	}

	// Juntar colunas e índices
	allParts := append(columns, indexes...)
	ddl.WriteString(strings.Join(allParts, ",\n"))

	ddl.WriteString("\n)")

	// Engine e charset (MySQL)
	if driver == "mysql" {
		ddl.WriteString(" ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	}

	ddl.WriteString(";")

	return ddl.String()
}

func (g *SQLGenerator) generateColumnDDL(field parser.FieldMetadata, driver string) string {
	columnName := field.JSONTag
	if columnName == "" {
		columnName = toSnakeCase(field.Name)
	}

	sqlType := parser.GetSQLType(field.Type, driver)

	var parts []string
	parts = append(parts, columnName, sqlType)

	// Primary Key
	if field.PrimaryKey {
		if driver == "mysql" {
			parts = append(parts, "AUTO_INCREMENT PRIMARY KEY")
		} else if driver == "postgres" {
			sqlType = "SERIAL PRIMARY KEY"
			parts = []string{columnName, sqlType}
		}
		return strings.Join(parts, " ")
	}

	// NOT NULL
	if field.Required {
		parts = append(parts, "NOT NULL")
	}

	// DEFAULT
	if field.Default != "" {
		parts = append(parts, fmt.Sprintf("DEFAULT '%s'", field.Default))
	}

	return strings.Join(parts, " ")
}
