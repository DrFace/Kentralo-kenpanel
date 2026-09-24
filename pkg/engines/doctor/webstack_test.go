package doctor

import (
	"context"
	"testing"
)

func TestWebStackDoctorMissingSocket(t *testing.T) {
	doc := NewWebStackDoctor()
	// Test with a non-existent socket path
	report, err := doc.Diagnose(context.Background(), "localhost", "/tmp/non-existent-php-socket.sock", "")
	if err != nil {
		t.Fatalf("unexpected error running diagnose: %v", err)
	}

	if report.IsHealthy {
		t.Errorf("expected report to flag missing socket as unhealthy")
	}

	foundSocketIssue := false
	for _, f := range report.Findings {
		if f.Code == "PHP_FPM_SOCKET_MISSING" {
			foundSocketIssue = true
			if f.Remediation == nil || f.Remediation.ID != "runtime:restart_php_fpm" {
				t.Errorf("expected remediation action for restarting php-fpm")
			}
		}
	}

	if !foundSocketIssue {
		t.Errorf("expected finding with code PHP_FPM_SOCKET_MISSING")
	}
}
