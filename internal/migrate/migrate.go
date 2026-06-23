package migrate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Up(ctx context.Context, db *pgxpool.Pool, dir string) error {
	if err := ensureSchemaMigrations(ctx, db); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {
		version, err := parseVersion(file)
		if err != nil {
			return err
		}

		applied, err := migrationApplied(ctx, db, version)
		if err != nil {
			return err
		}

		if applied {
			continue
		}

		if err := applyMigration(ctx, db, file, version); err != nil {
			return err
		}
	}

	return nil
}

func ensureSchemaMigrations(ctx context.Context, db *pgxpool.Pool) error {
	const query = `
CREATE TABLE IF NOT EXISTS schema_migrations(
    version BIGINT PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
)
`
	_, err := db.Exec(ctx, query)
	return err
}

func migrationApplied(ctx context.Context, db *pgxpool.Pool, version int64) (bool, error) {
	const query = `
SELECT EXISTS(
    SELECT 1
    FROM schema_migrations
    WHERE version=$1
)
`
	var exists bool
	err := db.QueryRow(ctx, query, version).Scan(&exists)

	return exists, err
}

func applyMigration(ctx context.Context, db *pgxpool.Pool, file string, version int64) error {
	sqlBytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
		return fmt.Errorf("%s: %w", file, err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES($1) `, version)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func parseVersion(path string) (int64, error) {
	name := filepath.Base(path)
	parts := strings.Split(name, "_")
	return strconv.ParseInt(parts[0], 10, 64)
}
