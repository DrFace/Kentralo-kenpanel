package webserver

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ApacheAdapter manages Apache HTTP Server backend virtual hosts.
type ApacheAdapter struct {
	configDir string
}

func NewApacheAdapter() *ApacheAdapter {
	return &ApacheAdapter{
		configDir: "/etc/apache2/sites-available",
	}
}

func (a *ApacheAdapter) Name() string {
	return "apache2"
}

func (a *ApacheAdapter) TestConfig(ctx context.Context) (bool, string, error) {
	cmd := exec.CommandContext(ctx, "apachectl", "configtest")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, string(output), nil
	}
	return true, string(output), nil
}

func (a *ApacheAdapter) DeployVHost(ctx context.Context, domain string, configContent []byte) error {
	vhostPath := filepath.Join(a.configDir, domain+".conf")
	if err := os.WriteFile(vhostPath, configContent, 0644); err != nil {
		return fmt.Errorf("failed writing apache vhost: %w", err)
	}

	_ = exec.CommandContext(ctx, "a2ensite", domain+".conf").Run()

	valid, out, err := a.TestConfig(ctx)
	if !valid || err != nil {
		_ = exec.CommandContext(ctx, "a2dissite", domain+".conf").Run()
		return fmt.Errorf("apache config syntax validation failed: %s", out)
	}

	return a.Reload(ctx)
}

func (a *ApacheAdapter) RemoveVHost(ctx context.Context, domain string) error {
	_ = exec.CommandContext(ctx, "a2dissite", domain+".conf").Run()
	_ = os.Remove(filepath.Join(a.configDir, domain+".conf"))
	return a.Reload(ctx)
}

func (a *ApacheAdapter) Reload(ctx context.Context) error {
	return exec.CommandContext(ctx, "systemctl", "reload", "apache2").Run()
}

// GenerateBackendVHostTemplate renders an Apache backend virtual host listening on a local port.
func GenerateBackendVHostTemplate(domain, rootDir string, backendPort int, phpSocket string) string {
	return fmt.Sprintf(`<VirtualHost 127.0.0.1:%d>
    ServerName %s
    DocumentRoot %s

    <Directory %s>
        Options -Indexes +FollowSymLinks
        AllowOverride All
        Require all granted
    </Directory>

    <FilesMatch \.php$>
        SetHandler "proxy:unix:%s|fcgi://localhost/"
    </FilesMatch>

    ErrorLog ${APACHE_LOG_DIR}/%s_error.log
    CustomLog ${APACHE_LOG_DIR}/%s_access.log combined
</VirtualHost>
`, backendPort, domain, rootDir, rootDir, phpSocket, domain, domain)
}
