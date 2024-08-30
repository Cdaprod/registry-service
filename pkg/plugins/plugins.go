package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/Cdaprod/registry-service/internal/registry"
)

// PluginLoader is responsible for loading and registering plugins
type PluginLoader struct {
	registry   registry.Registry
	pluginsDir string
}

// NewPluginLoader creates a new PluginLoader instance
func NewPluginLoader(reg registry.Registry, pluginsDir string) *PluginLoader {
	return &PluginLoader{
		registry:   reg,
		pluginsDir: pluginsDir,
	}
}

// LoadAll dynamically loads all plugins from the specified directory
func (pl *PluginLoader) LoadAll() error {
	err := filepath.Walk(pl.pluginsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".so" {
			return nil // Skip non-shared object files
		}

		// Open the plugin
		p, err := plugin.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open plugin: %v", err)
		}

		// Look up the Register function
		symRegister, err := p.Lookup("Register")
		if err != nil {
			return fmt.Errorf("failed to find Register function in %v: %v", path, err)
		}

		// Assert that the Register function has the correct signature
		registerFunc, ok := symRegister.(func(reg registry.Registry) error)
		if !ok {
			return fmt.Errorf("invalid Register function signature in plugin: %v", path)
		}

		// Call the Register function to register the plugin with the registry
		if err := registerFunc(pl.registry); err != nil {
			return fmt.Errorf("failed to register plugin: %v", err)
		}

		return nil
	})

	return err
}

// LoadPlugin dynamically loads a single plugin by file path
func (pl *PluginLoader) LoadPlugin(pluginPath string) error {
	if filepath.Ext(pluginPath) != ".so" {
		return fmt.Errorf("invalid plugin file: %s", pluginPath)
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %v", err)
	}

	symRegister, err := p.Lookup("Register")
	if err != nil {
		return fmt.Errorf("failed to find Register function: %v", err)
	}

	registerFunc, ok := symRegister.(func(reg registry.Registry) error)
	if !ok {
		return fmt.Errorf("invalid Register function signature in plugin: %v", pluginPath)
	}

	if err := registerFunc(pl.registry); err != nil {
		return fmt.Errorf("failed to register plugin: %v", err)
	}

	return nil
}
