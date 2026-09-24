package capsule

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCapsuleRuntimeDetection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "capsule-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create dummy package.json
	os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(`{"name":"test"}`), 0644)

	packager := NewCapsulePackager(tempDir)
	rt, ver, _, _ := packager.DetectRuntime()

	if rt != "nodejs" {
		t.Errorf("expected nodejs, got %s", rt)
	}
	if ver != "20" {
		t.Errorf("expected version 20, got %s", ver)
	}
}

func TestCapsuleSecretRedaction(t *testing.T) {
	packager := NewCapsulePackager("")
	envSample := "APP_NAME=MyApp\nDB_PASSWORD=supersecret123\nAPI_KEY=abcdef123456\nPORT=3000"

	envMap, secretRefs := packager.RedactSecrets(envSample)

	if envMap["DB_PASSWORD"] != "<KENPANEL_SECRET_REF:DB_PASSWORD>" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %s", envMap["DB_PASSWORD"])
	}
	if envMap["API_KEY"] != "<KENPANEL_SECRET_REF:API_KEY>" {
		t.Errorf("expected API_KEY to be redacted, got %s", envMap["API_KEY"])
	}
	if envMap["PORT"] != "3000" {
		t.Errorf("expected PORT to remain 3000, got %s", envMap["PORT"])
	}
	if len(secretRefs) != 2 {
		t.Errorf("expected 2 secret refs, got %d", len(secretRefs))
	}
}
