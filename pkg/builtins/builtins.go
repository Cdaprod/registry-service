package builtins

import (
    "fmt"
    "os"
    "path/filepath"
    "plugin"
    "github.com/Cdaprod/registry-service/internal/registry"
)

// BuiltinLoader manages loading and registering built-in plugins
type BuiltinLoader struct {
    registry   registry.Registry
    pluginsDir string
}

// NewBuiltinLoader initializes a new BuiltinLoader with the registry and plugins directory
func NewBuiltinLoader(reg registry.Registry, pluginsDir string) *BuiltinLoader {
    return &BuiltinLoader{
        registry:   reg,
        pluginsDir: pluginsDir,
    }
}

// LoadAll loads and registers all built-in plugins from the specified directory
func (bl *BuiltinLoader) LoadAll() error {
    err := filepath.Walk(bl.pluginsDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if filepath.Ext(path) != ".so" {
            return nil // Skip non-plugin files
        }

        p, err := plugin.Open(path)
        if err != nil {
            return fmt.Errorf("failed to open plugin: %v", err)
        }

        symRegister, err := p.Lookup("Register")
        if err != nil {
            return fmt.Errorf("failed to find Register function in %v: %v", path, err)
        }

        registerFunc, ok := symRegister.(func(reg registry.Registry) error)
        if !ok {
            return fmt.Errorf("invalid Register function signature in plugin: %v", path)
        }

        if err := registerFunc(bl.registry); err != nil {
            return fmt.Errorf("failed to register built-in plugin: %v", err)
        }

        return nil
    })

    return err
}