package doctor

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"
)

// IssueSeverity defines diagnostic alert levels.
type IssueSeverity string

const (
	SeverityInfo     IssueSeverity = "info"
	SeverityWarning  IssueSeverity = "warning"
	SeverityCritical IssueSeverity = "critical"
)

// DiagnosticFinding represents an evidence-backed issue detected in the stack.
type DiagnosticFinding struct {
	Component   string        `json:"component"` // "dns", "proxy", "runtime", "database", "storage"
	Code        string        `json:"code"`
	Severity    IssueSeverity `json:"severity"`
	Message     string        `json:"message"`
	Evidence    string        `json:"evidence"`
	Remediation *FixAction    `json:"remediation,omitempty"`
}

// FixAction represents a typed, safe remediation operation.
type FixAction struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Privilege      string `json:"privilege"`
	ServiceRestart bool   `json:"service_restart"`
	DowntimeSec    int    `json:"downtime_sec"`
}

// DiagnosticReport is the complete evaluation of a website/host dependency tree.
type DiagnosticReport struct {
	Domain      string              `json:"domain"`
	EvaluatedAt time.Time           `json:"evaluated_at"`
	IsHealthy   bool                `json:"is_healthy"`
	Findings    []DiagnosticFinding `json:"findings"`
}

// WebStackDoctor runs vertical dependency diagnostics across the web hosting stack.
type WebStackDoctor struct{}

func NewWebStackDoctor() *WebStackDoctor {
	return &WebStackDoctor{}
}

// Diagnose evaluates the vertical health graph for a given domain and socket target.
func (d *WebStackDoctor) Diagnose(ctx context.Context, domain string, phpSocket string, dbDSN string) (*DiagnosticReport, error) {
	report := &DiagnosticReport{
		Domain:      domain,
		EvaluatedAt: time.Now().UTC(),
		IsHealthy:   true,
		Findings:    make([]DiagnosticFinding, 0),
	}

	// 1. DNS Resolution Check
	ips, err := net.LookupHost(domain)
	if err != nil || len(ips) == 0 {
		report.IsHealthy = false
		report.Findings = append(report.Findings, DiagnosticFinding{
			Component: "dns",
			Code:      "DNS_UNRESOLVED",
			Severity:  SeverityCritical,
			Message:   fmt.Sprintf("Domain %s does not resolve to any public IP address", domain),
			Evidence:  fmt.Sprintf("LookupHost error: %v", err),
			Remediation: &FixAction{
				ID:             "dns:configure_records",
				Title:          "Verify and point DNS A/AAAA records to this server",
				Privilege:      "user",
				ServiceRestart: false,
				DowntimeSec:    0,
			},
		})
	}

	// 2. PHP-FPM Socket Check
	if phpSocket != "" {
		if _, err := os.Stat(phpSocket); os.IsNotExist(err) {
			report.IsHealthy = false
			report.Findings = append(report.Findings, DiagnosticFinding{
				Component: "runtime",
				Code:      "PHP_FPM_SOCKET_MISSING",
				Severity:  SeverityCritical,
				Message:   fmt.Sprintf("PHP-FPM socket %s does not exist on disk", phpSocket),
				Evidence:  "Socket file stat returned not found",
				Remediation: &FixAction{
					ID:             "runtime:restart_php_fpm",
					Title:          "Restart PHP-FPM service to re-create Unix domain socket",
					Privilege:      "root",
					ServiceRestart: true,
					DowntimeSec:    1,
				},
			})
		}
	}

	return report, nil
}
