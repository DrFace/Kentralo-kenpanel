package changeguard

import (
	"context"
	"fmt"
	"time"
)

// TransactionPhase tracks the current stage in a ChangeGuard transaction.
type TransactionPhase string

const (
	PhaseDiscover   TransactionPhase = "discover"
	PhaseValidate   TransactionPhase = "validate"
	PhasePreview    TransactionPhase = "preview"
	PhaseCheckpoint TransactionPhase = "checkpoint"
	PhaseTest       TransactionPhase = "test"
	PhaseApply      TransactionPhase = "apply"
	PhaseVerify     TransactionPhase = "verify"
	PhaseCommit     TransactionPhase = "commit"
	PhaseRollback   TransactionPhase = "rollback"
)

// StepHook represents an individual hook executed during a transaction phase.
type StepHook func(ctx context.Context) error

// Transaction models an atomic reversible operation.
type Transaction struct {
	ID          string           `json:"id"`
	Description string           `json:"description"`
	CurrentPhase TransactionPhase `json:"current_phase"`
	StartTime   time.Time        `json:"start_time"`
	EndTime     *time.Time       `json:"end_time,omitempty"`

	onCheckpoint StepHook
	onApply      StepHook
	onVerify     StepHook
	onRollback   StepHook
}

// NewTransaction initializes a new ChangeGuard transaction.
func NewTransaction(id, description string) *Transaction {
	return &Transaction{
		ID:           id,
		Description:  description,
		CurrentPhase: PhaseDiscover,
		StartTime:    time.Now().UTC(),
	}
}

// SetHooks registers the phase hooks for this transaction.
func (t *Transaction) SetHooks(checkpoint, apply, verify, rollback StepHook) {
	t.onCheckpoint = checkpoint
	t.onApply = apply
	t.onVerify = verify
	t.onRollback = rollback
}

// Execute runs the full transaction state machine with automatic rollback on verification failure.
func (t *Transaction) Execute(ctx context.Context) error {
	// 1. Checkpoint
	t.CurrentPhase = PhaseCheckpoint
	if t.onCheckpoint != nil {
		if err := t.onCheckpoint(ctx); err != nil {
			return fmt.Errorf("transaction aborted at checkpoint: %w", err)
		}
	}

	// 2. Apply
	t.CurrentPhase = PhaseApply
	if t.onApply != nil {
		if err := t.onApply(ctx); err != nil {
			t.rollback(ctx)
			return fmt.Errorf("transaction failed during apply (rolled back): %w", err)
		}
	}

	// 3. Verify Health
	t.CurrentPhase = PhaseVerify
	if t.onVerify != nil {
		if err := t.onVerify(ctx); err != nil {
			t.rollback(ctx)
			return fmt.Errorf("health verification failed post-apply (rolled back): %w", err)
		}
	}

	// 4. Commit
	t.CurrentPhase = PhaseCommit
	now := time.Now().UTC()
	t.EndTime = &now
	return nil
}

func (t *Transaction) rollback(ctx context.Context) {
	t.CurrentPhase = PhaseRollback
	if t.onRollback != nil {
		_ = t.onRollback(ctx)
	}
	now := time.Now().UTC()
	t.EndTime = &now
}
