package commands

import (
	"fmt"

	"github.com/Dalistor/gaver/pkg/migrations"

	"github.com/spf13/cobra"
)

func NewMigrationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "makemigrations",
		Short: "Detecta mudanças nos models e gera migrations",
		Long:  "Escaneia os models, compara com o schema do banco e gera arquivos SQL de migration.",
		RunE:  runMakeMigrations,
	}

	cmd.Flags().StringP("name", "n", "", "Nome descritivo para a migration")
	cmd.Flags().BoolP("dry-run", "d", false, "Apenas mostra as mudanças sem gerar arquivo")

	return cmd
}

func runMakeMigrations(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	fmt.Println("Detectando mudanças nos models...")

	detector := migrations.NewDetector()

	changes, err := detector.DetectChanges()
	if err != nil {
		return fmt.Errorf("erro ao detectar mudanças: %w", err)
	}

	if len(changes) == 0 {
		fmt.Println("✓ Nenhuma mudança detectada")
		return nil
	}

	fmt.Printf("Encontradas %d mudança(s):\n\n", len(changes))
	for _, change := range changes {
		fmt.Printf("  - %s: %s\n", change.Type, change.Description)
	}

	if dryRun {
		fmt.Println("\n(dry-run: nenhum arquivo foi gerado)")
		return nil
	}

	// Gerar arquivo de migration
	migrationFile, err := detector.GenerateMigrationFile(changes, name)
	if err != nil {
		return fmt.Errorf("erro ao gerar migration: %w", err)
	}

	fmt.Printf("\n✓ Migration gerada: %s\n", migrationFile)
	return nil
}

func NewMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Gerencia migrations do banco de dados",
		Long:  "Executa, reverte e mostra status das migrations.",
	}

	cmd.AddCommand(newMigrateUpCommand())
	cmd.AddCommand(newMigrateDownCommand())
	cmd.AddCommand(newMigrateStatusCommand())

	return cmd
}

func newMigrateUpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Aplica migrations pendentes",
		Long:  "Executa todas as migrations que ainda não foram aplicadas ao banco de dados.",
		RunE:  runMigrateUp,
	}

	cmd.Flags().IntP("steps", "s", 0, "Número de migrations para aplicar (0 = todas)")

	return cmd
}

func runMigrateUp(cmd *cobra.Command, args []string) error {
	steps, _ := cmd.Flags().GetInt("steps")

	fmt.Println("Aplicando migrations...")

	runner := migrations.NewRunner()

	applied, err := runner.MigrateUp(steps)
	if err != nil {
		return fmt.Errorf("erro ao aplicar migrations: %w", err)
	}

	if applied == 0 {
		fmt.Println("✓ Nenhuma migration pendente")
	} else {
		fmt.Printf("✓ %d migration(s) aplicada(s) com sucesso\n", applied)
	}

	return nil
}

func newMigrateDownCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "Reverte migrations",
		Long:  "Reverte a última migration ou um número específico de migrations.",
		RunE:  runMigrateDown,
	}

	cmd.Flags().IntP("steps", "s", 1, "Número de migrations para reverter")
	cmd.Flags().StringP("to", "t", "", "Reverter até uma versão específica")

	return cmd
}

func runMigrateDown(cmd *cobra.Command, args []string) error {
	steps, _ := cmd.Flags().GetInt("steps")
	toVersion, _ := cmd.Flags().GetString("to")

	fmt.Println("Revertendo migrations...")

	runner := migrations.NewRunner()

	var reverted int
	var err error

	if toVersion != "" {
		reverted, err = runner.MigrateDownTo(toVersion)
	} else {
		reverted, err = runner.MigrateDown(steps)
	}

	if err != nil {
		return fmt.Errorf("erro ao reverter migrations: %w", err)
	}

	fmt.Printf("✓ %d migration(s) revertida(s) com sucesso\n", reverted)
	return nil
}

func newMigrateStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Mostra status das migrations",
		Long:  "Lista todas as migrations e indica quais foram aplicadas.",
		RunE:  runMigrateStatus,
	}
}

func runMigrateStatus(cmd *cobra.Command, args []string) error {
	runner := migrations.NewRunner()

	status, err := runner.GetStatus()
	if err != nil {
		return fmt.Errorf("erro ao obter status: %w", err)
	}

	fmt.Println("\n=== Status das Migrations ===")

	if len(status.Applied) == 0 {
		fmt.Println("Aplicadas: (nenhuma)")
	} else {
		fmt.Println("Aplicadas:")
		for _, m := range status.Applied {
			fmt.Printf("  ✓ %s - %s\n", m.Version, m.Description)
		}
	}

	fmt.Println()

	if len(status.Pending) == 0 {
		fmt.Println("Pendentes: (nenhuma)")
	} else {
		fmt.Println("Pendentes:")
		for _, m := range status.Pending {
			fmt.Printf("  • %s - %s\n", m.Version, m.Description)
		}
	}

	fmt.Println()
	return nil
}
