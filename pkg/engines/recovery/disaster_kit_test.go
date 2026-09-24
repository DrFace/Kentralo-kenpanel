package recovery

import (
	"strings"
	"testing"
)

func TestDisasterKitRunbookGeneration(t *testing.T) {
	manifest := DisasterKitManifest{
		KitID:           "test-kit-001",
		Hostname:        "web01.example.com",
		OSDistribution:  "Ubuntu 24.04 LTS",
		KenPanelVersion: "2.3.0",
		Websites: []DisasterSiteItem{
			{Domain: "site1.com", DocumentRoot: "/var/www/site1", Webserver: "nginx", Runtime: "php", RuntimeVer: "8.2"},
		},
		Databases: []DisasterDBItem{
			{Name: "site1_db", Engine: "mariadb", SizeBytes: 1024000},
		},
	}

	generator := NewDisasterKitGenerator(manifest)
	generator.ComputeChecksum()
	htmlContent := generator.GenerateHTMLRunbook()

	if !strings.Contains(htmlContent, "web01.example.com") {
		t.Errorf("expected HTML to contain hostname")
	}
	if !strings.Contains(htmlContent, "site1.com") {
		t.Errorf("expected HTML to contain site domain")
	}
	if !strings.Contains(htmlContent, "site1_db") {
		t.Errorf("expected HTML to contain database name")
	}
	if !strings.Contains(htmlContent, "Offline Disaster Runbook") {
		t.Errorf("expected HTML to contain runbook title")
	}
}
