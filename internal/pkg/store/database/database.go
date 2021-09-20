package database

import (
	"context"
	"fmt"

	"github.com/SKorolchuk/dpio-workspace/internal/pkg/infra/database"

	"github.com/ardanlabs/darwin"
	"github.com/jmoiron/sqlx"
)

var (
	// go:embed sql/schema.sql
	workspaceSchemaScript string

	// go:embed sql/seed.sql
	workspaceSeedScript string

	// go:embed sql/drop.sql
	workspaceDropScript string
)

// Migrate will do the schema migration of workspace database.
func Migrate(ctx context.Context, connection *sqlx.DB) error {
	if err := database.StatusCheck(ctx, connection); err != nil {
		return fmt.Errorf("database is not available: %w", err)
	}

	driver, err := darwin.NewGenericDriver(connection.DB, darwin.PostgresDialect{})
	if err != nil {
		return fmt.Errorf("failed to initialize migration driver: %w", err)
	}

	processor := darwin.New(driver, darwin.ParseMigrations(workspaceSchemaScript))
	return processor.Migrate()
}

// Seed will generate initial data useful for development and testing purposes.
func Seed(ctx context.Context, connection *sqlx.DB) error {
	if err := database.StatusCheck(ctx, connection); err != nil {
		return fmt.Errorf("database is not available: %w", err)
	}

	transaction, err := connection.Begin()
	if err != nil {
		return err
	}

	if _, err := transaction.Exec(workspaceSeedScript); err != nil {
		if err := transaction.Rollback(); err != nil {
			return err
		}
	}

	return transaction.Commit()
}

// Drop cleans workspace database.
func Drop(ctx context.Context, connection *sqlx.DB) error {
	if err := database.StatusCheck(ctx, connection); err != nil {
		return fmt.Errorf("database is not available: %w", err)
	}

	transaction, err := connection.Begin()
	if err != nil {
		return err
	}

	if _, err := transaction.Exec(workspaceDropScript); err != nil {
		if err := transaction.Rollback(); err != nil {
			return err
		}
	}

	return transaction.Commit()
}
