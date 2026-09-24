package events

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// AuditEvent represents an immutable, tamper-evident log record of a system change.
type AuditEvent struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	ActorID   string    `json:"actor_id"`
	ActorIP   string    `json:"actor_ip"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	OldState  string    `json:"old_state_hash,omitempty"`
	NewState  string    `json:"new_state_hash"`
	PrevHash  string    `json:"prev_hash"`
	Hash      string    `json:"hash"`
}

// AuditLedger manages an append-only cryptographic chain of audit events.
type AuditLedger struct {
	mu       sync.RWMutex
	secret   []byte
	lastHash string
	events   []AuditEvent
}

// NewAuditLedger initializes an audit ledger with an HMAC secret key.
func NewAuditLedger(secret []byte) *AuditLedger {
	return &AuditLedger{
		secret:   secret,
		lastHash: "0000000000000000000000000000000000000000000000000000000000000000",
		events:   make([]AuditEvent, 0),
	}
}

// RecordEvent appends a new verified audit event to the ledger.
func (l *AuditLedger) RecordEvent(id, actorID, actorIP, action, resource, oldHash, newHash string) AuditEvent {
	l.mu.Lock()
	defer l.mu.Unlock()

	event := AuditEvent{
		ID:        id,
		Timestamp: time.Now().UTC(),
		ActorID:   actorID,
		ActorIP:   actorIP,
		Action:    action,
		Resource:  resource,
		OldState:  oldHash,
		NewState:  newHash,
		PrevHash:  l.lastHash,
	}

	event.Hash = l.computeHash(event)
	l.lastHash = event.Hash
	l.events = append(l.events, event)

	return event
}

// VerifyIntegrity checks every event in the ledger against its cryptographic chain link.
func (l *AuditLedger) VerifyIntegrity() (bool, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	expectedPrev := "0000000000000000000000000000000000000000000000000000000000000000"
	for idx, event := range l.events {
		if event.PrevHash != expectedPrev {
			return false, fmt.Errorf("chain broken at event index %d: expected prev %s, got %s", idx, expectedPrev, event.PrevHash)
		}

		computed := l.computeHash(event)
		if computed != event.Hash {
			return false, fmt.Errorf("hash mismatch at event index %d: computed %s, stored %s", idx, computed, event.Hash)
		}

		expectedPrev = event.Hash
	}

	return true, nil
}

func (l *AuditLedger) computeHash(e AuditEvent) string {
	payload := fmt.Sprintf("%s|%d|%s|%s|%s|%s|%s|%s|%s",
		e.ID, e.Timestamp.UnixNano(), e.ActorID, e.ActorIP, e.Action, e.Resource, e.OldState, e.NewState, e.PrevHash)

	h := hmac.New(sha256.New, l.secret)
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}
