package capsule

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// CapsuleManifest describes an application package.
type CapsuleManifest struct {
	CapsuleVersion string                 `json:"capsule_version"` // e.g. "1.0"
	AppID          string                 `json:"app_id"`
	AppName        string                 `json:"app_name"`
	CreatedAt      string                 `json:"created_at"`
	Runtime        string                 `json:"runtime"`         // php, nodejs, python, ruby, go, static
	RuntimeVersion string                 `json:"runtime_version"` // e.g. "8.2", "20.x"
	Entrypoint     string                 `json:"entrypoint"`      // e.g. "index.php", "server.js"
	DocumentRoot   string                 `json:"document_root"`   // e.g. "public"
	Dependencies   []string               `json:"dependencies"`    // composer.json, package.json packages
	SystemPackages []string               `json:"system_packages"` // e.g. libpng-dev, imagemagick
	CronJobs       []string               `json:"cron_jobs"`
	Volumes        []string               `json:"volumes"`         // persistent data dirs: uploads, storage
	DatabaseEngine string                 `json:"database_engine"` // mariadb, postgresql, sqlite, none
	DatabaseName   string                 `json:"database_name,omitempty"`
	EnvVariables   map[string]string      `json:"env_variables"`   // with secrets redacted
	SecretRefs     []string               `json:"secret_refs"`     // keys of redacted secrets
	ChecksumSHA256 string                 `json:"checksum_sha256"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// CapsulePackager handles creating, inspecting and diffing capsules.
type CapsulePackager struct {
	AppDir string
}

// NewCapsulePackager creates a new packager.
func NewCapsulePackager(appDir string) *CapsulePackager {
	return &CapsulePackager{AppDir: appDir}
}

// DetectRuntime inspects project files to automatically infer the runtime stack.
func (p *CapsulePackager) DetectRuntime() (runtime string, version string, entrypoint string, docRoot string) {
	if _, err := os.Stat(filepath.Join(p.AppDir, "composer.json")); err == nil {
		runtime = "php"
		version = "8.2"
		entrypoint = "index.php"
		if _, err := os.Stat(filepath.Join(p.AppDir, "public")); err == nil {
			docRoot = "public"
		} else {
			docRoot = "."
		}
		return
	}

	if _, err := os.Stat(filepath.Join(p.AppDir, "package.json")); err == nil {
		runtime = "nodejs"
		version = "20"
		entrypoint = "npm start"
		docRoot = "."
		return
	}

	if _, err := os.Stat(filepath.Join(p.AppDir, "requirements.txt")); err == nil || fileExists(p.AppDir, "pyproject.toml") {
		runtime = "python"
		version = "3.11"
		entrypoint = "main.py"
		docRoot = "."
		return
	}

	if _, err := os.Stat(filepath.Join(p.AppDir, "go.mod")); err == nil {
		runtime = "go"
		version = "1.23"
		entrypoint = "main"
		docRoot = "."
		return
	}

	// Default fallback to static site
	runtime = "static"
	version = "1.0"
	entrypoint = "index.html"
	docRoot = "."
	return
}

func fileExists(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}

// RedactSecrets parses environment configurations (.env) and redacts sensitive tokens.
func (p *CapsulePackager) RedactSecrets(envContent string) (map[string]string, []string) {
	envMap := make(map[string]string)
	secretRefs := make([]string, 0)

	secretKeyPattern := regexp.MustCompile(`(?i)(password|secret|key|token|auth|credential|api_key)`)

	lines := strings.Split(envContent, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])

		if secretKeyPattern.MatchString(k) {
			envMap[k] = fmt.Sprintf("<KENPANEL_SECRET_REF:%s>", k)
			secretRefs = append(secretRefs, k)
		} else {
			envMap[k] = v
		}
	}

	return envMap, secretRefs
}

// BuildManifest inspects the local directory and generates the Capsule manifest.
func (p *CapsulePackager) BuildManifest(appName string, envContent string) *CapsuleManifest {
	rt, ver, entry, docRoot := p.DetectRuntime()
	envVars, secretRefs := p.RedactSecrets(envContent)

	manifest := &CapsuleManifest{
		CapsuleVersion: "1.0",
		AppID:          fmt.Sprintf("app-%s-%d", strings.ToLower(appName), time.Now().Unix()),
		AppName:        appName,
		CreatedAt:      time.Now().UTC().Format(time.RFC3339),
		Runtime:        rt,
		RuntimeVersion: ver,
		Entrypoint:     entry,
		DocumentRoot:   docRoot,
		Dependencies:   []string{},
		SystemPackages: []string{},
		CronJobs:       []string{},
		Volumes:        []string{"storage", "uploads"},
		DatabaseEngine: "none",
		EnvVariables:   envVars,
		SecretRefs:     secretRefs,
		Metadata:       make(map[string]interface{}),
	}

	return manifest
}

// GenerateDockerfile creates a clean container definition for the application.
func (p *CapsulePackager) GenerateDockerfile(m *CapsuleManifest) string {
	switch m.Runtime {
	case "php":
		return fmt.Sprintf(`FROM php:%s-fpm-alpine
RUN docker-php-ext-install pdo pdo_mysql opcache
COPY . /var/www/html
WORKDIR /var/www/html
EXPOSE 9000
CMD ["php-fpm"]
`, m.RuntimeVersion)

	case "nodejs":
		return fmt.Sprintf(`FROM node:%s-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --production
COPY . .
EXPOSE 3000
CMD ["%s"]
`, m.RuntimeVersion, m.Entrypoint)

	case "python":
		return fmt.Sprintf(`FROM python:%s-slim
WORKDIR /app
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
EXPOSE 8000
CMD ["python", "%s"]
`, m.RuntimeVersion, m.Entrypoint)

	default:
		return `FROM nginx:alpine
COPY . /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
`
	}
}

// GenerateCompose creates a docker-compose.yml definition for export.
func (p *CapsulePackager) GenerateCompose(m *CapsuleManifest) string {
	var sb strings.Builder
	sb.WriteString("version: '3.8'\nservices:\n")
	sb.WriteString(fmt.Sprintf("  %s:\n", strings.ToLower(m.AppName)))
	sb.WriteString("    build: .\n")
	sb.WriteString("    restart: unless-stopped\n")
	sb.WriteString("    ports:\n      - \"8080:80\"\n")
	sb.WriteString("    environment:\n")
	for k, v := range m.EnvVariables {
		sb.WriteString(fmt.Sprintf("      - %s=%s\n", k, v))
	}
	return sb.String()
}

// CompareDrift compares two manifests to detect version, package, or config differences.
type DriftReport struct {
	HasDrift       bool     `json:"has_drift"`
	RuntimeDrift   string   `json:"runtime_drift,omitempty"`
	ConfigDiffs    []string `json:"config_diffs"`
	PackageDiffs   []string `json:"package_diffs"`
}

func CompareDrift(m1, m2 *CapsuleManifest) *DriftReport {
	report := &DriftReport{
		ConfigDiffs:  make([]string, 0),
		PackageDiffs: make([]string, 0),
	}

	if m1.Runtime != m2.Runtime || m1.RuntimeVersion != m2.RuntimeVersion {
		report.HasDrift = true
		report.RuntimeDrift = fmt.Sprintf("%s:%s -> %s:%s", m1.Runtime, m1.RuntimeVersion, m2.Runtime, m2.RuntimeVersion)
	}

	for k, v1 := range m1.EnvVariables {
		if v2, ok := m2.EnvVariables[k]; !ok {
			report.HasDrift = true
			report.ConfigDiffs = append(report.ConfigDiffs, fmt.Sprintf("Removed variable: %s", k))
		} else if v1 != v2 {
			report.HasDrift = true
			report.ConfigDiffs = append(report.ConfigDiffs, fmt.Sprintf("Changed variable: %s", k))
		}
	}

	for k := range m2.EnvVariables {
		if _, ok := m1.EnvVariables[k]; !ok {
			report.HasDrift = true
			report.ConfigDiffs = append(report.ConfigDiffs, fmt.Sprintf("Added variable: %s", k))
		}
	}

	return report
}
