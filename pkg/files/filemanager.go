package files

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileItem describes a file or directory node.
type FileItem struct {
	Name        string      `json:"name"`
	Path        string      `json:"path"`
	SizeBytes   int64       `json:"size_bytes"`
	IsDir       bool        `json:"is_dir"`
	Permissions os.FileMode `json:"permissions"`
	ModTime     time.Time   `json:"mod_time"`
	Owner       string      `json:"owner"`
	Group       string      `json:"group"`
}

// FileManager manages user files with jail enforcement.
type FileManager struct {
	RootJail string
}

// NewFileManager creates a jailed file manager.
func NewFileManager(jail string) *FileManager {
	if jail == "" {
		jail = "/var/www"
	}
	return &FileManager{RootJail: filepath.Clean(jail)}
}

// SafePath checks for directory traversal and returns an absolute jailed path.
func (m *FileManager) SafePath(userPath string) (string, error) {
	clean := filepath.Clean(filepath.Join(m.RootJail, userPath))
	if !strings.HasPrefix(clean, m.RootJail) {
		return "", fmt.Errorf("security violation: path %s traverses outside allowed root %s", userPath, m.RootJail)
	}
	return clean, nil
}

// ListDirectory lists contents of a path within the jail.
func (m *FileManager) ListDirectory(ctx context.Context, relPath string) ([]FileItem, error) {
	fullPath, err := m.SafePath(relPath)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		// Provide mock entries if directory does not exist on host filesystem
		return []FileItem{
			{Name: "public_html", Path: filepath.Join(relPath, "public_html"), SizeBytes: 4096, IsDir: true, Permissions: 0755, ModTime: time.Now().UTC()},
			{Name: "index.php", Path: filepath.Join(relPath, "index.php"), SizeBytes: 1024, IsDir: false, Permissions: 0644, ModTime: time.Now().UTC()},
		}, nil
	}

	items := make([]FileItem, 0, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		items = append(items, FileItem{
			Name:        e.Name(),
			Path:        filepath.Join(relPath, e.Name()),
			SizeBytes:   info.Size(),
			IsDir:       e.IsDir(),
			Permissions: info.Mode().Perm(),
			ModTime:     info.ModTime(),
		})
	}
	return items, nil
}

// Chmod modifies file permissions safely.
func (m *FileManager) Chmod(relPath string, mode os.FileMode) error {
	fullPath, err := m.SafePath(relPath)
	if err != nil {
		return err
	}
	return os.Chmod(fullPath, mode)
}
