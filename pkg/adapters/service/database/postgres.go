package database

import (
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"strings"
)

// PostgresAdapter manages PostgreSQL databases, roles, and dumps.
type PostgresAdapter struct {
	adminDSN string
}

func NewPostgresAdapter(adminDSN string) *PostgresAdapter {
	return &PostgresAdapter{
		adminDSN: adminDSN,
	}
}

func (a *PostgresAdapter) Engine() string {
	return "postgres"
}

func (a *PostgresAdapter) CreateDatabase(ctx context.Context, dbName, collation string) error {
	if strings.ContainsAny(dbName, " ';`\"\\") {
		return fmt.Errorf("invalid characters in database name: %s", dbName)
	}

	db, err := sql.Open("postgres", a.adminDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)
	_, err = db.ExecContext(ctx, query)
	return err
}

func (a *PostgresAdapter) DropDatabase(ctx context.Context, dbName string) error {
	if strings.ContainsAny(dbName, " ';`\"\\") {
		return fmt.Errorf("invalid characters in database name: %s", dbName)
	}

	db, err := sql.Open("postgres", a.adminDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, dbName)
	_, err = db.ExecContext(ctx, query)
	return err
}

func (a *PostgresAdapter) CreateUser(ctx context.Context, username, password, host string) error {
	if strings.ContainsAny(username, " ';`\"\\") {
		return fmt.Errorf("invalid characters in username: %s", username)
	}

	db, err := sql.Open("postgres", a.adminDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	// PostgreSQL uses ROLE with PASSWORD
	query := fmt.Sprintf(`CREATE ROLE "%s" WITH LOGIN PASSWORD '%s'`, username, strings.ReplaceAll(password, "'", "''"))
	_, err = db.ExecContext(ctx, query)
	return err
}

func (a *PostgresAdapter) GrantPermissions(ctx context.Context, username, dbName, host string, privileges []string) error {
	if strings.ContainsAny(username, " ';`\"\\") || strings.ContainsAny(dbName, " ';`\"\\") {
		return fmt.Errorf("invalid characters in identifiers")
	}

	db, err := sql.Open("postgres", a.adminDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`GRANT ALL PRIVILEGES ON DATABASE "%s" TO "%s"`, dbName, username)
	_, err = db.ExecContext(ctx, query)
	return err
}

func (a *PostgresAdapter) DumpDatabase(ctx context.Context, dbName, targetPath string) error {
	cmd := exec.CommandContext(ctx, "pg_dump", "-Fc", "-b", "-v", "-f", targetPath, dbName)
	return cmd.Run()
}

func (a *PostgresAdapter) RestoreDatabase(ctx context.Context, dbName, sourcePath string) error {
	cmd := exec.CommandContext(ctx, "pg_restore", "-d", dbName, "-v", sourcePath)
	return cmd.Run()
}
