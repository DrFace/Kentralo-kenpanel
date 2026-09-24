package dns

import (
	"strings"
	"testing"
)

func TestDNSZoneCreationAndBINDExport(t *testing.T) {
	mgr := NewDNSManager()
	zone := mgr.CreateZone("example.org", "ns1.example.org.", "203.0.113.10")

	if zone.Domain != "example.org" {
		t.Errorf("domain mismatch")
	}

	ds, err := mgr.EnableDNSSEC("example.org")
	if err != nil {
		t.Fatalf("unexpected error enabling DNSSEC: %v", err)
	}
	if !strings.Contains(ds, "IN DS") {
		t.Errorf("expected DS record")
	}

	bindZone, err := mgr.ExportBINDZone("example.org")
	if err != nil {
		t.Fatalf("unexpected error exporting BIND: %v", err)
	}
	if !strings.Contains(bindZone, "SOA ns1.example.org.") {
		t.Errorf("expected SOA in BIND zone output")
	}
}
