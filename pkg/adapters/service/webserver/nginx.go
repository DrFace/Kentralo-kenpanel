package webserver

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// NginxAdapter manages Nginx reverse proxy virtual hosts and upstream configurations.
type NginxAdapter struct {
	configDir    string
	sitesEnabled string
}

func NewNginxAdapter() *NginxAdapter {
	return &NginxAdapter{
		configDir:    "/etc/nginx/sites-available",
		sitesEnabled: "/etc/nginx/sites-enabled",
	}
}

func (a *NginxAdapter) Name() string {
	return "nginx"
}

func (a *NginxAdapter) TestConfig(ctx context.Context) (bool, string, error) {
	cmd := exec.CommandContext(ctx, "nginx", "-t")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, string(output), nil
	}
	return true, string(output), nil
}

func (a *NginxAdapter) DeployVHost(ctx context.Context, domain string, configContent []byte) error {
	vhostPath := filepath.Join(a.configDir, domain+".conf")
	symlinkPath := filepath.Join(a.sitesEnabled, domain+".conf")

	// 1. Write config
	if err := os.WriteFile(vhostPath, configContent, 0644); err != nil {
		return fmt.Errorf("failed to write nginx vhost: %w", err)
	}

	// 2. Enable symlink
	_ = os.Remove(symlinkPath)
	if err := os.Symlink(vhostPath, symlinkPath); err != nil {
		return fmt.Errorf("failed creating nginx enabled symlink: %w", err)
	}

	// 3. Test config
	valid, out, err := a.TestConfig(ctx)
	if !valid || err != nil {
		_ = os.Remove(symlinkPath)
		return fmt.Errorf("nginx config syntax validation failed: %s", out)
	}

	return nil
}

func (a *NginxAdapter) RemoveVHost(ctx context.Context, domain string) error {
	symlinkPath := filepath.Join(a.sitesEnabled, domain+".conf")
	vhostPath := filepath.Join(a.configDir, domain+".conf")

	_ = os.Remove(symlinkPath)
	_ = os.Remove(vhostPath)

	return a.Reload(ctx)
}

func (a *NginxAdapter) Reload(ctx context.Context) error {
	return exec.CommandContext(ctx, "systemctl", "reload", "nginx").Run()
}

// GenerateReverseProxyTemplate renders an Nginx vhost proxying to Apache or PHP-FPM.
func GenerateReverseProxyTemplate(domain, rootDir, upstreamAddr string, sslCert, sslKey string) string {
	sslBlock := ""
	if sslCert != "" && sslKey != "" {
		sslBlock = fmt.Sprintf(`
    listen 443 ssl http2;
    ssl_certificate %s;
    ssl_certificate_key %s;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
`, sslCert, sslKey)
	}

	return fmt.Sprintf(`server {
    listen 80;
    server_name %s;
    %s

    root %s;
    index index.php index.html;

    location / {
        proxy_pass http://%s;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location ~ /\.ht {
        deny all;
    }
}
`, domain, sslBlock, rootDir, upstreamAddr)
}
