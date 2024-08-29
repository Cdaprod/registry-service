package plugins

import (
    "github.com/Cdaprod/registry-service/internal/registry"
    "fmt"
)

// APIPlugin implements the Plugin interface
type APIPlugin struct{}

func (p *APIPlugin) Register(reg registry.Registry) error {
    api := &BuiltinAPI{ID: "api", Type: "API", Name: "Generic API"}
    if err := reg.Register(api); err != nil {
        return fmt.Errorf("failed to register API plugin: %w", err)
    }
    return nil
}