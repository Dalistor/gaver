package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
	// TODO: Implementar execução real de migrations
	// Por enquanto, apenas simula

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

	for _, migration := range toApply {
		if err := r.executeMigrationUp(migration); err != nil {
			return 0, fmt.Errorf("erro ao aplicar %s: %w", migration.Version, err)
		}
	}

	return len(toApply), nil
}

// MigrateDown reverte migrations
func (r *Runner) MigrateDown(steps int) (int, error) {
	// TODO: Implementar reversão real de migrations
	return 0, fmt.Errorf("não implementado ainda")
}

// MigrateDownTo reverte até uma versão específica
func (r *Runner) MigrateDownTo(version string) (int, error) {
	// TODO: Implementar reversão até versão
	return 0, fmt.Errorf("não implementado ainda")
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
	// TODO: Buscar do banco de dados
	// Por enquanto, retorna vazio
	return []MigrationInfo{}, nil
}

func (r *Runner) getPendingMigrations() ([]MigrationInfo, error) {
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

		info := r.parseMigrationFilename(file.Name())
		pending = append(pending, info)
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

func (r *Runner) executeMigrationUp(migration MigrationInfo) error {
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

	// TODO: Executar SQL no banco de dados
	// Por enquanto, apenas mostra o que seria executado
	fmt.Printf("Executando: %s\n", migration.Version)

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

