package rpc

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicConfigWrite(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "rpc-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	server := NewAgentRPCServer([]string{tempDir})

	targetPath := filepath.Join(tempDir, "nginx.conf")
	content := []byte("events {} http { server { listen 80; } }")

	resp, err := server.WriteConfigAtomic(context.Background(), WriteConfigRequest{
		TargetPath:  targetPath,
		Content:     content,
		Permissions: 0644,
	})
	if err != nil {
		t.Fatalf("unexpected error writing config: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected write to succeed")
	}

	readBack, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("failed to read back written file: %v", err)
	}
	if string(readBack) != string(content) {
		t.Errorf("content mismatch")
	}
}

func TestRejectArbitraryShell(t *testing.T) {
	server := NewAgentRPCServer(nil)
	err := server.RejectArbitraryShell("rm -rf /")
	if err == nil {
		t.Errorf("expected error rejecting arbitrary shell execution")
	}
}
