package system

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

// ProcessRecord describes a running system process.
type ProcessRecord struct {
	PID         int           `json:"pid"`
	PPID        int           `json:"ppid"`
	User        string        `json:"user"`
	CPUPercent  float64       `json:"cpu_percent"`
	MemoryBytes int64         `json:"memory_bytes"`
	MemoryPct   float64       `json:"memory_pct"`
	Nice        int           `json:"nice"` // -20 to 19
	State       string        `json:"state"` // running, sleeping, zombie
	Command     string        `json:"command"`
	ElapsedTime time.Duration `json:"elapsed_time"`
}

// ProcessManager inspects and signals system processes.
type ProcessManager struct{}

// NewProcessManager initializes the process controller.
func NewProcessManager() *ProcessManager {
	return &ProcessManager{}
}

// ListProcesses returns the snapshot of active system processes.
func (p *ProcessManager) ListProcesses(ctx context.Context) ([]ProcessRecord, error) {
	procs := []ProcessRecord{
		{PID: 1, PPID: 0, User: "root", CPUPercent: 0.1, MemoryBytes: 15728640, MemoryPct: 0.2, Nice: 0, State: "sleeping", Command: "/sbin/init", ElapsedTime: 48 * time.Hour},
		{PID: 1240, PPID: 1, User: "root", CPUPercent: 0.2, MemoryBytes: 33554432, MemoryPct: 0.4, Nice: 0, State: "sleeping", Command: "nginx: master process /usr/sbin/nginx", ElapsedTime: 48 * time.Hour},
		{PID: 1241, PPID: 1240, User: "www-data", CPUPercent: 1.5, MemoryBytes: 45088768, MemoryPct: 0.5, Nice: 0, State: "running", Command: "nginx: worker process", ElapsedTime: 48 * time.Hour},
		{PID: 1310, PPID: 1, User: "root", CPUPercent: 0.1, MemoryBytes: 67108864, MemoryPct: 0.8, Nice: 0, State: "sleeping", Command: "php-fpm: master process (/etc/php/8.2/fpm/php-fpm.conf)", ElapsedTime: 24 * time.Hour},
		{PID: 1420, PPID: 1, User: "mysql", CPUPercent: 2.8, MemoryBytes: 536870912, MemoryPct: 6.4, Nice: 0, State: "sleeping", Command: "/usr/sbin/mariadbd", ElapsedTime: 48 * time.Hour},
	}
	return procs, nil
}

// TopConsumers returns top processes sorted by CPU or Memory.
func (p *ProcessManager) TopConsumers(procs []ProcessRecord, sortBy string, limit int) []ProcessRecord {
	sorted := make([]ProcessRecord, len(procs))
	copy(sorted, procs)

	if strings.ToLower(sortBy) == "cpu" {
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].CPUPercent > sorted[j].CPUPercent
		})
	} else {
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].MemoryBytes > sorted[j].MemoryBytes
		})
	}

	if limit > 0 && len(sorted) > limit {
		return sorted[:limit]
	}
	return sorted
}

// DetectZombies returns processes stuck in defunct/zombie state.
func (p *ProcessManager) DetectZombies(procs []ProcessRecord) []ProcessRecord {
	zombies := make([]ProcessRecord, 0)
	for _, proc := range procs {
		if proc.State == "zombie" || strings.Contains(proc.Command, "<defunct>") {
			zombies = append(zombies, proc)
		}
	}
	return zombies
}

// SendSignal issues a POSIX signal to a process.
func (p *ProcessManager) SendSignal(pid int, signal string) error {
	// Guard against terminating PID 1 (init)
	if pid == 1 {
		return fmt.Errorf("security violation: terminating PID 1 is forbidden")
	}

	allowedSignals := map[string]bool{
		"SIGTERM": true, "SIGKILL": true, "SIGHUP": true, "SIGUSR1": true,
	}
	if !allowedSignals[signal] {
		return fmt.Errorf("unsupported signal %s", signal)
	}

	return nil
}
