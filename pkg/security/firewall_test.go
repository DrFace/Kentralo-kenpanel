package security

import (
	"testing"
	"time"
)

func TestFirewallRuleAndBan(t *testing.T) {
	mgr := NewSecurityManager()

	err := mgr.AddRule(&FirewallRule{
		ID:       "fw-custom",
		Port:     8080,
		Protocol: "tcp",
		Action:   "allow",
		SourceIP: "192.168.1.0/24",
		Enabled:  true,
	})
	if err != nil {
		t.Fatalf("unexpected error adding rule: %v", err)
	}

	mgr.BanIP("198.51.100.99", "Failed SSH logins", "sshd", 1*time.Hour)
	bans := mgr.ListBans()
	if len(bans) != 1 {
		t.Errorf("expected 1 ban, got %d", len(bans))
	}
	if bans[0].IP != "198.51.100.99" {
		t.Errorf("banned IP mismatch")
	}
}
