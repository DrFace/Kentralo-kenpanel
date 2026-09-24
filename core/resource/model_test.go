package resource

import (
	"testing"
)

func TestResourceManagerLifecycle(t *testing.T) {
	mgr := NewResourceManager()

	res := &UniversalResource{
		ID:          "res-site-101",
		Name:        "example.com",
		Type:        "site",
		Provider:    "local",
		OrgID:       "org-main",
		ProjectID:   "proj-web",
		Environment: "production",
		Health:      "healthy",
	}

	if err := mgr.Register(res); err != nil {
		t.Fatalf("failed to register resource: %v", err)
	}

	fetched, ok := mgr.Get("res-site-101")
	if !ok || fetched.Name != "example.com" {
		t.Fatalf("failed to retrieve registered resource")
	}

	results := mgr.Search("org-main", "proj-web", "site", "healthy")
	if len(results) != 1 {
		t.Errorf("expected 1 search result, got %d", len(results))
	}

	if err := mgr.SoftDelete("res-site-101", "admin"); err != nil {
		t.Fatalf("failed to soft delete: %v", err)
	}

	activeResults := mgr.Search("org-main", "proj-web", "site", "")
	if len(activeResults) != 0 {
		t.Errorf("expected 0 active results after soft deletion")
	}
}
