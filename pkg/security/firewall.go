package security

import (
	"fmt"
	"sync"
	"time"
)

// FirewallRule describes a packet filtering directive.
type FirewallRule struct {
	ID        string `json:"id"`
	Port      int    `json:"port"`
	Protocol  string `json:"protocol"` // tcp, udp
	Action    string `json:"action"`   // allow, deny, reject
	SourceIP  string `json:"source_ip"` // "any" or CIDR e.g. "198.51.100.0/24"
	Comment   string `json:"comment"`
	Enabled   bool   `json:"enabled"`
}

// BannedIP tracks an address blocked for suspicious activity.
type BannedIP struct {
	IP        string    `json:"ip"`
	Reason    string    `json:"reason"`
	Jail      string    `json:"jail"` // sshd, kenpanel-auth, nginx-botsearch
	BannedAt  time.Time `json:"banned_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SecurityManager oversees firewall rules, brute-force defenses and WAF policies.
type SecurityManager struct {
	mu      sync.RWMutex
	rules   map[string]*FirewallRule
	bans    map[string]*BannedIP
	wafMode string // off, detection, blocking
}

// NewSecurityManager creates security controller.
func NewSecurityManager() *SecurityManager {
	mgr := &SecurityManager{
		rules:   make(map[string]*FirewallRule),
		bans:    make(map[string]*BannedIP),
		wafMode: "blocking",
	}

	// Default baseline
	defaultRules := []*FirewallRule{
		{ID: "fw-ssh", Port: 22, Protocol: "tcp", Action: "allow", SourceIP: "any", Comment: "SSH remote management", Enabled: true},
		{ID: "fw-http", Port: 80, Protocol: "tcp", Action: "allow", SourceIP: "any", Comment: "Web HTTP traffic", Enabled: true},
		{ID: "fw-https", Port: 443, Protocol: "tcp", Action: "allow", SourceIP: "any", Comment: "Web HTTPS traffic", Enabled: true},
		{ID: "fw-panel", Port: 2083, Protocol: "tcp", Action: "allow", SourceIP: "any", Comment: "KenPanel control plane", Enabled: true},
	}
	for _, r := range defaultRules {
		mgr.rules[r.ID] = r
	}

	return mgr
}

// AddRule adds an incoming firewall rule.
func (s *SecurityManager) AddRule(rule *FirewallRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if rule.Port <= 0 || rule.Port > 65535 {
		return fmt.Errorf("invalid port %d", rule.Port)
	}

	s.rules[rule.ID] = rule
	return nil
}

// BanIP blocks an IP across system firewalls and intrusion prevention.
func (s *SecurityManager) BanIP(ip, reason, jail string, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.bans[ip] = &BannedIP{
		IP:        ip,
		Reason:    reason,
		Jail:      jail,
		BannedAt:  time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(duration),
	}
}

// ListBans returns currently active IP blocks.
func (s *SecurityManager) ListBans() []*BannedIP {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*BannedIP, 0, len(s.bans))
	now := time.Now().UTC()
	for _, b := range s.bans {
		if now.Before(b.ExpiresAt) {
			list = append(list, b)
		}
	}
	return list
}
