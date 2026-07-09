package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AutoMigrate checks if the database schema exists. If not, it executes the
// SQL migration file at sqlPath to create all tables. It never drops or
// recreates existing tables — it is idempotent.
func AutoMigrate(ctx context.Context, pool *pgxpool.Pool, sqlPath string) error {
	exists, err := schemaExists(ctx, pool)
	if err != nil {
		return fmt.Errorf("check schema exists: %w", err)
	}

	if exists {
		log.Println("Database schema already exists — skipping migration")
		return nil
	}

	log.Printf("Database schema not found — running migration from %s", sqlPath)

	data, err := os.ReadFile(sqlPath)
	if err != nil {
		return fmt.Errorf("read migration file %s: %w", sqlPath, err)
	}

	// Split on semicolons to execute each statement separately.
	// pgx.Exec does not support multi-statement strings in all drivers,
	// so we split and execute one by one.
	statements := splitSQL(string(data))
	for i, stmt := range statements {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}
		// Skip meta-commands like \connect — those are psql-specific
		if strings.HasPrefix(trimmed, `\`) {
			log.Printf("  Skipping psql meta-command: %s", strings.SplitN(trimmed, "\n", 2)[0])
			continue
		}
		if _, err := pool.Exec(ctx, trimmed); err != nil {
			return fmt.Errorf("execute statement %d: %w\nSQL: %s", i+1, err, trimForLog(trimmed, 120))
		}
	}

	log.Println("Database migration complete")
	return nil
}

// schemaExists checks whether the public_calls_proposals table exists.
func schemaExists(ctx context.Context, pool *pgxpool.Pool) (bool, error) {
	query := `SELECT EXISTS (
		SELECT FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = 'public_calls_proposals'
	)`
	var exists bool
	if err := pool.QueryRow(ctx, query).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// splitSQL splits a SQL script on semicolons, preserving statement boundaries.
func splitSQL(data string) []string {
	var parts []string
	var current strings.Builder
	inString := false
	var stringChar byte

	for i := 0; i < len(data); i++ {
		ch := data[i]
		if inString {
			current.WriteByte(ch)
			if ch == stringChar && (i+1 >= len(data) || data[i+1] != stringChar) {
				inString = false
			}
			continue
		}
		if ch == '\'' || ch == '"' {
			inString = true
			stringChar = ch
			current.WriteByte(ch)
			continue
		}
		if ch == '-' && i+1 < len(data) && data[i+1] == '-' {
			// Single-line comment — skip to end of line
			for i < len(data) && data[i] != '\n' {
				i++
			}
			continue
		}
		if ch == '/' && i+1 < len(data) && data[i+1] == '*' {
			// Block comment — skip to */
			i += 2
			for i+1 < len(data) && !(data[i] == '*' && data[i+1] == '/') {
				i++
			}
			i++ // skip the /
			continue
		}
		if ch == ';' {
			parts = append(parts, current.String())
			current.Reset()
			continue
		}
		current.WriteByte(ch)
	}
	// Last statement (may not have trailing semicolon)
	remaining := strings.TrimSpace(current.String())
	if remaining != "" {
		parts = append(parts, remaining)
	}
	return parts
}

func trimForLog(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}