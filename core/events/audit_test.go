package events

import (
	"testing"
)

func TestAuditLedgerIntegrity(t *testing.T) {
	secret := []byte("kenpanel-test-secret-key-32bytes!")
	ledger := NewAuditLedger(secret)

	// Record sequence of events
	ledger.RecordEvent("evt-1", "user-admin", "192.168.1.10", "create_vhost", "domain:example.com", "", "hash_state_1")
	ledger.RecordEvent("evt-2", "user-admin", "192.168.1.10", "issue_ssl", "cert:example.com", "hash_state_1", "hash_state_2")
	ledger.RecordEvent("evt-3", "user-dev", "10.0.0.5", "deploy_app", "app:php-8.3", "hash_state_2", "hash_state_3")

	valid, err := ledger.VerifyIntegrity()
	if err != nil || !valid {
		t.Fatalf("expected audit chain to be valid, got err: %v", err)
	}

	// Tamper with event 2
	ledger.events[1].Action = "tampered_action"
	validTampered, _ := ledger.VerifyIntegrity()
	if validTampered {
		t.Fatalf("expected tampered ledger to fail integrity verification")
	}
}
