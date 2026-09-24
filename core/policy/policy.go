package policy

import (
	"strings"
	"sync"
)

// Role represents a standard pre-defined system role.
type Role string

const (
	RoleSuperAdmin         Role = "superadmin"
	RoleOrgAdmin           Role = "org_admin"
	RoleInfrastructureAdmin Role = "infra_admin"
	RoleDeveloper          Role = "developer"
	RoleHostingCustomer    Role = "hosting_customer"
	RoleAuditor            Role = "auditor"
)

// Action defines the operation requested on a resource.
type Action string

const (
	ActionRead     Action = "read"
	ActionWrite    Action = "write"
	ActionExecute  Action = "execute"
	ActionDelete   Action = "delete"
	ActionAdmin    Action = "admin"
)

// PolicyRule defines an authorization grant.
type PolicyRule struct {
	Role     Role   `json:"role"`
	Resource string `json:"resource"` // Glob support, e.g., "servers/*", "websites/*"
	Action   Action `json:"action"`
	Allowed  bool   `json:"allowed"`
}

// Engine evaluates authorization queries against configured rules.
type Engine struct {
	mu    sync.RWMutex
	rules []PolicyRule
}

// NewEngine creates a policy engine pre-configured with secure default roles.
func NewEngine() *Engine {
	e := &Engine{
		rules: make([]PolicyRule, 0),
	}
	e.loadDefaultRules()
	return e
}

func (e *Engine) loadDefaultRules() {
	e.rules = []PolicyRule{
		// SuperAdmin: Full access to everything
		{Role: RoleSuperAdmin, Resource: "*", Action: ActionAdmin, Allowed: true},

		// Auditor: Read-only access to all resources and audit logs
		{Role: RoleAuditor, Resource: "*", Action: ActionRead, Allowed: true},

		// Developer: Read/Write access to apps, websites, databases; cannot modify host firewall or agent
		{Role: RoleDeveloper, Resource: "websites/*", Action: ActionWrite, Allowed: true},
		{Role: RoleDeveloper, Resource: "databases/*", Action: ActionWrite, Allowed: true},
		{Role: RoleDeveloper, Resource: "apps/*", Action: ActionExecute, Allowed: true},

		// Hosting Customer: Scoped to their own assigned tenant resources
		{Role: RoleHostingCustomer, Resource: "websites/*", Action: ActionWrite, Allowed: true},
		{Role: RoleHostingCustomer, Resource: "databases/*", Action: ActionWrite, Allowed: true},
		{Role: RoleHostingCustomer, Resource: "mail/*", Action: ActionWrite, Allowed: true},
	}
}

// IsAuthorized evaluates whether a subject role is permitted to perform action on resource.
func (e *Engine) IsAuthorized(role Role, resource string, action Action) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if role == RoleSuperAdmin {
		return true
	}

	for _, rule := range e.rules {
		if rule.Role == role {
			if matchResource(rule.Resource, resource) {
				if rule.Action == ActionAdmin || rule.Action == action {
					return rule.Allowed
				}
			}
		}
	}

	return false
}

func matchResource(pattern, target string) bool {
	if pattern == "*" || pattern == target {
		return true
	}
	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		return strings.HasPrefix(target, prefix)
	}
	return false
}
