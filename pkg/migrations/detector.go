package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Dalistor/gaver/pkg/parser"
)

// SchemaChange representa uma mudança no schema
type SchemaChange struct {
	Type        string // CREATE_TABLE, ADD_COLUMN, DROP_COLUMN, ALTER_COLUMN, etc
	ModelName   string
	TableName   string
	Field       string
	Description string
	OldValue    interface{}
	NewValue    interface{}
}

// Detector detecta mudanças nos models
type Detector struct {
	modelsPath string
}

// NewDetector cria um novo detector
func NewDetector() *Detector {
	return &Detector{
		modelsPath: "modules",
	}
}

// DetectChanges escaneia os models e detecta mudanças
func (d *Detector) DetectChanges() ([]SchemaChange, error) {
	// 1. Escanear todos os models
	models, err := d.scanModels()
	if err != nil {
		return nil, err
	}

	if len(models) == 0 {
		return []SchemaChange{}, nil
	}

	// 2. Ler schema atual do banco (se existir)
	dbSchema, err := d.readDatabaseSchema()
	if err != nil {
		// Se o banco não existe ou não está configurado, todas as tabelas são novas
		return d.allTablesAsNew(models), nil
	}

	// 3. Comparar e detectar mudanças
	changes := d.compareSchemas(models, dbSchema)

	return changes, nil
}

// scanModels escaneia todos os arquivos de models
func (d *Detector) scanModels() ([]*parser.ModelMetadata, error) {
	var models []*parser.ModelMetadata

	// Escanear pasta modules/*/models/*.go
	err := filepath.Walk(d.modelsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignorar diretórios e arquivos que não são .go
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		// Verificar se está em uma pasta models
		dir := filepath.Dir(path)
		if filepath.Base(dir) != "models" {
			return nil
		}

		// Parsear o arquivo
		metadata, err := parser.ParseModelFile(path)
		if err != nil {
			// Ignorar arquivos que não podem ser parseados
			return nil
		}

		models = append(models, metadata)
		return nil
	})

	return models, err
}

// readDatabaseSchema lê o schema atual do banco
func (d *Detector) readDatabaseSchema() (map[string]*TableSchema, error) {
	// TODO: Implementar leitura real do banco de dados
	// Por enquanto, retorna vazio para simular banco novo
	return make(map[string]*TableSchema), nil
}

// allTablesAsNew cria mudanças CREATE_TABLE para todos os models
func (d *Detector) allTablesAsNew(models []*parser.ModelMetadata) []SchemaChange {
	changes := []SchemaChange{}

	for _, model := range models {
		changes = append(changes, SchemaChange{
			Type:        "CREATE_TABLE",
			ModelName:   model.Name,
			TableName:   model.TableName,
			Description: fmt.Sprintf("Criar tabela %s", model.TableName),
		})
	}

	return changes
}

// compareSchemas compara models com schema do banco
func (d *Detector) compareSchemas(models []*parser.ModelMetadata, dbSchema map[string]*TableSchema) []SchemaChange {
	changes := []SchemaChange{}

	// Verificar models novos e alterados
	for _, model := range models {
		table, exists := dbSchema[model.TableName]

		if !exists {
			// Tabela nova
			changes = append(changes, SchemaChange{
				Type:        "CREATE_TABLE",
				ModelName:   model.Name,
				TableName:   model.TableName,
				Description: fmt.Sprintf("Criar tabela %s", model.TableName),
			})
		} else {
			// Comparar campos
			fieldChanges := d.compareFields(model, table)
			changes = append(changes, fieldChanges...)
		}
	}

	// TODO: Detectar tabelas removidas (DROP_TABLE)

	return changes
}

// compareFields compara campos do model com a tabela do banco
func (d *Detector) compareFields(model *parser.ModelMetadata, table *TableSchema) []SchemaChange {
	changes := []SchemaChange{}

	// Criar mapa de colunas do banco
	dbColumns := make(map[string]*ColumnSchema)
	for _, col := range table.Columns {
		dbColumns[col.Name] = col
	}

	// Verificar campos novos e alterados
	for _, field := range model.Fields {
		columnName := field.JSONTag
		if columnName == "" {
			columnName = toSnakeCase(field.Name)
		}

		dbCol, exists := dbColumns[columnName]

		if !exists {
			// Coluna nova
			changes = append(changes, SchemaChange{
				Type:        "ADD_COLUMN",
				ModelName:   model.Name,
				TableName:   model.TableName,
				Field:       columnName,
				Description: fmt.Sprintf("Adicionar coluna %s em %s", columnName, model.TableName),
			})
		} else {
			// Verificar se o tipo mudou
			if d.typeChanged(field, dbCol) {
				changes = append(changes, SchemaChange{
					Type:        "ALTER_COLUMN",
					ModelName:   model.Name,
					TableName:   model.TableName,
					Field:       columnName,
					Description: fmt.Sprintf("Alterar tipo da coluna %s em %s", columnName, model.TableName),
					OldValue:    dbCol.Type,
					NewValue:    field.Type,
				})
			}
		}
	}

	// TODO: Detectar colunas removidas (DROP_COLUMN)

	return changes
}

func (d *Detector) typeChanged(field parser.FieldMetadata, col *ColumnSchema) bool {
	// TODO: Implementar comparação de tipos
	return false
}

// GenerateMigrationFile gera arquivo de migration
func (d *Detector) GenerateMigrationFile(changes []SchemaChange, name string) (string, error) {
	// Gerar timestamp
	timestamp := time.Now().Format("20060102_150405")

	// Nome do arquivo
	if name == "" {
		name = "auto_migration"
	}
	filename := fmt.Sprintf("%s_%s.sql", timestamp, name)
	filepath := filepath.Join("migrations", filename)

	// Gerar conteúdo SQL
	sqlGenerator := NewSQLGenerator()
	upSQL, downSQL := sqlGenerator.Generate(changes, "mysql") // TODO: Detectar driver

	// Criar conteúdo do arquivo
	content := fmt.Sprintf(`-- Migration: %s
-- Generated at: %s

-- ========== UP ==========
%s

-- ========== DOWN ==========
%s
`, filename, time.Now().Format("2006-01-02 15:04:05"), upSQL, downSQL)

	// Criar pasta migrations se não existir
	os.MkdirAll("migrations", 0755)

	// Salvar arquivo
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

// TableSchema representa o schema de uma tabela no banco
type TableSchema struct {
	Name    string
	Columns []*ColumnSchema
	Indexes []*IndexSchema
}

// ColumnSchema representa uma coluna no banco
type ColumnSchema struct {
	Name       string
	Type       string
	Nullable   bool
	Default    interface{}
	PrimaryKey bool
	Unique     bool
}

// IndexSchema representa um índice no banco
type IndexSchema struct {
	Name    string
	Columns []string
	Unique  bool
}
