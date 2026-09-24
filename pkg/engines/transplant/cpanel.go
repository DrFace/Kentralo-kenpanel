package transplant

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// CPanelDiscoveryParser parses cPanel/WHM backups (cpmove or full backup tar.gz archives)
// and extracts user accounts, virtual hosts, databases, mailboxes, and DNS records into ManifestV1.
type CPanelDiscoveryParser struct {
	ArchivePath string
	ExtractDir  string
}

// NewCPanelDiscoveryParser creates a new parser instance.
func NewCPanelDiscoveryParser(archivePath string) *CPanelDiscoveryParser {
	return &CPanelDiscoveryParser{
		ArchivePath: archivePath,
	}
}

// Parse extracts archive metadata without modifying anything on the source server.
func (p *CPanelDiscoveryParser) Parse(ctx context.Context) (*ManifestV1, error) {
	if _, err := os.Stat(p.ArchivePath); err != nil {
		return nil, fmt.Errorf("cpanel archive not found: %w", err)
	}

	manifest := &ManifestV1{
		Version:     "1.0",
		ExportedAt:  time.Now().UTC().Format(time.RFC3339),
		SourceType:  "cpanel",
		SourceHost:  "unknown",
		Websites:    []WebsiteResource{},
		Databases:   []DatabaseResource{},
		Mailboxes:   []MailboxResource{},
		DNSZones:    []DNSZoneResource{},
		Metadata:    make(map[string]interface{}),
	}

	file, err := os.Open(p.ArchivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	var tr *tar.Reader
	if strings.HasSuffix(p.ArchivePath, ".tar.gz") || strings.HasSuffix(p.ArchivePath, ".tgz") {
		gz, err := gzip.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize gzip reader: %w", err)
		}
		defer gz.Close()
		tr = tar.NewReader(gz)
	} else {
		tr = tar.NewReader(file)
	}

	username := ""
	cpuserdata := make(map[string]string)
	mysqlDumps := make([]string, 0)
	emailAccounts := make([]string, 0)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading tar archive: %w", err)
		}

		name := header.Name
		// Clean leading path ./ or cpmove-username/
		parts := strings.Split(filepath.ToSlash(name), "/")
		if len(parts) > 1 && strings.HasPrefix(parts[0], "cpmove-") {
			username = strings.TrimPrefix(parts[0], "cpmove-")
		}

		// Detect cpuser / userdata
		if strings.Contains(name, "userdata/main") || strings.HasSuffix(name, "cp/user") {
			buf := new(strings.Builder)
			io.Copy(buf, tr)
			p.parseYAMLKeyValue(buf.String(), cpuserdata)
		}

		// Detect SQL dumps
		if strings.Contains(name, "mysql/") && strings.HasSuffix(name, ".sql") {
			dbName := filepath.Base(name)
			dbName = strings.TrimSuffix(dbName, ".sql")
			mysqlDumps = append(mysqlDumps, dbName)
		}

		// Detect DNS zones
		if strings.Contains(name, "dnszones/") && !header.FileInfo().IsDir() {
			domain := filepath.Base(name)
			buf := new(strings.Builder)
			io.Copy(buf, tr)
			manifest.DNSZones = append(manifest.DNSZones, DNSZoneResource{
				Domain:  domain,
				RawZone: buf.String(),
			})
		}

		// Detect email shadow / passwd
		if strings.Contains(name, "homedir/etc/") && strings.HasSuffix(name, "/shadow") {
			parts := strings.Split(name, "/")
			for i, p := range parts {
				if p == "etc" && i+1 < len(parts) {
					emailDomain := parts[i+1]
					scanner := bufio.NewScanner(tr)
					for scanner.Scan() {
						line := strings.TrimSpace(scanner.Text())
						if line != "" {
							user := strings.Split(line, ":")[0]
							emailAccounts = append(emailAccounts, user+"@"+emailDomain)
						}
					}
				}
			}
		}
	}

	if username == "" {
		if u, ok := cpuserdata["user"]; ok {
			username = u
		} else {
			username = "cpanel_imported"
		}
	}

	manifest.SourceHost = username

	// Construct Website resource
	mainDomain := cpuserdata["main_domain"]
	if mainDomain == "" {
		mainDomain = username + ".cpanel.local"
	}

	docRoot := cpuserdata["document_root"]
	if docRoot == "" {
		docRoot = "/home/" + username + "/public_html"
	}

	manifest.Websites = append(manifest.Websites, WebsiteResource{
		ID:           fmt.Sprintf("web-%s", username),
		Domain:       mainDomain,
		DocumentRoot: docRoot,
		Runtime:      "php",
		RuntimeVer:   "8.2",
		Webserver:    "nginx",
		SSL:          false,
	})

	// Construct Databases
	for _, db := range mysqlDumps {
		manifest.Databases = append(manifest.Databases, DatabaseResource{
			ID:          fmt.Sprintf("db-%s", db),
			Name:        db,
			Engine:      "mariadb",
			Charset:     "utf8mb4",
			DumpPath:    fmt.Sprintf("mysql/%s.sql", db),
			SizeBytes:   0,
		})
	}

	// Construct Mailboxes
	for _, email := range emailAccounts {
		parts := strings.Split(email, "@")
		manifest.Mailboxes = append(manifest.Mailboxes, MailboxResource{
			ID:        fmt.Sprintf("mail-%s", email),
			Email:     email,
			Domain:    parts[1],
			QuotaMB:   1024,
			Format:    "maildir",
		})
	}

	manifest.Metadata["total_websites"] = len(manifest.Websites)
	manifest.Metadata["total_databases"] = len(manifest.Databases)
	manifest.Metadata["total_mailboxes"] = len(manifest.Mailboxes)
	manifest.Metadata["total_dns_zones"] = len(manifest.DNSZones)

	return manifest, nil
}

func (p *CPanelDiscoveryParser) parseYAMLKeyValue(content string, target map[string]string) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			target[k] = v
		}
	}
}

// PleskDiscoveryParser parses Plesk XML and backup dumps
type PleskDiscoveryParser struct {
	ArchivePath string
}

// NewPleskDiscoveryParser creates a new Plesk parser.
func NewPleskDiscoveryParser(archivePath string) *PleskDiscoveryParser {
	return &PleskDiscoveryParser{ArchivePath: archivePath}
}

// Parse extracts Plesk domains, clients, databases and mail accounts
func (p *PleskDiscoveryParser) Parse(ctx context.Context) (*ManifestV1, error) {
	if _, err := os.Stat(p.ArchivePath); err != nil {
		return nil, fmt.Errorf("plesk archive not found: %w", err)
	}

	manifest := &ManifestV1{
		Version:     "1.0",
		ExportedAt:  time.Now().UTC().Format(time.RFC3339),
		SourceType:  "plesk",
		SourceHost:  "plesk-export",
		Websites:    []WebsiteResource{},
		Databases:   []DatabaseResource{},
		Mailboxes:   []MailboxResource{},
		DNSZones:    []DNSZoneResource{},
		Metadata:    make(map[string]interface{}),
	}

	// For standard Plesk backup archives (tar, zip, or xml dump)
	file, err := os.Open(p.ArchivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plesk archive: %w", err)
	}
	defer file.Close()

	// Regex heuristics for Plesk dump XML
	domainRegex := regexp.MustCompile(`<domain name="([^"]+)"`)
	docRootRegex := regexp.MustCompile(`<www-root>([^<]+)</www-root>`)
	dbRegex := regexp.MustCompile(`<database name="([^"]+)" type="([^"]+)"`)
	mailRegex := regexp.MustCompile(`<mail-name name="([^"]+)"`)

	scanner := bufio.NewScanner(file)
	currentDomain := ""
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		line := scanner.Text()

		if m := domainRegex.FindStringSubmatch(line); len(m) > 1 {
			currentDomain = m[1]
			manifest.Websites = append(manifest.Websites, WebsiteResource{
				ID:           fmt.Sprintf("web-%s", currentDomain),
				Domain:       currentDomain,
				DocumentRoot: fmt.Sprintf("/var/www/vhosts/%s/httpdocs", currentDomain),
				Runtime:      "php",
				RuntimeVer:   "8.2",
				Webserver:    "nginx",
				SSL:          true,
			})
		}

		if m := docRootRegex.FindStringSubmatch(line); len(m) > 1 && len(manifest.Websites) > 0 {
			manifest.Websites[len(manifest.Websites)-1].DocumentRoot = m[1]
		}

		if m := dbRegex.FindStringSubmatch(line); len(m) > 2 {
			dbName := m[1]
			dbType := m[2]
			manifest.Databases = append(manifest.Databases, DatabaseResource{
				ID:          fmt.Sprintf("db-%s", dbName),
				Name:        dbName,
				Engine:      dbType,
				Charset:     "utf8mb4",
				SizeBytes:   0,
			})
		}

		if m := mailRegex.FindStringSubmatch(line); len(m) > 1 && currentDomain != "" {
			mailUser := m[1]
			fullEmail := fmt.Sprintf("%s@%s", mailUser, currentDomain)
			manifest.Mailboxes = append(manifest.Mailboxes, MailboxResource{
				ID:      fmt.Sprintf("mail-%s", fullEmail),
				Email:   fullEmail,
				Domain:  currentDomain,
				QuotaMB: 2048,
				Format:  "maildir",
			})
		}
	}

	manifest.Metadata["total_websites"] = len(manifest.Websites)
	manifest.Metadata["total_databases"] = len(manifest.Databases)
	manifest.Metadata["total_mailboxes"] = len(manifest.Mailboxes)

	return manifest, nil
}

// SaveManifestJSON writes a manifest to disk.
func SaveManifestJSON(manifest *ManifestV1, targetPath string) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(targetPath, data, 0644)
}
