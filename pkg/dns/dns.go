package dns

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// RecordType defines standard DNS resource record types.
type RecordType string

const (
	TypeA     RecordType = "A"
	TypeAAAA  RecordType = "AAAA"
	TypeCNAME RecordType = "CNAME"
	TypeMX    RecordType = "MX"
	TypeTXT   RecordType = "TXT"
	TypeSRV   RecordType = "SRV"
	TypeCAA   RecordType = "CAA"
	TypeNS    RecordType = "NS"
	TypePTR   RecordType = "PTR"
)

// DNSRecord represents an individual DNS entry.
type DNSRecord struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"` // "@" or "mail" or "subdomain"
	Type     RecordType `json:"type"`
	Value    string     `json:"value"`
	TTL      int        `json:"ttl"`
	Priority int        `json:"priority,omitempty"` // For MX and SRV
}

// DNSZone represents a complete zone file.
type DNSZone struct {
	Domain        string      `json:"domain"`
	SOAEmail      string      `json:"soa_email"`
	Serial        uint32      `json:"serial"`
	PrimaryNS     string      `json:"primary_ns"`
	SecondaryNS   []string    `json:"secondary_ns"`
	DNSSECEnabled bool        `json:"dnssec_enabled"`
	DSRecord      string      `json:"ds_record,omitempty"`
	Records       []DNSRecord `json:"records"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// DNSManager manages local BIND/PowerDNS zones and external provider sync.
type DNSManager struct {
	mu    sync.RWMutex
	zones map[string]*DNSZone
}

// NewDNSManager creates a DNS controller.
func NewDNSManager() *DNSManager {
	return &DNSManager{
		zones: make(map[string]*DNSZone),
	}
}

// CreateZone initializes a forward zone with default apex, MX, and NS records.
func (m *DNSManager) CreateZone(domain, primaryNS, mailServerIP string) *DNSZone {
	m.mu.Lock()
	defer m.mu.Unlock()

	domain = strings.ToLower(strings.TrimSuffix(domain, "."))
	serial := uint32(time.Now().Unix())

	zone := &DNSZone{
		Domain:        domain,
		SOAEmail:      fmt.Sprintf("admin.%s.", domain),
		Serial:        serial,
		PrimaryNS:     primaryNS,
		SecondaryNS:   []string{primaryNS},
		DNSSECEnabled: false,
		Records: []DNSRecord{
			{ID: "rec-1", Name: "@", Type: TypeA, Value: mailServerIP, TTL: 3600},
			{ID: "rec-2", Name: "www", Type: TypeCNAME, Value: domain + ".", TTL: 3600},
			{ID: "rec-3", Name: "mail", Type: TypeA, Value: mailServerIP, TTL: 3600},
			{ID: "rec-4", Name: "@", Type: TypeMX, Value: fmt.Sprintf("mail.%s.", domain), Priority: 10, TTL: 3600},
			{ID: "rec-5", Name: "@", Type: TypeTXT, Value: fmt.Sprintf("v=spf1 mx a:%s ~all", mailServerIP), TTL: 3600},
		},
		UpdatedAt: time.Now().UTC(),
	}

	m.zones[domain] = zone
	return zone
}

// EnableDNSSEC generates DNSSEC key-signing and DS records.
func (m *DNSManager) EnableDNSSEC(domain string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	zone, ok := m.zones[domain]
	if !ok {
		return "", fmt.Errorf("zone %s not found", domain)
	}

	zone.DNSSECEnabled = true
	// Mock DS record generation (Algorithm 13: ECDSA Curve P-256 with SHA-256)
	zone.DSRecord = fmt.Sprintf("%s. IN DS 2371 13 2 49FD...MOCK_HASH...", domain)
	return zone.DSRecord, nil
}

// ExportBINDZone renders standard RFC 1035 zone file syntax.
func (m *DNSManager) ExportBINDZone(domain string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	zone, ok := m.zones[domain]
	if !ok {
		return "", fmt.Errorf("zone %s not found", domain)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("$ORIGIN %s.\n$TTL 3600\n", zone.Domain))
	sb.WriteString(fmt.Sprintf("@ IN SOA %s %s (\n\t%d ; serial\n\t7200 ; refresh\n\t3600 ; retry\n\t1209600 ; expire\n\t3600 ; min ttl\n)\n",
		zone.PrimaryNS, zone.SOAEmail, zone.Serial))

	for _, r := range zone.Records {
		if r.Type == TypeMX {
			sb.WriteString(fmt.Sprintf("%s\tIN\tMX\t%d\t%s\n", r.Name, r.Priority, r.Value))
		} else {
			sb.WriteString(fmt.Sprintf("%s\tIN\t%s\t%s\n", r.Name, r.Type, r.Value))
		}
	}

	return sb.String(), nil
}
