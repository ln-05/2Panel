package docker

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	DockerContainerApi
	DockerImageApi
	DockerNetworkApi
	DockerVolumeApi
	DockerRegistryApi
	DockerConfigApi
	DockerOverviewApi
	DockerDiagnosticApi
}

var (
	dockerContainerService  = service.ServiceGroupApp.DockerServiceGroup.DockerContainerService
	dockerImageService      = service.ServiceGroupApp.DockerServiceGroup.DockerImageService
	dockerNetworkService    = service.ServiceGroupApp.DockerServiceGroup.DockerNetworkService
	dockerVolumeService     = service.ServiceGroupApp.DockerServiceGroup.DockerVolumeService
	dockerRegistryService   = service.ServiceGroupApp.DockerServiceGroup.DockerRegistryService
	dockerConfigService     = service.ServiceGroupApp.DockerServiceGroup.DockerConfigService
	dockerOverviewService   = service.ServiceGroupApp.DockerServiceGroup.DockerOverviewService
	dockerDiagnosticService = service.ServiceGroupApp.DockerServiceGroup.DockerDiagnosticService
)
