package containers

import (
	"context"
	"testing"
)

func TestDockerManagerLifecycle(t *testing.T) {
	mgr := NewDockerManager()

	stack, err := mgr.DeployStack(context.Background(), "my-stack", "/var/www/apps/my-stack", "version: '3.8'")
	if err != nil {
		t.Fatalf("unexpected error deploying stack: %v", err)
	}

	if !stack.IsActive || stack.Name != "my-stack" {
		t.Errorf("stack deployment state mismatch")
	}
}
