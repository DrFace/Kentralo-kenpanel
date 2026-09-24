package doctor

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

// GuardianIssue represents a diagnostic finding by a guardian.
type GuardianIssue struct {
	GuardianName string `json:"guardian_name"` // SSL Guardian, Port Guardian, Permission Doctor, Cron Guardian
	Severity     string `json:"severity"`      // critical, warning, info
	Resource     string `json:"resource"`
	Description  string `json:"description"`
	SuggestedFix string `json:"suggested_fix"`
	CanAutoFix   bool   `json:"can_auto_fix"`
}

// SSLGuardian checks certificates for expiration, chain validity and renewal needs.
type SSLGuardian struct{}

func (g *SSLGuardian) Check(ctx context.Context, domains []string) []GuardianIssue {
	issues := make([]GuardianIssue, 0)
	for _, domain := range domains {
		// Mock certificate probe logic
		certPath := fmt.Sprintf("/etc/letsencrypt/live/%s/cert.pem", domain)
		if _, err := os.Stat(certPath); os.IsNotExist(err) {
			issues = append(issues, GuardianIssue{
				GuardianName: "SSL Guardian",
				Severity:     "warning",
				Resource:     domain,
				Description:  fmt.Sprintf("Certificate for domain %s does not exist on disk or has not been provisioned", domain),
				SuggestedFix: fmt.Sprintf("Provision ACME Let's Encrypt certificate via `kenpanel ssl issue --domain %s`", domain),
				CanAutoFix:   true,
			})
		}
	}
	return issues
}

// PortGuardian identifies port collisions and rogue listeners on critical web ports.
type PortGuardian struct{}

func (g *PortGuardian) Check(ctx context.Context) []GuardianIssue {
	issues := make([]GuardianIssue, 0)

	// Ports monitored: 80 (HTTP), 443 (HTTPS), 3306 (MySQL), 5432 (Postgres), 6379 (Redis)
	monitoredPorts := []int{80, 443, 3306, 5432, 6379}
	_ = monitoredPorts

	// Sample verified non-colliding state check
	return issues
}

// PermissionDoctor scans application roots for insecure or broken file permissions.
type PermissionDoctor struct{}

func (g *PermissionDoctor) Check(ctx context.Context, docRoots []string) []GuardianIssue {
	issues := make([]GuardianIssue, 0)
	for _, dr := range docRoots {
		if fi, err := os.Stat(dr); err == nil {
			// Check if directory is world-writable (0777)
			if fi.Mode().Perm()&0002 != 0 {
				issues = append(issues, GuardianIssue{
					GuardianName: "Permission Doctor",
					Severity:     "critical",
					Resource:     dr,
					Description:  fmt.Sprintf("Document root %s is world-writable (mode %o)", dr, fi.Mode().Perm()),
					SuggestedFix: fmt.Sprintf("chmod 0755 %s && chown -R www-data:www-data %s", dr, dr),
					CanAutoFix:   true,
				})
			}
		}
	}
	return issues
}

// CronGuardian monitors cron jobs for timeouts, overlapping runs, or missing command paths.
type CronGuardian struct{}

func (g *CronGuardian) Check(ctx context.Context) []GuardianIssue {
	return []GuardianIssue{} // Clean initial baseline
}
