package database

import (
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"strings"
)

// MariaDBAdapter manages MariaDB and MySQL databases, users, grants, and dumps.
type MariaDBAdapter struct {
	adminDSN string
}

func NewMariaDBAdapter(adminDSN string) *MariaDBAdapter {
	return &MariaDBAdapter{
		adminDSN: adminDSN,
	}
}

func (a *MariaDBAdapter) Engine() string {
	return "mariadb"
}

func (a *MariaDBAdapter) CreateDatabase(ctx context.Context, dbName, collation string) error {
	// Validate database name identifier to prevent SQL injection
	if strings.ContainsAny(dbName, " ';`\"\\") {
		return fmt.Errorf("invalid characters in database name: %s", dbName)
	}

	db, err := sql.Open("mysql", a.adminDSN)
	if err != nil {
		return fmt.Errorf("failed connecting to database server: %w", err)
	}
	defer db.Close()

	if collation == "" {
		collation = "utf8mb4_unicode_ci"
	}

	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE %s", dbName, collation)
	_, err = db.ExecContext(ctx, query)
	return err
}

func (a *MariaDBAdapter) DropDatabase(ctx context.Context, dbName string) error {
	if strings.ContainsAny(dbName, " ';`\"\\") {
		return fmt.Errorf("invalid characters in database name: %s", dbName)
	}

	db, err := sql.Open("mysql", a.adminDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName)
	_, err = db.ExecContext(ctx, query)
	return err
}

func (a *MariaDBAdapter) CreateUser(ctx context.Context, username, password, host string) error {
	if strings.ContainsAny(username, " ';`\"\\") {
		return fmt.Errorf("invalid characters in username: %s", username)
	}
	if host == "" {
		host = "localhost"
	}

	db, err := sql.Open("mysql", a.adminDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	// Parameterized user creation
	query := fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'%s' IDENTIFIED BY ?", username, host)
	_, err = db.ExecContext(ctx, query, password)
	return err
}

func (a *MariaDBAdapter) GrantPermissions(ctx context.Context, username, dbName, host string, privileges []string) error {
	if strings.ContainsAny(username, " ';`\"\\") || strings.ContainsAny(dbName, " ';`\"\\") {
		return fmt.Errorf("invalid characters in identifiers")
	}
	if host == "" {
		host = "localhost"
	}

	db, err := sql.Open("mysql", a.adminDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	privs := "ALL PRIVILEGES"
	if len(privileges) > 0 {
		privs = strings.Join(privileges, ", ")
	}

	query := fmt.Sprintf("GRANT %s ON `%s`.* TO '%s'@'%s'", privs, dbName, username, host)
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}

	_, _ = db.ExecContext(ctx, "FLUSH PRIVILEGES")
	return nil
}

func (a *MariaDBAdapter) DumpDatabase(ctx context.Context, dbName, targetPath string) error {
	if strings.ContainsAny(dbName, " ';`\"\\") {
		return fmt.Errorf("invalid characters in database name: %s", dbName)
	}

	cmd := exec.CommandContext(ctx, "mysqldump", "--single-transaction", "--quick", dbName, "-r", targetPath)
	return cmd.Run()
}

func (a *MariaDBAdapter) RestoreDatabase(ctx context.Context, dbName, sourcePath string) error {
	if strings.ContainsAny(dbName, " ';`\"\\") {
		return fmt.Errorf("invalid characters in database name: %s", dbName)
	}

	cmd := exec.CommandContext(ctx, "mysql", dbName, "-e", fmt.Sprintf("source %s", sourcePath))
	return cmd.Run()
}
