package twofactoractions

import (
	"fmt"
)

type RegistryGroup struct {
	Registry *Registry
	path     string
}

func (registry *Registry) Group(relativePath string) *RegistryGroup {
	return &RegistryGroup{
		path:     relativePath,
		Registry: registry,
	}
}
func (group *RegistryGroup) Group(relativePath string) *RegistryGroup {
	return &RegistryGroup{
		path:     joinPaths(group.path, relativePath),
		Registry: group.Registry,
	}
}

func (group *RegistryGroup) RegisterAction(action *ActionDefinition) {
	action.ID = joinPaths(group.path, action.ID)
	group.Registry.RegisterAction(action)
}

func joinPaths(path1 string, path2 string) string {
	if path1 == "" {
		return path2
	}
	if path2 == "" {
		return path1
	}
	return fmt.Sprintf("%s/%s", path1, path2)
}
