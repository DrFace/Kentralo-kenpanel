package mail

import (
	"context"
	"strings"
	"testing"
)

func TestMailboxCreationAndAuthRecords(t *testing.T) {
	mgr := NewMailManager()

	acc, err := mgr.CreateMailbox(context.Background(), "info@company.com", "SecretPass123!", 2048)
	if err != nil {
		t.Fatalf("unexpected error creating mailbox: %v", err)
	}
	if acc.Domain != "company.com" {
		t.Errorf("domain mismatch")
	}

	records := mgr.GenerateAuthenticationRecords("company.com", "mail.company.com", "MIGfMA0GCS...")
	if !strings.Contains(records["SPF"], "v=spf1") {
		t.Errorf("expected SPF record")
	}
	if !strings.Contains(records["DKIM"], "v=DKIM1") {
		t.Errorf("expected DKIM record")
	}
	if !strings.Contains(records["DMARC"], "v=DMARC1") {
		t.Errorf("expected DMARC record")
	}
}
