package database

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// SupabaseConnection represents a connection to Supabase PostgreSQL database
type SupabaseConnection struct {
	DB *sqlx.DB
}

// ConnectSupabase establishes a connection to Supabase PostgreSQL
// and automatically runs migrations at startup
func ConnectSupabase() (*SupabaseConnection, error) {
	log.Println("Connecting to Supabase database...")

	// Support both explicit env vars and 1Password CLI format
	dbPassword := os.Getenv("SUPABASE_DB_PASSWORD")
	projectRef := os.Getenv("SUPABASE_PROJECT_REF")
	postgresURL := os.Getenv("POSTGRES_URL")

	var connectionString string

	if postgresURL != "" {
		// Use direct POSTGRES_URL if available (1Password CLI format)
		connectionString = postgresURL
		log.Println("Using POSTGRES_URL from environment")
	} else if dbPassword != "" && projectRef != "" {
		// Build connection string from components
		connectionString = fmt.Sprintf(
			"postgres://postgres:%s@db.%s.supabase.co:5432/postgres?sslmode=require",
			dbPassword,
			projectRef,
		)
	} else {
		// Try 1Password format
		dbPassword = os.Getenv("SUPABASE_PRIV_KEY")
		if dbPassword != "" && projectRef != "" {
			connectionString = fmt.Sprintf(
				"postgres://postgres:%s@db.%s.supabase.co:5432/postgres?sslmode=require",
				dbPassword,
				projectRef,
			)
		} else {
			return nil, fmt.Errorf("database credentials not found: set POSTGRES_URL or (SUPABASE_DB_PASSWORD and SUPABASE_PROJECT_REF)")
		}
	}

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Supabase database: %w", err)
	}

	connection := &SupabaseConnection{DB: db}

	if err := connection.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging Supabase database: %w", err)
	}

	log.Println("Connected to Supabase database")

	// Run migrations automatically at startup
	log.Println("Running database migrations...")
	if err := connection.runMigrations(); err != nil {
		return nil, fmt.Errorf("error running migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return connection, nil
}

// runMigrations executes embedded migrations using goose
func (c *SupabaseConnection) runMigrations() error {
	// Set up goose with embedded migrations
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("error setting goose dialect: %w", err)
	}

	// Run all pending migrations
	if err := goose.Up(c.DB.DB, "migrations"); err != nil {
		return fmt.Errorf("error running goose migrations: %w", err)
	}

	return nil
}

// Close closes the Supabase database connection
func (c *SupabaseConnection) Close() error {
	log.Println("Closing Supabase database connection...")
	return c.DB.Close()
}

// Ping verifies the connection to Supabase is alive
func (c *SupabaseConnection) Ping() error {
	log.Println("Pinging Supabase database connection...")
	return c.DB.Ping()
}

// HealthCheck performs a basic health check on the database
func (c *SupabaseConnection) HealthCheck() error {
	var result int
	err := c.DB.Get(&result, "SELECT 1")
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	return nil
}
