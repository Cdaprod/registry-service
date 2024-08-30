package plugins

import (
	"fmt"

	"github.com/Cdaprod/registry-service/internal/registry"
)

// DockerPlugin implements the Plugin interface
type DockerPlugin struct{}

func (p *DockerPlugin) Register(reg registry.Registry) error {
	dockerAPI := &BuiltinAPI{ID: "docker", Type: "API", Name: "Docker API"}
	if err := reg.Register(dockerAPI); err != nil {
		return fmt.Errorf("failed to register Docker plugin: %w", err)
	}
	return nil
}
