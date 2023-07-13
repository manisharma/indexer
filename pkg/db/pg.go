package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"indexer/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

const migrationVersion = 1

//go:embed migrations/*.sql
var files embed.FS

// connects to postgres & returns a pool
func Connect(ctx context.Context, cfg config.PgCfg) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, cfg.String())
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func Migrate(cfg config.PgCfg) error {
	driver, err := iofs.New(files, "migrations")
	if err != nil {
		return fmt.Errorf("load migration files: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", driver, cfg.String())
	if err != nil {
		return fmt.Errorf("mount iofs instance: %w", err)
	}
	if err := m.Migrate(migrationVersion); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
