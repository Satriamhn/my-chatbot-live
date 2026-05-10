package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"my-chatbot-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection and runs migrations
func InitDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default fallback for local development if not provided
		dsn = "host=localhost user=postgres password=postgres dbname=chatbot port=5432 sslmode=disable"
	}

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Connection pooling settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// RunMigrations handles GORM auto-migrations and custom SQL migrations
func RunMigrations(db *gorm.DB) error {
	// 0. Enable pgvector extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		log.Printf("Warning: failed to create vector extension: %v", err)
		// Proceeding anyway, maybe DB admin needs to do it
	}

	// 1. Auto Migrate models
	err := db.AutoMigrate(
		&models.Organization{},
		&models.User{},
		&models.BotSetting{},
		&models.KnowledgeBase{},
		&models.KnowledgeBaseItem{},
		&models.ChatSession{},
		&models.Message{},
	)
	if err != nil {
		return err
	}

	// 1.5. Drop NOT NULL constraint on chat_sessions.user_id to allow anonymous widget sessions
	if err := db.Exec("ALTER TABLE chat_sessions ALTER COLUMN user_id DROP NOT NULL").Error; err != nil {
		log.Printf("Warning: failed to drop not null constraint on chat_sessions.user_id: %v", err)
	}

	// 2. Enable RLS on tenant-scoped tables
	tenantTables := []string{
		"users",
		"bot_settings",
		"knowledge_bases",
		"knowledge_base_items",
		"chat_sessions",
	}

	for _, table := range tenantTables {
		if err := setupRLS(db, table); err != nil {
			log.Printf("Warning: failed to setup RLS for %s: %v", table, err)
			return err
		}
	}

	return nil
}

func setupRLS(db *gorm.DB, tableName string) error {
	// Enable RLS on the table
	if err := db.Exec(fmt.Sprintf("ALTER TABLE %s ENABLE ROW LEVEL SECURITY", tableName)).Error; err != nil {
		return err
	}

	// Create a policy that restricts access based on app.current_tenant setting
	// Note: We use a session variable 'app.current_tenant' to pass the organization_id
	policyName := fmt.Sprintf("tenant_isolation_policy_%s", tableName)

	// Drop existing policy if it exists to make this idempotent
	db.Exec(fmt.Sprintf("DROP POLICY IF EXISTS %s ON %s", policyName, tableName))

	sql := fmt.Sprintf(`
		CREATE POLICY %s ON %s
		USING (organization_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid)
		WITH CHECK (organization_id = NULLIF(current_setting('app.current_tenant', true), '')::uuid)
	`, policyName, tableName)

	return db.Exec(sql).Error
}

// TenantScope returns a GORM scope that filters by organization_id
// and sets the Postgres session variable for RLS enforcement.
func TenantScope(orgID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if orgID == "" {
			return db
		}
		// Set session variable for RLS
		db.Exec(fmt.Sprintf("SET LOCAL app.current_tenant = '%s'", orgID))
		return db.Where("organization_id = ?", orgID)
	}
}
