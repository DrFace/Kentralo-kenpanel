package mailbridge

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"
)

// MailMigrationOptions holds configuration for a mail migration run.
type MailMigrationOptions struct {
	SourceServer       string   `json:"source_server"`
	SourcePort         int      `json:"source_port"` // e.g. 993 IMAP SSL
	TargetServer       string   `json:"target_server"`
	TargetPort         int      `json:"target_port"`
	Domains            []string `json:"domains"`
	PreserveHierarchy  bool     `json:"preserve_hierarchy"`
	PreserveFlags      bool     `json:"preserve_flags"`      // \Seen, \Flagged, etc.
	PreserveInternalDate bool   `json:"preserve_internal_date"`
	RegenerateDKIM     bool     `json:"regenerate_dkim"`
	IncrementalMode    bool     `json:"incremental_mode"`
	ExcludeFolders     []string `json:"exclude_folders"` // e.g. Trash, Spam
}

// MailboxAccount represents a single account being migrated.
type MailboxAccount struct {
	Email         string            `json:"email"`
	Domain        string            `json:"domain"`
	QuotaBytes    int64             `json:"quota_bytes"`
	Aliases       []string          `json:"aliases"`
	Forwarders    []string          `json:"forwarders"`
	CatchAll      bool              `json:"catch_all"`
	Folders       []MailFolder      `json:"folders"`
	CustomFilters map[string]string `json:"custom_filters,omitempty"`
}

// MailFolder represents an IMAP folder/Maildir directory.
type MailFolder struct {
	Name         string `json:"name"` // e.g. INBOX, Sent, Archive
	MessageCount int    `json:"message_count"`
	SizeBytes    int64  `json:"size_bytes"`
	UIDValidity  uint32 `json:"uid_validity"`
}

// MessageTransferRecord tracks individual message sync without logging message body content.
type MessageTransferRecord struct {
	MessageID   string    `json:"message_id"` // Message-ID header or SHA256 of headers
	Folder      string    `json:"folder"`
	SizeBytes   int64     `json:"size_bytes"`
	Flags       []string  `json:"flags"`
	SyncStatus  string    `json:"sync_status"` // synced, skipped, failed
	Transferred time.Time `json:"transferred"`
}

// DNSCutoverPlan specifies required DNS changes for mail delivery.
type DNSCutoverPlan struct {
	Domain      string   `json:"domain"`
	MXRecords   []string `json:"mx_records"`
	SPFRecord   string   `json:"spf_record"`
	DKIMSelector string  `json:"dkim_selector"`
	DKIMRecord  string   `json:"dkim_record"`
	DMARCRecord string   `json:"dmarc_record"`
	TTLSeconds  int      `json:"ttl_seconds"`
}

// MailBridgeEngine orchestrates mailbox streaming and migration.
type MailBridgeEngine struct {
	Options MailMigrationOptions
	mu      sync.RWMutex
	synced  map[string]bool // Deduplication map: sha256(mailbox+folder+msgid)
}

// NewMailBridgeEngine creates a new MailBridge instance.
func NewMailBridgeEngine(opts MailMigrationOptions) *MailBridgeEngine {
	return &MailBridgeEngine{
		Options: opts,
		synced:  make(map[string]bool),
	}
}

// GenerateDNSCutover creates DNS records for SPF, DKIM, DMARC and MX.
func (e *MailBridgeEngine) GenerateDNSCutover(domain string, targetMailServer string, dkimPubKey string) *DNSCutoverPlan {
	if dkimPubKey == "" {
		dkimPubKey = "p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC0...mock...pubkey"
	}
	return &DNSCutoverPlan{
		Domain: domain,
		MXRecords: []string{
			fmt.Sprintf("10 mail.%s.", domain),
		},
		SPFRecord:   fmt.Sprintf("v=spf1 mx a:%s ~all", targetMailServer),
		DKIMSelector: "kenpanel",
		DKIMRecord:  fmt.Sprintf("v=DKIM1; k=rsa; %s", dkimPubKey),
		DMARCRecord: fmt.Sprintf("v=DMARC1; p=quarantine; rua=mailto:dmarc@%s; pct=100", domain),
		TTLSeconds:  300, // Short TTL recommendation prior to migration
	}
}

// SyncMailbox performs streaming sync of an account's folders and messages.
// CRITICAL PRIVACY REQUIREMENT (Spec Section 49.3): Never store message content in logs.
func (e *MailBridgeEngine) SyncMailbox(ctx context.Context, account *MailboxAccount, progressFn func(folder string, synced, total int)) ([]MessageTransferRecord, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	records := make([]MessageTransferRecord, 0)

	for _, folder := range account.Folders {
		select {
		case <-ctx.Done():
			return records, ctx.Err()
		default:
		}

		// Check if folder is excluded
		if e.isFolderExcluded(folder.Name) {
			continue
		}

		// Stream messages in folder
		for i := 0; i < folder.MessageCount; i++ {
			msgUID := fmt.Sprintf("%s-%s-%d", account.Email, folder.Name, i+1)
			hashKey := e.hashMessageKey(account.Email, folder.Name, msgUID)

			if e.Options.IncrementalMode && e.synced[hashKey] {
				// Already transferred in previous sync run; skip to avoid duplicates
				records = append(records, MessageTransferRecord{
					MessageID:   msgUID,
					Folder:      folder.Name,
					SizeBytes:   2048,
					Flags:       []string{"\\Seen"},
					SyncStatus:  "skipped",
					Transferred: time.Now(),
				})
				continue
			}

			// Simulate transfer of metadata and raw RFC822 stream directly without saving body to disk
			e.synced[hashKey] = true
			records = append(records, MessageTransferRecord{
				MessageID:   msgUID,
				Folder:      folder.Name,
				SizeBytes:   4096,
				Flags:       []string{"\\Seen"},
				SyncStatus:  "synced",
				Transferred: time.Now(),
			})

			if progressFn != nil && (i%20 == 0 || i == folder.MessageCount-1) {
				progressFn(folder.Name, i+1, folder.MessageCount)
			}
		}
	}

	return records, nil
}

func (e *MailBridgeEngine) isFolderExcluded(folder string) bool {
	for _, ex := range e.Options.ExcludeFolders {
		if strings.EqualFold(ex, folder) {
			return true
		}
	}
	return false
}

func (e *MailBridgeEngine) hashMessageKey(email, folder, msgID string) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%s", email, folder, msgID)))
	return hex.EncodeToString(h.Sum(nil))
}
