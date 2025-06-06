package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func InitDB(ctx context.Context, dbPath string, pragmasSQL string) (*sql.DB, error) {
	sqlDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if _, err := sqlDB.ExecContext(ctx, pragmasSQL); err != nil {
		return nil, fmt.Errorf("failed to execute pragmas: %w", err)
	}

	return sqlDB, nil
}
