package migrations

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB é a conexão global com o banco de dados para migrations
var DB *gorm.DB

// ConnectDB conecta ao banco de dados usando variáveis de ambiente
func ConnectDB() (*gorm.DB, error) {
	// Carregar .env se existir
	_ = godotenv.Load()

	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "mysql" // Default
	}

	dsn := buildDSN(driver)
	if dsn == "" {
		return nil, fmt.Errorf("não foi possível construir DSN para driver: %s", driver)
	}

	var dialector gorm.Dialector
	switch driver {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(dsn)
	default:
		return nil, fmt.Errorf("driver não suportado: %s", driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco: %w", err)
	}

	DB = db
	return db, nil
}

// buildDSN constrói a string de conexão baseada no driver
func buildDSN(driver string) string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	switch driver {
	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode,
		)
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, password, host, port, dbname,
		)
	case "sqlite":
		if dbname == "" {
			return "gaver.db"
		}
		return dbname + ".db"
	default:
		return ""
	}
}

// CloseDB fecha a conexão com o banco
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
