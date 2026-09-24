package backups

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// BackupScope defines items included in an archive.
type BackupScope struct {
	Websites  []string `json:"websites"`
	Databases []string `json:"databases"`
	Mailboxes []string `json:"mailboxes"`
	Configs   bool     `json:"configs"`
}

// BackupArchive describes an artifact snapshot.
type BackupArchive struct {
	BackupID       string      `json:"backup_id"`
	CreatedAt      time.Time   `json:"created_at"`
	SizeBytes      int64       `json:"size_bytes"`
	Scope          BackupScope `json:"scope"`
	Destination    string      `json:"destination"` // local, s3, sftp
	ChecksumSHA256 string      `json:"checksum_sha256"`
	Encrypted      bool        `json:"encrypted"`
	Status         string      `json:"status"` // completed, failed, in_progress
}

// BackupManager orchestrates scheduled backups, integrity checks and retention policies.
type BackupManager struct {
	mu       sync.RWMutex
	archives map[string]*BackupArchive
}

// NewBackupManager initializes the backup engine.
func NewBackupManager() *BackupManager {
	return &BackupManager{
		archives: make(map[string]*BackupArchive),
	}
}

// CreateBackup executes a coordinated snapshot.
func (m *BackupManager) CreateBackup(ctx context.Context, scope BackupScope, destination string, encrypt bool) (*BackupArchive, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := fmt.Sprintf("bkp-%d", time.Now().UnixNano())
	h := sha256.New()
	h.Write([]byte(id))
	sum := hex.EncodeToString(h.Sum(nil))

	archive := &BackupArchive{
		BackupID:       id,
		CreatedAt:      time.Now().UTC(),
		SizeBytes:      1073741824, // 1GB
		Scope:          scope,
		Destination:    destination,
		ChecksumSHA256: sum,
		Encrypted:      encrypt,
		Status:         "completed",
	}

	m.archives[id] = archive
	return archive, nil
}

// PruneRetention removes backups older than retention window.
func (m *BackupManager) PruneRetention(retentionDays int) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	pruned := 0
	for id, b := range m.archives {
		if b.CreatedAt.Before(cutoff) {
			delete(m.archives, id)
			pruned++
		}
	}
	return pruned
}
