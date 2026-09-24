package registry

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"
)

// Status represents the development/release status of a requirement.
type Status string

const (
	StatusPlanned      Status = "Planned"
	StatusDesigned     Status = "Designed"
	StatusExperimental Status = "Experimental"
	StatusAlpha        Status = "Alpha"
	StatusBeta         Status = "Beta"
	StatusStable       Status = "Stable"
	StatusDeprecated   Status = "Deprecated"
	StatusUnsupported  Status = "Unsupported"
)

// RiskClass defines the operational risk category of an action or subsystem.
type RiskClass string

const (
	RiskReadOnly    RiskClass = "read-only"
	RiskLow         RiskClass = "low"
	RiskModerate    RiskClass = "moderate"
	RiskHigh        RiskClass = "high"
	RiskDestructive RiskClass = "destructive"
)

// RollbackClass defines how an operation can be reversed.
type RollbackClass string

const (
	RollbackAutomatic         RollbackClass = "automatic"
	RollbackCheckpointRestore RollbackClass = "checkpoint_restore"
	RollbackManual            RollbackClass = "manual"
	RollbackNotReversible     RollbackClass = "not_reversible"
)

// ExecutionMode defines how the requirement can run.
type ExecutionMode string

const (
	ModeIntegrated       ExecutionMode = "integrated"
	ModeStandaloneCLI    ExecutionMode = "standalone_cli"
	ModeStandaloneDaemon ExecutionMode = "standalone_daemon"
	ModeBootableRecovery ExecutionMode = "bootable_recovery"
	ModeMultiple         ExecutionMode = "multiple"
)

// Requirement represents an authoritative normative specification item.
type Requirement struct {
	ID                          string          `json:"id"`
	Title                       string          `json:"title"`
	Module                      string          `json:"module"`
	Subsystem                   string          `json:"subsystem"`
	Description                 string          `json:"description"`
	Status                      Status          `json:"status"`
	TargetMilestone             string          `json:"target_milestone"`
	PrivilegeLevel              string          `json:"privilege_level"`
	RiskClass                   RiskClass       `json:"risk_class"`
	ServiceRestartRequired      bool            `json:"service_restart_required"`
	RebootRequired              bool            `json:"reboot_required"`
	PotentialDowntimeSeconds    int             `json:"potential_downtime_seconds"`
	RollbackClass               RollbackClass   `json:"rollback_class"`
	OfflineCapable              bool            `json:"offline_capable"`
	ExternalNetworkDependencies []string        `json:"external_network_dependencies"`
	ExecutionModes              []ExecutionMode `json:"execution_modes"`
	TargetOS                    []string        `json:"target_os"`
	Architectures               []string        `json:"architectures"`
	CodePaths                   []string        `json:"code_paths"`
	APIEndpoints                []string        `json:"api_endpoints"`
	CLICommands                 []string        `json:"cli_commands"`
	UIRoutes                    []string        `json:"ui_routes"`
	Tests                       []string        `json:"tests"`
	GuideRoute                  string          `json:"guide_route"`
	KnownLimitations            []string        `json:"known_limitations"`
	FirstShipped                *string         `json:"first_shipped"`
	DeprecatedBy                *string         `json:"deprecated_by"`
}

// Registry stores and provides indexed queries over all product requirements.
type Registry struct {
	mu           sync.RWMutex
	requirements map[string]Requirement
}

// NewRegistry creates an empty registry instance.
func NewRegistry() *Registry {
	return &Registry{
		requirements: make(map[string]Requirement),
	}
}

// LoadFromJSON parses raw JSON requirement records into the registry.
func (r *Registry) LoadFromJSON(data []byte) error {
	var list []Requirement
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("failed to unmarshal requirement JSON: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	for _, req := range list {
		r.requirements[req.ID] = req
	}
	return nil
}

// Get retrieves a requirement by its unique ID.
func (r *Registry) Get(id string) (Requirement, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	req, ok := r.requirements[id]
	return req, ok
}

// All returns a slice of all registered requirements.
func (r *Registry) All() []Requirement {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]Requirement, 0, len(r.requirements))
	for _, req := range r.requirements {
		list = append(list, req)
	}
	return list
}

// FilterBySubsystem returns all requirements belonging to a named subsystem.
func (r *Registry) FilterBySubsystem(subsystem string) []Requirement {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []Requirement
	for _, req := range r.requirements {
		if req.Subsystem == subsystem {
			list = append(list, req)
		}
	}
	return list
}
