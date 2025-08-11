package jobs

import "github.com/hedgehog125/project-reboot/jobs/jobscommon"

type RegistryGroup struct {
	Registry *Registry
	Path     string
}

func (registry *Registry) Group(relativePath string) *RegistryGroup {
	return &RegistryGroup{
		Path:     relativePath,
		Registry: registry,
	}
}
func (group *RegistryGroup) Group(relativePath string) *RegistryGroup {
	return &RegistryGroup{
		Path:     jobscommon.JoinPaths(group.Path, relativePath),
		Registry: group.Registry,
	}
}

func (group *RegistryGroup) Register(action *Definition) {
	action.ID = jobscommon.JoinPaths(group.Path, action.ID)
	group.Registry.Register(action)
}
