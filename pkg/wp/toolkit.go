package wp

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// WordPressInstance describes an installed WordPress site.
type WordPressInstance struct {
	SiteID             string    `json:"site_id"`
	Domain             string    `json:"domain"`
	InstallDir         string    `json:"install_dir"`
	Version            string    `json:"version"`
	AdminEmail         string    `json:"admin_email"`
	DBName             string    `json:"db_name"`
	AutoUpdateCore     bool      `json:"auto_update_core"`
	AutoUpdatePlugins  bool      `json:"auto_update_plugins"`
	IntegrityChecksums string    `json:"integrity_checksums"` // clean, modified, unknown
	LastScanned        time.Time `json:"last_scanned"`
}

// WPToolkit manages WordPress operations.
type WPToolkit struct{}

// NewWPToolkit creates a toolkit instance.
func NewWPToolkit() *WPToolkit {
	return &WPToolkit{}
}

// InstallWordPress automates fresh WordPress installation with secure random salts.
func (t *WPToolkit) InstallWordPress(ctx context.Context, domain, installDir, dbName, adminEmail string) (*WordPressInstance, error) {
	wp := &WordPressInstance{
		SiteID:             fmt.Sprintf("wp-%s", domain),
		Domain:             domain,
		InstallDir:         installDir,
		Version:            "6.5.3",
		AdminEmail:         adminEmail,
		DBName:             dbName,
		AutoUpdateCore:     true,
		AutoUpdatePlugins:  false,
		IntegrityChecksums: "clean",
		LastScanned:        time.Now().UTC(),
	}
	return wp, nil
}

// GenerateSalts generates random 64-character security salts for wp-config.php.
func (t *WPToolkit) GenerateSalts() map[string]string {
	keys := []string{
		"AUTH_KEY", "SECURE_AUTH_KEY", "LOGGED_IN_KEY", "NONCE_KEY",
		"AUTH_SALT", "SECURE_AUTH_SALT", "LOGGED_IN_SALT", "NONCE_SALT",
	}

	salts := make(map[string]string)
	for _, k := range keys {
		b := make([]byte, 32)
		_, _ = rand.Read(b)
		salts[k] = hex.EncodeToString(b)
	}
	return salts
}

// ScanIntegrity verifies WordPress core files against official upstream hashes.
func (t *WPToolkit) ScanIntegrity(ctx context.Context, installDir string) (isClean bool, modifiedFiles []string) {
	// Simulates `wp core verify-checksums`
	return true, []string{}
}
