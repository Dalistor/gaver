package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Migration representa um registro na tabela de controle de migrations
type Migration struct {
	ID         uint      `gorm:"primaryKey"`
	Name       string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Batch      int       `gorm:"not null"`
	ExecutedAt time.Time `gorm:"autoCreateTime"`
}

// Runner executa migrations
type Runner struct {
	migrationsPath string
}

// NewRunner cria um novo runner
func NewRunner() *Runner {
	return &Runner{
		migrationsPath: "migrations",
	}
}

// MigrationStatus representa o status das migrations
type MigrationStatus struct {
	Applied []MigrationInfo
	Pending []MigrationInfo
}

// MigrationInfo contém informações sobre uma migration
type MigrationInfo struct {
	Version     string
	Description string
	AppliedAt   *time.Time
}

// MigrateUp aplica migrations pendentes
func (r *Runner) MigrateUp(steps int) (int, error) {
	// Conectar ao banco
	if DB == nil {
		_, err := ConnectDB()
		if err != nil {
			return 0, fmt.Errorf("erro ao conectar ao banco: %w", err)
		}
	}

	// Criar tabela de migrations se não existir
	if err := createMigrationsTable(); err != nil {
		return 0, fmt.Errorf("erro ao criar tabela de migrations: %w", err)
	}

	// Buscar migrations pendentes
	pending, err := r.getPendingMigrations()
	if err != nil {
		return 0, err
	}

	if len(pending) == 0 {
		return 0, nil
	}

	toApply := pending
	if steps > 0 && steps < len(pending) {
		toApply = pending[:steps]
	}

	// Obter próximo batch number
	batch, err := getNextBatch()
	if err != nil {
		return 0, fmt.Errorf("erro ao obter batch number: %w", err)
	}

	// Aplicar migrations
	for _, migration := range toApply {
		if err := r.executeMigrationUp(migration, batch); err != nil {
			return 0, fmt.Errorf("erro ao aplicar %s: %w", migration.Version, err)
		}
	}

	return len(toApply), nil
}

// MigrateDown reverte migrations
func (r *Runner) MigrateDown(steps int) (int, error) {
	// Conectar ao banco
	if DB == nil {
		_, err := ConnectDB()
		if err != nil {
			return 0, fmt.Errorf("erro ao conectar ao banco: %w", err)
		}
	}

	// Buscar último batch
	lastBatch, err := getLastBatch()
	if err != nil {
		return 0, fmt.Errorf("nenhuma migration para reverter")
	}

	// Buscar migrations do último batch
	migrations, err := getMigrationsByBatch(lastBatch)
	if err != nil {
		return 0, fmt.Errorf("erro ao buscar migrations: %w", err)
	}

	if len(migrations) == 0 {
		return 0, fmt.Errorf("nenhuma migration para reverter")
	}

	// Aplicar steps se especificado
	toRevert := migrations
	if steps > 0 && steps < len(migrations) {
		toRevert = migrations[:steps]
	}

	// Reverter migrations (ordem reversa já vem do banco)
	for _, migration := range toRevert {
		if err := r.executeMigrationDown(migration); err != nil {
			return 0, fmt.Errorf("erro ao reverter %s: %w", migration.Name, err)
		}
	}

	return len(toRevert), nil
}

// MigrateDownTo reverte até uma versão específica
func (r *Runner) MigrateDownTo(version string) (int, error) {
	// Conectar ao banco
	if DB == nil {
		_, err := ConnectDB()
		if err != nil {
			return 0, fmt.Errorf("erro ao conectar ao banco: %w", err)
		}
	}

	// Buscar todas migrations aplicadas
	var allMigrations []Migration
	if err := DB.Order("id DESC").Find(&allMigrations).Error; err != nil {
		return 0, fmt.Errorf("erro ao buscar migrations: %w", err)
	}

	// Encontrar migrations após a versão especificada
	var toRevert []Migration
	for _, m := range allMigrations {
		if strings.HasPrefix(m.Name, version) {
			break
		}
		toRevert = append(toRevert, m)
	}

	if len(toRevert) == 0 {
		return 0, nil
	}

	// Reverter migrations
	for _, migration := range toRevert {
		if err := r.executeMigrationDown(migration); err != nil {
			return 0, fmt.Errorf("erro ao reverter %s: %w", migration.Name, err)
		}
	}

	return len(toRevert), nil
}

// GetStatus retorna o status das migrations
func (r *Runner) GetStatus() (*MigrationStatus, error) {
	applied, err := r.getAppliedMigrations()
	if err != nil {
		return nil, err
	}

	pending, err := r.getPendingMigrations()
	if err != nil {
		return nil, err
	}

	return &MigrationStatus{
		Applied: applied,
		Pending: pending,
	}, nil
}

func (r *Runner) getAppliedMigrations() ([]MigrationInfo, error) {
	if DB == nil {
		_, err := ConnectDB()
		if err != nil {
			return []MigrationInfo{}, nil // Se não conectar, assume que não há migrations
		}
	}

	var migrations []Migration
	err := DB.Order("id ASC").Find(&migrations).Error
	if err != nil {
		return []MigrationInfo{}, nil
	}

	var result []MigrationInfo
	for _, m := range migrations {
		info := r.parseMigrationFilename(m.Name)
		info.AppliedAt = &m.ExecutedAt
		result = append(result, info)
	}

	return result, nil
}

func (r *Runner) getPendingMigrations() ([]MigrationInfo, error) {
	// Conectar ao banco se ainda não conectou
	if DB == nil {
		_, _ = ConnectDB()
	}

	// Buscar migrations já aplicadas
	appliedMap := make(map[string]bool)
	if DB != nil {
		var applied []Migration
		if err := DB.Find(&applied).Error; err == nil {
			for _, m := range applied {
				appliedMap[m.Name] = true
			}
		}
	}

	// Ler arquivos de migration
	files, err := os.ReadDir(r.migrationsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []MigrationInfo{}, nil
		}
		return nil, err
	}

	var pending []MigrationInfo

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		// Verificar se já foi aplicada
		if !appliedMap[file.Name()] {
			info := r.parseMigrationFilename(file.Name())
			pending = append(pending, info)
		}
	}

	return pending, nil
}

func (r *Runner) parseMigrationFilename(filename string) MigrationInfo {
	// Formato: 20060102_150405_description.sql
	name := strings.TrimSuffix(filename, ".sql")
	parts := strings.SplitN(name, "_", 3)

	version := name
	description := ""

	if len(parts) >= 3 {
		version = parts[0] + "_" + parts[1]
		description = strings.ReplaceAll(parts[2], "_", " ")
	}

	return MigrationInfo{
		Version:     version,
		Description: description,
	}
}

func (r *Runner) executeMigrationUp(migration MigrationInfo, batch int) error {
	// Ler arquivo
	filePath := filepath.Join(r.migrationsPath, migration.Version+"_*.sql")
	matches, err := filepath.Glob(filePath)
	if err != nil || len(matches) == 0 {
		return fmt.Errorf("arquivo de migration não encontrado")
	}

	content, err := os.ReadFile(matches[0])
	if err != nil {
		return err
	}

	// Extrair parte UP
	upSQL := extractUpSQL(string(content))
	if upSQL == "" {
		return fmt.Errorf("SQL UP não encontrado na migration")
	}

	// Executar SQL no banco de dados (dividir por statements)
	if err := executeSQL(upSQL); err != nil {
		return fmt.Errorf("erro ao executar SQL: %w", err)
	}

	// Registrar migration na tabela de controle
	filename := filepath.Base(matches[0])
	if err := recordMigration(filename, batch); err != nil {
		return fmt.Errorf("erro ao registrar migration: %w", err)
	}

	fmt.Printf("✓ Aplicado: %s\n", filename)

	return nil
}

func extractUpSQL(content string) string {
	// Procurar por -- ========== UP ==========
	upMarker := "-- ========== UP =========="
	downMarker := "-- ========== DOWN =========="

	upIdx := strings.Index(content, upMarker)
	downIdx := strings.Index(content, downMarker)

	if upIdx == -1 {
		return ""
	}

	start := upIdx + len(upMarker)
	end := len(content)
	if downIdx > upIdx {
		end = downIdx
	}

	return strings.TrimSpace(content[start:end])
}

func extractDownSQL(content string) string {
	// Procurar por -- ========== DOWN ==========
	downMarker := "-- ========== DOWN =========="

	downIdx := strings.Index(content, downMarker)
	if downIdx == -1 {
		return ""
	}

	return strings.TrimSpace(content[downIdx+len(downMarker):])
}

// executeMigrationDown executa o SQL DOWN de uma migration
func (r *Runner) executeMigrationDown(migration Migration) error {
	// Encontrar arquivo da migration
	filePath := filepath.Join(r.migrationsPath, migration.Name)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	// Extrair parte DOWN
	downSQL := extractDownSQL(string(content))
	if downSQL == "" {
		return fmt.Errorf("SQL DOWN não encontrado na migration")
	}

	// Executar SQL no banco de dados (dividir por statements)
	if err := executeSQL(downSQL); err != nil {
		return fmt.Errorf("erro ao executar SQL: %w", err)
	}

	// Remover registro da tabela de migrations
	if err := removeMigrationRecord(migration.Name); err != nil {
		return fmt.Errorf("erro ao remover registro: %w", err)
	}

	fmt.Printf("✓ Revertido: %s\n", migration.Name)

	return nil
}

// createMigrationsTable cria a tabela de controle se não existir
func createMigrationsTable() error {
	if DB == nil {
		return fmt.Errorf("banco de dados não conectado")
	}
	return DB.AutoMigrate(&Migration{})
}

// recordMigration registra uma migration executada
func recordMigration(name string, batch int) error {
	if DB == nil {
		return fmt.Errorf("banco de dados não conectado")
	}

	migration := &Migration{
		Name:  name,
		Batch: batch,
	}

	return DB.Create(migration).Error
}

// getNextBatch retorna o próximo número de batch
func getNextBatch() (int, error) {
	if DB == nil {
		return 0, fmt.Errorf("banco de dados não conectado")
	}

	var lastMigration Migration
	result := DB.Order("batch DESC").First(&lastMigration)

	if result.Error != nil {
		// Se não há registros, retorna batch 1
		return 1, nil
	}

	return lastMigration.Batch + 1, nil
}

// getLastBatch retorna o número do último batch
func getLastBatch() (int, error) {
	if DB == nil {
		return 0, fmt.Errorf("banco de dados não conectado")
	}

	var lastMigration Migration
	result := DB.Order("batch DESC").First(&lastMigration)

	if result.Error != nil {
		return 0, result.Error
	}

	return lastMigration.Batch, nil
}

// getMigrationsByBatch retorna migrations de um batch específico
func getMigrationsByBatch(batch int) ([]Migration, error) {
	if DB == nil {
		return nil, fmt.Errorf("banco de dados não conectado")
	}

	var migrations []Migration
	err := DB.Where("batch = ?", batch).Order("id DESC").Find(&migrations).Error
	return migrations, err
}

// removeMigrationRecord remove um registro de migration
func removeMigrationRecord(name string) error {
	if DB == nil {
		return fmt.Errorf("banco de dados não conectado")
	}

	return DB.Where("name = ?", name).Delete(&Migration{}).Error
}

// executeSQL executa SQL dividindo em statements individuais
func executeSQL(sql string) error {
	if DB == nil {
		return fmt.Errorf("banco de dados não conectado")
	}

	// Dividir SQL em statements individuais (separados por ;)
	statements := splitSQLStatements(sql)

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		// Executar statement
		if err := DB.Exec(stmt).Error; err != nil {
			return fmt.Errorf("erro ao executar statement: %w\nSQL: %s", err, stmt)
		}
	}

	return nil
}

// splitSQLStatements divide SQL em statements individuais
func splitSQLStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	inString := false
	var stringChar rune

	for i, char := range sql {
		current.WriteRune(char)

		// Detectar strings
		if (char == '\'' || char == '"') && (i == 0 || sql[i-1] != '\\') {
			if inString && char == stringChar {
				inString = false
			} else if !inString {
				inString = true
				stringChar = char
			}
		}

		// Detectar fim de statement (;)
		if char == ';' && !inString {
			statements = append(statements, current.String())
			current.Reset()
		}
	}

	// Adicionar último statement se não terminar com ;
	if current.Len() > 0 {
		statements = append(statements, current.String())
	}

	return statements
}
