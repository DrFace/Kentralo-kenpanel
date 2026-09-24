package recovery

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strconv"
	"strings"
)

// DisasterKitManifest represents the full machine-readable disaster inventory.
type DisasterKitManifest struct {
	KitID           string                 `json:"kit_id"`
	GeneratedAt     string                 `json:"generated_at"`
	Hostname        string                 `json:"hostname"`
	OSDistribution  string                 `json:"os_distribution"`
	KenPanelVersion string                 `json:"kenpanel_version"`
	RestoreOrder    []string               `json:"restore_order"`
	Websites        []DisasterSiteItem     `json:"websites"`
	Databases       []DisasterDBItem       `json:"databases"`
	Mailboxes       []DisasterMailItem     `json:"mailboxes"`
	SystemPackages  []string               `json:"system_packages"`
	EncryptedVault  string                 `json:"encrypted_vault,omitempty"` // AES-GCM ciphertext if user opted-in
	ChecksumSHA256  string                 `json:"checksum_sha256"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type DisasterSiteItem struct {
	Domain       string `json:"domain"`
	DocumentRoot string `json:"document_root"`
	Webserver    string `json:"webserver"`
	Runtime      string `json:"runtime"`
	RuntimeVer   string `json:"runtime_ver"`
	DatabaseName string `json:"database_name,omitempty"`
}

type DisasterDBItem struct {
	Name      string `json:"name"`
	Engine    string `json:"engine"`
	SizeBytes int64  `json:"size_bytes"`
}

type DisasterMailItem struct {
	Domain        string   `json:"domain"`
	TotalAccounts int      `json:"total_accounts"`
	Mailboxes     []string `json:"mailboxes"`
}

// DisasterKitGenerator builds runbooks and verification manifests.
type DisasterKitGenerator struct {
	Manifest DisasterKitManifest
}

// NewDisasterKitGenerator initializes generator.
func NewDisasterKitGenerator(m DisasterKitManifest) *DisasterKitGenerator {
	if len(m.RestoreOrder) == 0 {
		m.RestoreOrder = []string{
			"1. Provision host OS and base packages",
			"2. Install KenPanel agent / runtime dependencies",
			"3. Restore database engines and schemas",
			"4. Reconstruct web server virtual hosts and SSL certificates",
			"5. Synchronize document roots and user files",
			"6. Restore mail domains and mailbox hierarchies",
			"7. Perform DNS cutover and execute WebStack Doctor validation",
		}
	}
	return &DisasterKitGenerator{Manifest: m}
}

// GenerateHTMLRunbook renders a self-contained, offline HTML runbook with zero external dependencies.
func (g *DisasterKitGenerator) GenerateHTMLRunbook() string {
	m := g.Manifest
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>KenPanel Disaster Recovery Runbook — `)
	sb.WriteString(html.EscapeString(m.Hostname))
	sb.WriteString(`</title>
<style>
  :root { --bg: #090d16; --card: #111827; --text: #f3f4f6; --muted: #9ca3af; --accent: #2563eb; --border: #1f2937; --success: #10b981; }
  body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: var(--bg); color: var(--text); margin: 0; padding: 2rem; line-height: 1.6; }
  .container { max-width: 960px; margin: 0 auto; }
  .badge { display: inline-block; padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; }
  .badge-primary { background: var(--accent); color: white; }
  .badge-success { background: rgba(16,185,129,0.2); color: var(--success); border: 1px solid var(--success); }
  .card { background: var(--card); border: 1px solid var(--border); border-radius: 0.75rem; padding: 1.5rem; margin-bottom: 1.5rem; }
  h1, h2, h3 { color: white; margin-top: 0; }
  table { width: 100%; border-collapse: collapse; margin-top: 1rem; }
  th, td { text-align: left; padding: 0.75rem; border-bottom: 1px solid var(--border); }
  th { color: var(--muted); font-size: 0.875rem; }
  code { background: #1e293b; padding: 0.2rem 0.4rem; border-radius: 0.25rem; font-family: monospace; font-size: 0.9em; }
  ol { padding-left: 1.25rem; }
  li { margin-bottom: 0.5rem; }
</style>
</head>
<body>
<div class="container">
  <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem;">
    <div>
      <span class="badge badge-primary">Offline Disaster Runbook</span>
      <h1>`)
	sb.WriteString(html.EscapeString(m.Hostname))
	sb.WriteString(`</h1>
      <p style="color: var(--muted); margin: 0;">Generated: `)
	sb.WriteString(html.EscapeString(m.GeneratedAt))
	sb.WriteString(` | KenPanel `)
	sb.WriteString(html.EscapeString(m.KenPanelVersion))
	sb.WriteString(`</p>
    </div>
    <span class="badge badge-success">Self-Contained</span>
  </div>

  <div class="card">
    <h2>1. Recommended Restoration Order</h2>
    <ol>
`)

	for _, step := range m.RestoreOrder {
		sb.WriteString("      <li>")
		sb.WriteString(html.EscapeString(step))
		sb.WriteString("</li>\n")
	}

	sb.WriteString(`    </ol>
  </div>

  <div class="card">
    <h2>2. Inventory: Websites & Applications (`)
	sb.WriteString(strconv.Itoa(len(m.Websites)))
	sb.WriteString(`)</h2>
    <table>
      <thead>
        <tr><th>Domain</th><th>Document Root</th><th>Webserver</th><th>Runtime</th><th>Database</th></tr>
      </thead>
      <tbody>
`)

	for _, w := range m.Websites {
		fmt.Fprintf(&sb, "        <tr><td><strong>%s</strong></td><td><code>%s</code></td><td>%s</td><td>%s %s</td><td><code>%s</code></td></tr>\n",
			html.EscapeString(w.Domain),
			html.EscapeString(w.DocumentRoot),
			html.EscapeString(w.Webserver),
			html.EscapeString(w.Runtime),
			html.EscapeString(w.RuntimeVer),
			html.EscapeString(w.DatabaseName),
		)
	}

	sb.WriteString(`      </tbody>
    </table>
  </div>

  <div class="card">
    <h2>3. Inventory: Databases (`)
	sb.WriteString(strconv.Itoa(len(m.Databases)))
	sb.WriteString(`)</h2>
    <table>
      <thead>
        <tr><th>Database Name</th><th>Engine</th><th>Size (Bytes)</th></tr>
      </thead>
      <tbody>
`)

	for _, db := range m.Databases {
		fmt.Fprintf(&sb, "        <tr><td><code>%s</code></td><td>%s</td><td>%d</td></tr>\n",
			html.EscapeString(db.Name),
			html.EscapeString(db.Engine),
			db.SizeBytes,
		)
	}

	sb.WriteString(`      </tbody>
    </table>
  </div>

  <div class="card">
    <h2>4. Offline Verification & Integrity</h2>
    <p>Kit Manifest SHA256 Checksum: <code>`)
	sb.WriteString(html.EscapeString(m.ChecksumSHA256))
	sb.WriteString(`</code></p>
    <p>To verify offline using KenPanel CLI:</p>
    <pre><code>kenpanel disaster-kit verify --manifest disaster-kit-manifest.json</code></pre>
  </div>
</div>
</body>
</html>`)

	return sb.String()
}

// ComputeChecksum calculates SHA256 of the manifest data.
func (g *DisasterKitGenerator) ComputeChecksum() string {
	data, _ := json.Marshal(g.Manifest)
	h := sha256.New()
	h.Write(data)
	g.Manifest.ChecksumSHA256 = hex.EncodeToString(h.Sum(nil))
	return g.Manifest.ChecksumSHA256
}

// VerifyKit validates an offline disaster kit package.
func VerifyKit(manifestPath string) (bool, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return false, fmt.Errorf("failed to read disaster kit manifest: %w", err)
	}

	var m DisasterKitManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return false, fmt.Errorf("invalid manifest JSON: %w", err)
	}

	if m.KitID == "" || m.Hostname == "" {
		return false, fmt.Errorf("corrupted manifest: missing required identity fields")
	}

	return true, nil
}
