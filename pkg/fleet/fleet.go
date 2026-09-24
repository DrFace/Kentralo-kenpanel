package fleet

import (
	"fmt"
	"sync"
	"time"
)

// FleetNode represents a registered server in the cluster.
type FleetNode struct {
	NodeID       string            `json:"node_id"`
	Hostname     string            `json:"hostname"`
	Address      string            `json:"address"` // IP or FQDN
	Port         int               `json:"port"`    // agent RPC port (e.g. 2084)
	Role         string            `json:"role"`    // controller, worker, database, edge
	Status       string            `json:"status"`  // healthy, unreachable, draining
	LastSeen     time.Time         `json:"last_seen"`
	Labels       map[string]string `json:"labels"`
	CapacityCPUs int               `json:"capacity_cpus"`
	CapacityRAM  int64             `json:"capacity_ram_mb"`
}

// FleetManager coordinates multi-server operations and node state.
type FleetManager struct {
	mu    sync.RWMutex
	nodes map[string]*FleetNode
}

// NewFleetManager initializes fleet manager.
func NewFleetManager() *FleetManager {
	return &FleetManager{
		nodes: make(map[string]*FleetNode),
	}
}

// RegisterNode enrolls a node into the cluster.
func (f *FleetManager) RegisterNode(node *FleetNode) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if node.NodeID == "" || node.Address == "" {
		return fmt.Errorf("node ID and address are required")
	}

	node.Status = "healthy"
	node.LastSeen = time.Now().UTC()
	f.nodes[node.NodeID] = node
	return nil
}

// Heartbeat updates node liveness.
func (f *FleetManager) Heartbeat(nodeID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	node, ok := f.nodes[nodeID]
	if !ok {
		return fmt.Errorf("node %s not recognized", nodeID)
	}

	node.LastSeen = time.Now().UTC()
	node.Status = "healthy"
	return nil
}

// ListNodes returns all enrolled nodes.
func (f *FleetManager) ListNodes() []*FleetNode {
	f.mu.RLock()
	defer f.mu.RUnlock()

	list := make([]*FleetNode, 0, len(f.nodes))
	for _, n := range f.nodes {
		list = append(list, n)
	}
	return list
}
