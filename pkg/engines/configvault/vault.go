package configvault

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ConfigRevision represents a snapshot of a configuration file in the vault.
type ConfigRevision struct {
	RevisionID  string    `json:"revision_id"`
	FilePath    string    `json:"file_path"`
	CreatedAt   time.Time `json:"created_at"`
	Actor       string    `json:"actor"`        // e.g. "admin", "changeguard", "upgrade_job"
	JobID       string    `json:"job_id"`
	ContentHash string    `json:"content_hash"`
	Content     string    `json:"content"`      // with secret redaction applied
	Tag         string    `json:"tag"`          // "known-good", "pre-upgrade", "manual"
}

// DiffResult shows changes between two configuration revisions.
type DiffResult struct {
	OriginalRevID string   `json:"original_rev_id"`
	NewRevID      string   `json:"new_rev_id"`
	AddedLines    []string `json:"added_lines"`
	RemovedLines  []string `json:"removed_lines"`
	Identical     bool     `json:"identical"`
}

// ConfigVault manages versioning and rollback of server configuration files.
type ConfigVault struct {
	StorageDir string
	revisions  map[string][]ConfigRevision // key: file path
}

// NewConfigVault creates a new vault.
func NewConfigVault(storageDir string) *ConfigVault {
	if storageDir == "" {
		storageDir = "/var/lib/kenpanel/configvault"
	}
	return &ConfigVault{
		StorageDir: storageDir,
		revisions:  make(map[string][]ConfigRevision),
	}
}

// Snapshot records the current content of a configuration file.
func (v *ConfigVault) Snapshot(filePath, content, actor, jobID, tag string) ConfigRevision {
	cleanContent := v.redactSecrets(content)

	h := sha256.New()
	h.Write([]byte(cleanContent))
	hash := hex.EncodeToString(h.Sum(nil))

	rev := ConfigRevision{
		RevisionID:  fmt.Sprintf("rev-%d", time.Now().UnixNano()),
		FilePath:    filepath.Clean(filePath),
		CreatedAt:   time.Now().UTC(),
		Actor:       actor,
		JobID:       jobID,
		ContentHash: hash,
		Content:     cleanContent,
		Tag:         tag,
	}

	v.revisions[rev.FilePath] = append(v.revisions[rev.FilePath], rev)
	return rev
}

// GetRevisions returns historical snapshots for a file.
func (v *ConfigVault) GetRevisions(filePath string) []ConfigRevision {
	return v.revisions[filepath.Clean(filePath)]
}

// ComputeDiff calculates unified-style diffs between two revisions.
func (v *ConfigVault) ComputeDiff(revA, revB ConfigRevision) *DiffResult {
	res := &DiffResult{
		OriginalRevID: revA.RevisionID,
		NewRevID:      revB.RevisionID,
		AddedLines:    make([]string, 0),
		RemovedLines:  make([]string, 0),
	}

	linesA := strings.Split(revA.Content, "\n")
	linesB := strings.Split(revB.Content, "\n")

	setA := make(map[string]bool)
	for _, l := range linesA {
		setA[strings.TrimSpace(l)] = true
	}

	setB := make(map[string]bool)
	for _, l := range linesB {
		setB[strings.TrimSpace(l)] = true
	}

	for _, l := range linesB {
		trim := strings.TrimSpace(l)
		if trim != "" && !setA[trim] {
			res.AddedLines = append(res.AddedLines, l)
		}
	}

	for _, l := range linesA {
		trim := strings.TrimSpace(l)
		if trim != "" && !setB[trim] {
			res.RemovedLines = append(res.RemovedLines, l)
		}
	}

	res.Identical = len(res.AddedLines) == 0 && len(res.RemovedLines) == 0
	return res
}

func (v *ConfigVault) redactSecrets(content string) string {
	// Redact DB passwords, private keys, API secrets
	lines := strings.Split(content, "\n")
	var out []string
	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "password") || strings.Contains(lower, "secret") || strings.Contains(lower, "private_key") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				out = append(out, fmt.Sprintf("%s = \"<REDACTED_BY_CONFIG_VAULT>\"", parts[0]))
				continue
			}
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}
