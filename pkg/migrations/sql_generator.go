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
	// TODO: Buscar metadata do model para gerar DDL completo
	// Por enquanto, gera uma estrutura básica

	up = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`, change.TableName)

	down = fmt.Sprintf("DROP TABLE IF EXISTS %s;", change.TableName)

	return up, down
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
