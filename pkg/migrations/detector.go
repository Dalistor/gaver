package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Dalistor/gaver/pkg/parser"
	"gorm.io/gorm"
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
	Model       *parser.ModelMetadata // Metadados completos do model
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
	// Conectar ao banco
	db, err := ConnectDB()
	if err != nil {
		// Se não conseguir conectar, assume que é banco novo
		return make(map[string]*TableSchema), nil
	}

	schema := make(map[string]*TableSchema)

	// Buscar todas as tabelas (detectar driver)
	var tables []string
	driver := os.Getenv("DB_DRIVER")

	var query string
	switch driver {
	case "postgres":
		query = "SELECT tablename FROM pg_tables WHERE schemaname = 'public'"
	case "sqlite":
		query = "SELECT name FROM sqlite_master WHERE type='table'"
	default: // mysql
		query = "SHOW TABLES"
	}

	if err := db.Raw(query).Scan(&tables).Error; err != nil {
		// Se der erro, assume banco vazio
		return make(map[string]*TableSchema), nil
	}

	// Para cada tabela, buscar colunas
	for _, tableName := range tables {
		// Ignorar tabela de controle de migrations
		if tableName == "migrations" {
			continue
		}

		columns, err := d.getTableColumns(db, tableName)
		if err != nil {
			continue
		}

		schema[tableName] = &TableSchema{
			Name:    tableName,
			Columns: columns,
		}
	}

	return schema, nil
}

// getTableColumns busca colunas de uma tabela
func (d *Detector) getTableColumns(db *gorm.DB, tableName string) ([]*ColumnSchema, error) {
	driver := os.Getenv("DB_DRIVER")

	switch driver {
	case "postgres":
		return d.getColumnsPostgres(db, tableName)
	case "sqlite":
		return d.getColumnsSQLite(db, tableName)
	default: // mysql
		return d.getColumnsMySQL(db, tableName)
	}
}

// getColumnsMySQL busca colunas usando DESCRIBE (MySQL)
func (d *Detector) getColumnsMySQL(db *gorm.DB, tableName string) ([]*ColumnSchema, error) {
	type columnInfo struct {
		Field   string
		Type    string
		Null    string
		Key     string
		Default *string
		Extra   string
	}

	var columns []*ColumnSchema
	var results []columnInfo

	if err := db.Raw(fmt.Sprintf("DESCRIBE %s", tableName)).Scan(&results).Error; err != nil {
		return nil, err
	}

	for _, col := range results {
		columns = append(columns, &ColumnSchema{
			Name:       col.Field,
			Type:       col.Type,
			Nullable:   col.Null == "YES",
			Default:    col.Default,
			PrimaryKey: col.Key == "PRI",
		})
	}

	return columns, nil
}

// getColumnsPostgres busca colunas usando information_schema (PostgreSQL)
func (d *Detector) getColumnsPostgres(db *gorm.DB, tableName string) ([]*ColumnSchema, error) {
	type columnInfo struct {
		ColumnName             string
		DataType               string
		IsNullable             string
		ColumnDefault          *string
		CharacterMaximumLength *int
	}

	var columns []*ColumnSchema
	var results []columnInfo

	query := `
		SELECT column_name, data_type, is_nullable, column_default, character_maximum_length
		FROM information_schema.columns
		WHERE table_name = ?
		ORDER BY ordinal_position
	`

	if err := db.Raw(query, tableName).Scan(&results).Error; err != nil {
		return nil, err
	}

	for _, col := range results {
		columns = append(columns, &ColumnSchema{
			Name:       col.ColumnName,
			Type:       col.DataType,
			Nullable:   col.IsNullable == "YES",
			Default:    col.ColumnDefault,
			PrimaryKey: false, // TODO: Detectar primary key
		})
	}

	return columns, nil
}

// getColumnsSQLite busca colunas usando PRAGMA (SQLite)
func (d *Detector) getColumnsSQLite(db *gorm.DB, tableName string) ([]*ColumnSchema, error) {
	type columnInfo struct {
		CID          int
		Name         string
		Type         string
		NotNull      int
		DefaultValue *string
		PK           int
	}

	var columns []*ColumnSchema
	var results []columnInfo

	if err := db.Raw(fmt.Sprintf("PRAGMA table_info(%s)", tableName)).Scan(&results).Error; err != nil {
		return nil, err
	}

	for _, col := range results {
		columns = append(columns, &ColumnSchema{
			Name:       col.Name,
			Type:       col.Type,
			Nullable:   col.NotNull == 0,
			Default:    col.DefaultValue,
			PrimaryKey: col.PK == 1,
		})
	}

	return columns, nil
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
			Model:       model, // Passar metadados completos
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
				Model:       model, // Passar metadados completos
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
		// Ignorar campos que não têm coluna física no banco
		if d.shouldSkipField(field) {
			continue
		}

		columnName := field.JSONTag
		if columnName == "" || columnName == "-" {
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

// shouldSkipField verifica se um campo deve ser ignorado na comparação
func (d *Detector) shouldSkipField(field parser.FieldMetadata) bool {
	// Ignorar campos marcados explicitamente
	if field.Ignore || (field.IgnoreRead && field.IgnoreWrite) {
		return true
	}

	// Ignorar se JSON tag é "-"
	if field.JSONTag == "-" {
		return true
	}

	// Ignorar relacionamentos que não têm coluna física
	// (não terminam com ID/Id e não são tipos básicos)
	if field.Relation != nil && !strings.HasSuffix(field.Name, "ID") && !strings.HasSuffix(field.Name, "Id") {
		return true
	}

	// Ignorar tipos complexos que não são colunas (struct, array, etc)
	if strings.Contains(field.Type, ".") && !strings.HasPrefix(field.Type, "time.") && !strings.HasPrefix(field.Type, "uuid.") {
		// É um tipo customizado (não primitivo)
		return true
	}

	return false
}

func (d *Detector) typeChanged(field parser.FieldMetadata, col *ColumnSchema) bool {
	// Não detectar mudança se for interface{} ou tipo desconhecido
	if field.Type == "interface{}" {
		return false
	}

	// Converter tipo Go para SQL esperado
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "mysql"
	}

	sqlGen := NewSQLGenerator()
	expectedType := sqlGen.goTypeToSQL(field.Type, field.GORMTag, driver)

	// Normalizar tipos para comparação
	expectedNorm := normalizeSQLType(expectedType)
	actualNorm := normalizeSQLType(col.Type)

	// Comparar tipos normalizados
	return expectedNorm != actualNorm
}

// normalizeSQLType normaliza tipos SQL para comparação
func normalizeSQLType(sqlType string) string {
	sqlType = strings.ToUpper(sqlType)

	// Remover tamanhos, defaults e constraints
	sqlType = strings.Split(sqlType, " ")[0]

	// Normalizar variações comuns
	switch {
	case strings.HasPrefix(sqlType, "VARCHAR"):
		return "VARCHAR"
	case strings.HasPrefix(sqlType, "CHAR"):
		return "CHAR"
	case strings.HasPrefix(sqlType, "BIGINT"):
		return "BIGINT"
	case strings.HasPrefix(sqlType, "INT"):
		return "INT"
	case strings.HasPrefix(sqlType, "TIMESTAMP"):
		return "TIMESTAMP"
	case strings.HasPrefix(sqlType, "DATETIME"):
		return "TIMESTAMP"
	case strings.HasPrefix(sqlType, "TINYINT"):
		return "TINYINT"
	case strings.HasPrefix(sqlType, "TEXT"):
		return "TEXT"
	case strings.HasPrefix(sqlType, "BLOB"):
		return "BLOB"
	case strings.HasPrefix(sqlType, "DOUBLE"):
		return "DOUBLE"
	case strings.HasPrefix(sqlType, "FLOAT"):
		return "FLOAT"
	default:
		return sqlType
	}
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
