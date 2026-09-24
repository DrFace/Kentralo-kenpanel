package runtimes

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// PHPVersionManager manages isolated PHP-FPM versions and user pools.
type PHPVersionManager struct {
	availableVersions []string
}

func NewPHPVersionManager() *PHPVersionManager {
	return &PHPVersionManager{
		availableVersions: []string{"7.4", "8.1", "8.2", "8.3", "8.4"},
	}
}

// CreatePool provisions an isolated PHP-FPM worker pool for a user/site.
func (m *PHPVersionManager) CreatePool(ctx context.Context, version, user, siteDomain string, maxChildren int) (string, error) {
	poolDir := fmt.Sprintf("/etc/php/%s/fpm/pool.d", version)
	socketPath := fmt.Sprintf("/run/php/php%s-fpm-%s.sock", version, user)
	poolConfigPath := filepath.Join(poolDir, fmt.Sprintf("%s.conf", user))

	content := fmt.Sprintf(`[%s]
user = %s
group = %s
listen = %s
listen.owner = www-data
listen.group = www-data
listen.mode = 0660

pm = dynamic
pm.max_children = %d
pm.start_servers = 2
pm.min_spare_servers = 1
pm.max_spare_servers = 3
pm.max_requests = 500

php_admin_value[memory_limit] = 256M
php_admin_value[upload_max_filesize] = 64M
php_admin_value[post_max_size] = 64M
php_admin_value[max_execution_time] = 300
`, user, user, user, socketPath, maxChildren)

	if err := os.MkdirAll(poolDir, 0755); err != nil {
		return "", fmt.Errorf("failed creating pool dir: %w", err)
	}

	if err := os.WriteFile(poolConfigPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed writing php-fpm pool config: %w", err)
	}

	// Reload PHP-FPM service
	serviceName := fmt.Sprintf("php%s-fpm", version)
	_ = exec.CommandContext(ctx, "systemctl", "reload", serviceName).Run()

	return socketPath, nil
}

// RemovePool deletes a user pool and reloads the respective PHP-FPM service.
func (m *PHPVersionManager) RemovePool(ctx context.Context, version, user string) error {
	poolConfigPath := fmt.Sprintf("/etc/php/%s/fpm/pool.d/%s.conf", version, user)
	_ = os.Remove(poolConfigPath)

	serviceName := fmt.Sprintf("php%s-fpm", version)
	return exec.CommandContext(ctx, "systemctl", "reload", serviceName).Run()
}
