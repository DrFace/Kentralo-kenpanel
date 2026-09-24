package plugins

import (
	"context"
	"fmt"
	"sync"
)

// PluginManifest defines metadata and permission requirements for an extension.
type PluginManifest struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"` // e.g. "sites:read", "dns:write"
	Hooks       []string `json:"hooks"`       // e.g. "on_site_create"
	Enabled     bool     `json:"enabled"`
}

// HookHandler defines callback signature for extension points.
type HookHandler func(ctx context.Context, payload map[string]interface{}) error

// PluginEngine manages extensions and dispatches event hooks.
type PluginEngine struct {
	mu       sync.RWMutex
	plugins  map[string]*PluginManifest
	handlers map[string][]HookHandler
}

// NewPluginEngine creates the plugin controller.
func NewPluginEngine() *PluginEngine {
	return &PluginEngine{
		plugins:  make(map[string]*PluginManifest),
		handlers: make(map[string][]HookHandler),
	}
}

// RegisterPlugin registers an extension after validating its manifest.
func (e *PluginEngine) RegisterPlugin(manifest *PluginManifest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if manifest.ID == "" {
		return fmt.Errorf("plugin ID cannot be empty")
	}

	manifest.Enabled = true
	e.plugins[manifest.ID] = manifest
	return nil
}

// Subscribe attaches a callback to an event hook.
func (e *PluginEngine) Subscribe(hookName string, handler HookHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.handlers[hookName] = append(e.handlers[hookName], handler)
}

// Dispatch executes all handlers subscribed to a hook.
func (e *PluginEngine) Dispatch(ctx context.Context, hookName string, payload map[string]interface{}) []error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var errs []error
	for _, h := range e.handlers[hookName] {
		if err := h(ctx, payload); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
