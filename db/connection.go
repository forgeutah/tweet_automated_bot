package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Connection struct {
	DB *sqlx.DB
}

func Connect(ctx context.Context) (*Connection, error) {
	log.Println("Connecting to database...")

	rootCertPath, err := loadCockroachRootCert(ctx)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("sslmode", "verify-full")
	params.Set("sslrootcert", rootCertPath)
	params.Set("options", os.Getenv("DB_OPTIONS"))

	connectionString := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD")),
		Host:     os.Getenv("DB_HOST"),
		Path:     os.Getenv("DB_NAME"),
		RawQuery: params.Encode(),
	}
	log.Println(connectionString.String())

	db, err := sqlx.Connect("postgres", connectionString.String())
	if err != nil {

		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	connection := &Connection{DB: db}
	connection.Migrate(ctx)

	if err := connection.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	log.Println("Connected to database")

	return connection, nil
}

func (c *Connection) Close(ctx context.Context) error {
	log.Println("Closing database connection...")
	err := removeCert(ctx)
	if err != nil {
		log.Println(fmt.Errorf("error removing cert: %w", err))
	}
	return c.DB.Close()
}

func (c *Connection) Ping() error {
	log.Println("Pinging database connection...")
	return c.DB.Ping()
}

func (c *Connection) Migrate(ctx context.Context) {
	log.Println("Migrating database...")

	c.DB.MustExecContext(ctx, create_query)
	// check if table exists
	var count int
	row := c.DB.QueryRowx("SELECT count(*) FROM yt_videos LIMIT 1")
	if row == nil {
		log.Println("Table does not exist")
		// insert data
		result := c.DB.MustExecContext(ctx, videoInsert)
		fmt.Println(result)
	} else {
		row.Scan(&count)
		if count == 0 {
			log.Println("Table is empty")
			// insert data
			result := c.DB.MustExecContext(ctx, videoInsert)
			fmt.Println(result)
		} else {
			log.Println("Table is not empty")
		}
	}
}
