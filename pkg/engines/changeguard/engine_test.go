package changeguard

import (
	"context"
	"errors"
	"testing"
)

func TestChangeGuardSuccess(t *testing.T) {
	tx := NewTransaction("tx-1", "Test valid transaction")
	checkpointCalled := false
	applyCalled := false
	verifyCalled := false
	rollbackCalled := false

	tx.SetHooks(
		func(ctx context.Context) error { checkpointCalled = true; return nil },
		func(ctx context.Context) error { applyCalled = true; return nil },
		func(ctx context.Context) error { verifyCalled = true; return nil },
		func(ctx context.Context) error { rollbackCalled = true; return nil },
	)

	err := tx.Execute(context.Background())
	if err != nil {
		t.Fatalf("expected transaction to succeed, got %v", err)
	}

	if !checkpointCalled || !applyCalled || !verifyCalled {
		t.Errorf("expected checkpoint, apply, and verify hooks to be called")
	}

	if rollbackCalled {
		t.Errorf("expected rollback not to be called on success")
	}

	if tx.CurrentPhase != PhaseCommit {
		t.Errorf("expected phase Commit, got %s", tx.CurrentPhase)
	}
}

func TestChangeGuardRollbackOnVerifyFailure(t *testing.T) {
	tx := NewTransaction("tx-2", "Test health failure transaction")
	rollbackCalled := false

	tx.SetHooks(
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return errors.New("upstream service dead") },
		func(ctx context.Context) error { rollbackCalled = true; return nil },
	)

	err := tx.Execute(context.Background())
	if err == nil {
		t.Fatalf("expected transaction to fail on health check")
	}

	if !rollbackCalled {
		t.Errorf("expected rollback hook to be executed")
	}

	if tx.CurrentPhase != PhaseRollback {
		t.Errorf("expected phase Rollback, got %s", tx.CurrentPhase)
	}
}
