package plugins

import (
    "github.com/Cdaprod/registry-service/internal/registry"
    "fmt"
)

// GitPlugin implements the Plugin interface
type GitPlugin struct{}

func (p *GitPlugin) Register(reg registry.Registry) error {
    gitAPI := &BuiltinAPI{ID: "git", Type: "API", Name: "Git API"}
    if err := reg.Register(gitAPI); err != nil {
        return fmt.Errorf("failed to register Git plugin: %w", err)
    }
    return nil
}