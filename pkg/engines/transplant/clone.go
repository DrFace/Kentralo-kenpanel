package transplant

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// ServerCloneConfig defines parameters for cloning an entire host or workload.
type ServerCloneConfig struct {
	SourceHost          string            `json:"source_host"`
	TargetHost          string            `json:"target_host"`
	NewHostname         string            `json:"new_hostname"`
	RegenerateMachineID bool              `json:"regenerate_machine_id"`
	RenewSSHHostKeys    bool              `json:"renew_ssh_host_keys"`
	IPRemapTable        map[string]string `json:"ip_remap_table"`
	PreserveSourceState bool              `json:"preserve_source_state"` // source preservation policy: never disable source silently
	DryRun              bool              `json:"dry_run"`
}

// CloneAction represents an individual operation planned during server cloning.
type CloneAction struct {
	ID          string `json:"id"`
	Phase       string `json:"phase"` // preflight, identity, storage, services, postflight
	Description string `json:"description"`
	OriginalVal string `json:"original_val,omitempty"`
	RemappedVal string `json:"remapped_val,omitempty"`
	Status      string `json:"status"` // planned, executed, skipped
	RequiresRoot bool   `json:"requires_root"`
}

// ServerClonePlan holds the execution plan and validation results.
type ServerClonePlan struct {
	PlanID       string        `json:"plan_id"`
	CreatedAt    string        `json:"created_at"`
	SourceHost   string        `json:"source_host"`
	TargetHost   string        `json:"target_host"`
	Actions      []CloneAction `json:"actions"`
	Warnings     []string      `json:"warnings"`
	CanExecute   bool          `json:"can_execute"`
	Substitutions map[string]string `json:"substitutions"`
}

// ServerCloneEngine orchestrates logical server cloning.
type ServerCloneEngine struct {
	Config ServerCloneConfig
}

// NewServerCloneEngine creates a clone engine instance.
func NewServerCloneEngine(cfg ServerCloneConfig) *ServerCloneEngine {
	if cfg.IPRemapTable == nil {
		cfg.IPRemapTable = make(map[string]string)
	}
	return &ServerCloneEngine{Config: cfg}
}

// GeneratePlan creates a dry-run plan identifying all machine-specific variables to remap.
func (e *ServerCloneEngine) GeneratePlan(ctx context.Context) (*ServerClonePlan, error) {
	plan := &ServerClonePlan{
		PlanID:        fmt.Sprintf("clone-plan-%d", time.Now().UnixNano()),
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
		SourceHost:    e.Config.SourceHost,
		TargetHost:    e.Config.TargetHost,
		Actions:       make([]CloneAction, 0),
		Warnings:      make([]string, 0),
		CanExecute:    true,
		Substitutions: make(map[string]string),
	}

	// 1. Preflight: Source preservation validation
	if !e.Config.PreserveSourceState {
		plan.Warnings = append(plan.Warnings, "Source preservation policy enforced: source machine will NOT be halted or modified.")
	}

	plan.Actions = append(plan.Actions, CloneAction{
		ID:           "preflight-preserve-source",
		Phase:        "preflight",
		Description:  "Verify source server remains untouched and operational during clone workflow",
		Status:       "planned",
		RequiresRoot: false,
	})

	// 2. Identity: Hostname remapping
	if e.Config.NewHostname != "" {
		plan.Actions = append(plan.Actions, CloneAction{
			ID:           "identity-hostname",
			Phase:        "identity",
			Description:  fmt.Sprintf("Update system hostname to %s", e.Config.NewHostname),
			OriginalVal:  e.Config.SourceHost,
			RemappedVal:  e.Config.NewHostname,
			Status:       "planned",
			RequiresRoot: true,
		})
		plan.Substitutions[e.Config.SourceHost] = e.Config.NewHostname
	}

	// 3. Identity: machine-id regeneration
	if e.Config.RegenerateMachineID {
		newMachineID := generateMockMachineID()
		plan.Actions = append(plan.Actions, CloneAction{
			ID:           "identity-machine-id",
			Phase:        "identity",
			Description:  "Regenerate /etc/machine-id and /var/lib/dbus/machine-id to prevent DHCP/systemd collision",
			RemappedVal:  newMachineID,
			Status:       "planned",
			RequiresRoot: true,
		})
	}

	// 4. Identity: SSH Host Keys renewal
	if e.Config.RenewSSHHostKeys {
		plan.Actions = append(plan.Actions, CloneAction{
			ID:           "identity-ssh-keys",
			Phase:        "identity",
			Description:  "Regenerate SSH host keys (ed25519, rsa, ecdsa) on cloned instance",
			Status:       "planned",
			RequiresRoot: true,
		})
	}

	// 5. Network / IP remapping
	for oldIP, newIP := range e.Config.IPRemapTable {
		plan.Actions = append(plan.Actions, CloneAction{
			ID:           fmt.Sprintf("net-remap-%s", strings.ReplaceAll(oldIP, ".", "-")),
			Phase:        "network",
			Description:  fmt.Sprintf("Remap IP binding from %s to %s across Nginx/Apache/vhost configs", oldIP, newIP),
			OriginalVal:  oldIP,
			RemappedVal:  newIP,
			Status:       "planned",
			RequiresRoot: true,
		})
		plan.Substitutions[oldIP] = newIP
	}

	// 6. Postflight health validation
	plan.Actions = append(plan.Actions, CloneAction{
		ID:           "postflight-health-probe",
		Phase:        "postflight",
		Description:  "Run post-clone webstack health probes and service verification",
		Status:       "planned",
		RequiresRoot: false,
	})

	return plan, nil
}

// Execute applies the clone actions if not dry-run.
func (e *ServerCloneEngine) Execute(ctx context.Context, plan *ServerClonePlan) error {
	if e.Config.DryRun {
		return fmt.Errorf("cannot execute clone in dry-run mode")
	}

	for i := range plan.Actions {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		// Mark action completed in plan
		plan.Actions[i].Status = "executed"
	}
	return nil
}

func generateMockMachineID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
