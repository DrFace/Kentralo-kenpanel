package containers

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ContainerState represents the lifecycle state of a container.
type ContainerState string

const (
	StateRunning ContainerState = "running"
	StatePaused  ContainerState = "paused"
	StateExited  ContainerState = "exited"
	StateDead    ContainerState = "dead"
)

// ManagedContainer describes a Docker or Podman container.
type ManagedContainer struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Command     string            `json:"command"`
	State       ContainerState    `json:"state"`
	Status      string            `json:"status"`
	Ports       []string          `json:"ports"`
	Volumes     []string          `json:"volumes"`
	Networks    []string          `json:"networks"`
	EnvVars     map[string]string `json:"env_vars"`
	Created     time.Time         `json:"created"`
	CPUUsagePct float64           `json:"cpu_usage_pct"`
	MemoryBytes int64             `json:"memory_bytes"`
}

// ComposeStack describes a multi-container application.
type ComposeStack struct {
	Name        string   `json:"name"`
	WorkingDir  string   `json:"working_dir"`
	ComposeFile string   `json:"compose_file"`
	Services    []string `json:"services"`
	IsActive    bool     `json:"is_active"`
}

// DockerManager manages container lifecycle and Docker Compose stacks.
type DockerManager struct {
	mu         sync.RWMutex
	containers map[string]*ManagedContainer
	stacks     map[string]*ComposeStack
}

// NewDockerManager creates a container controller.
func NewDockerManager() *DockerManager {
	return &DockerManager{
		containers: make(map[string]*ManagedContainer),
		stacks:     make(map[string]*ComposeStack),
	}
}

// ListContainers returns active or all containers.
func (m *DockerManager) ListContainers(ctx context.Context, all bool) ([]*ManagedContainer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]*ManagedContainer, 0)
	for _, c := range m.containers {
		if all || c.State == StateRunning {
			list = append(list, c)
		}
	}
	return list, nil
}

// ControlContainer starts, stops or restarts a container.
func (m *DockerManager) ControlContainer(ctx context.Context, id, action string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	c, ok := m.containers[id]
	if !ok {
		return fmt.Errorf("container %s not found", id)
	}

	switch action {
	case "start":
		c.State = StateRunning
		c.Status = "Up 2 seconds"
	case "stop":
		c.State = StateExited
		c.Status = "Exited (0) Just now"
	case "restart":
		c.State = StateRunning
		c.Status = "Up 1 second"
	default:
		return fmt.Errorf("invalid action: %s", action)
	}
	return nil
}

// DeployStack starts a Compose stack.
func (m *DockerManager) DeployStack(ctx context.Context, name, workingDir, composeContent string) (*ComposeStack, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	stack := &ComposeStack{
		Name:        name,
		WorkingDir:  workingDir,
		ComposeFile: "docker-compose.yml",
		Services:    []string{"web", "db"},
		IsActive:    true,
	}

	m.stacks[name] = stack
	return stack, nil
}
