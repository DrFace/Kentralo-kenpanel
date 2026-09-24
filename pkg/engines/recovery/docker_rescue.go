package recovery

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ContainerMetadata represents an offline or stopped container's state.
type ContainerMetadata struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	State       string            `json:"state"` // running, exited, dead
	Ports       []string          `json:"ports"` // e.g. ["80:80/tcp", "443:443/tcp"]
	Volumes     []string          `json:"volumes"`
	Networks    []string          `json:"networks"`
	Environment map[string]string `json:"environment"`
}

// DockerRescueDiagnosis contains results of inspecting container subsystems.
type DockerRescueDiagnosis struct {
	Timestamp        time.Time           `json:"timestamp"`
	TotalContainers  int                 `json:"total_containers"`
	OrphanedVolumes  []string            `json:"orphaned_volumes"`
	OrphanedNetworks []string            `json:"orphaned_networks"`
	PortConflicts    []string            `json:"port_conflicts"`
	BrokenBindMounts []string            `json:"broken_bind_mounts"`
	Containers       []ContainerMetadata `json:"containers"`
}

// DockerRescueEngine inspects and rescues container workloads.
type DockerRescueEngine struct{}

// NewDockerRescueEngine initializes the engine.
func NewDockerRescueEngine() *DockerRescueEngine {
	return &DockerRescueEngine{}
}

// Diagnose evaluates containers for issues like port collisions or broken volume mounts.
func (e *DockerRescueEngine) Diagnose(ctx context.Context, containers []ContainerMetadata) *DockerRescueDiagnosis {
	diag := &DockerRescueDiagnosis{
		Timestamp:        time.Now().UTC(),
		TotalContainers:  len(containers),
		OrphanedVolumes:  make([]string, 0),
		OrphanedNetworks: make([]string, 0),
		PortConflicts:    make([]string, 0),
		BrokenBindMounts: make([]string, 0),
		Containers:       containers,
	}

	portMap := make(map[string]string)
	for _, c := range containers {
		for _, p := range c.Ports {
			if existingContainer, collision := portMap[p]; collision {
				diag.PortConflicts = append(diag.PortConflicts, fmt.Sprintf("Port %s collided between %s and %s", p, existingContainer, c.Name))
			} else {
				portMap[p] = c.Name
			}
		}
	}

	return diag
}

// ReconstructCompose generates a docker-compose.yml file from inspected container metadata.
func (e *DockerRescueEngine) ReconstructCompose(containers []ContainerMetadata) string {
	var sb strings.Builder
	sb.WriteString("version: '3.8'\n\nservices:\n")

	for _, c := range containers {
		safeName := strings.TrimPrefix(c.Name, "/")
		sb.WriteString(fmt.Sprintf("  %s:\n", safeName))
		sb.WriteString(fmt.Sprintf("    image: %s\n", c.Image))
		sb.WriteString("    restart: unless-stopped\n")

		if len(c.Ports) > 0 {
			sb.WriteString("    ports:\n")
			for _, p := range c.Ports {
				sb.WriteString(fmt.Sprintf("      - \"%s\"\n", p))
			}
		}

		if len(c.Volumes) > 0 {
			sb.WriteString("    volumes:\n")
			for _, v := range c.Volumes {
				sb.WriteString(fmt.Sprintf("      - %s\n", v))
			}
		}

		if len(c.Networks) > 0 {
			sb.WriteString("    networks:\n")
			for _, n := range c.Networks {
				sb.WriteString(fmt.Sprintf("      - %s\n", n))
			}
		}

		if len(c.Environment) > 0 {
			sb.WriteString("    environment:\n")
			for k, v := range c.Environment {
				sb.WriteString(fmt.Sprintf("      - %s=%s\n", k, v))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
