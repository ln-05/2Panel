package docker

type ServiceGroup struct {
	DockerContainerService
	DockerImageService
	DockerNetworkService
	DockerVolumeService
	DockerRegistryService
	DockerConfigService
	DockerOverviewService
	DockerDiagnosticService
}
