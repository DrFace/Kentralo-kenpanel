package security

import (
	"testing"
)

func TestSecretVaultEncryption(t *testing.T) {
	vault, err := NewSecretVault("")
	if err != nil {
		t.Fatalf("failed to create secret vault: %v", err)
	}

	secretKey := "db_master_password"
	secretVal := "Super$ecurePassword#2026"

	if err := vault.Put(secretKey, secretVal); err != nil {
		t.Fatalf("failed to put secret: %v", err)
	}

	decrypted, err := vault.Get(secretKey)
	if err != nil {
		t.Fatalf("failed to get secret: %v", err)
	}

	if decrypted != secretVal {
		t.Errorf("decrypted secret mismatch: got %s, want %s", decrypted, secretVal)
	}
}
