package website

import (
	"context"
	"testing"
)

func TestWebsiteLifecycle(t *testing.T) {
	engine := NewLifecycleEngine()

	site, err := engine.CreateSite(context.Background(), "mytestsite.com", SiteTypePHP, "8.2")
	if err != nil {
		t.Fatalf("unexpected error creating site: %v", err)
	}

	if site.PrimaryDomain != "mytestsite.com" {
		t.Errorf("domain mismatch")
	}
	if site.Status != "active" {
		t.Errorf("expected active status")
	}

	// Test maintenance mode toggle
	if err := engine.SetMaintenanceMode(site.ID, true, "<h1>Under Maintenance</h1>"); err != nil {
		t.Fatalf("unexpected error enabling maintenance mode: %v", err)
	}
	if !site.MaintenanceMode {
		t.Errorf("expected maintenance mode to be enabled")
	}

	// Test cloning
	cloned, err := engine.CloneSite(site.ID, "staging.mytestsite.com")
	if err != nil {
		t.Fatalf("unexpected error cloning site: %v", err)
	}
	if cloned.PrimaryDomain != "staging.mytestsite.com" {
		t.Errorf("cloned domain mismatch")
	}
}
