package mail

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// MailboxAccount describes an active email address.
type MailboxAccount struct {
	Email       string    `json:"email"`
	Domain      string    `json:"domain"`
	QuotaMB     int       `json:"quota_mb"`
	UsedMB      int       `json:"used_mb"`
	IsSuspended bool      `json:"is_suspended"`
	CreatedAt   time.Time `json:"created_at"`
}

// MailAlias routes emails from an alias address to destinations.
type MailAlias struct {
	SourceEmail  string   `json:"source_email"`
	Destinations []string `json:"destinations"`
}

// MailManager orchestrates mail services (Postfix, Dovecot, Rspamd).
type MailManager struct {
	mu        sync.RWMutex
	mailboxes map[string]*MailboxAccount
	aliases   map[string]*MailAlias
}

// NewMailManager creates a mail stack manager.
func NewMailManager() *MailManager {
	return &MailManager{
		mailboxes: make(map[string]*MailboxAccount),
		aliases:   make(map[string]*MailAlias),
	}
}

// CreateMailbox provisions a new mailbox.
func (m *MailManager) CreateMailbox(ctx context.Context, email, rawPassword string, quotaMB int) (*MailboxAccount, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	email = strings.ToLower(strings.TrimSpace(email))
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid email format: %s", email)
	}

	if _, exists := m.mailboxes[email]; exists {
		return nil, fmt.Errorf("mailbox %s already exists", email)
	}

	acc := &MailboxAccount{
		Email:       email,
		Domain:      parts[1],
		QuotaMB:     quotaMB,
		UsedMB:      0,
		IsSuspended: false,
		CreatedAt:   time.Now().UTC(),
	}

	m.mailboxes[email] = acc
	return acc, nil
}

// GenerateAuthenticationRecords creates verified DNS instructions for SPF, DKIM, and DMARC.
func (m *MailManager) GenerateAuthenticationRecords(domain, mailHost, dkimPublicKey string) map[string]string {
	return map[string]string{
		"SPF":   fmt.Sprintf("v=spf1 mx a:%s ~all", mailHost),
		"DKIM":  fmt.Sprintf("v=DKIM1; k=rsa; p=%s", dkimPublicKey),
		"DMARC": fmt.Sprintf("v=DMARC1; p=quarantine; rua=mailto:dmarc@%s; pct=100", domain),
		"MX":    fmt.Sprintf("10 %s.", mailHost),
	}
}
