package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

type Connection struct {
	DB *sqlx.DB
}

func Connect(ctx context.Context) (*Connection, error) {
	log.Println("Connecting to database...")
	db, err := sqlx.ConnectContext(ctx, "postgres", os.Getenv("COCKROACH_DB_URL"))
	if err != nil {

		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	connection := &Connection{DB: db}
	connection.Migrate(ctx)

	if err := connection.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	return connection, nil
}

func (c *Connection) Close(ctx context.Context) error {
	log.Println("Closing database connection...")
	c.DB.MustExecContext(ctx, `DROP TABLE IF EXISTS yt_videos CASCADE`)
	return c.DB.Close()
}

func (c *Connection) Ping() error {
	log.Println("Pinging database connection...")
	return c.DB.Ping()
}

func (c *Connection) Migrate(ctx context.Context) {
	log.Println("Migrating database...")

	c.DB.MustExecContext(ctx, create_query)
}
