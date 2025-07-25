package docker

type ServiceGroup struct {
	DockerContainerService
	DockerImageService
	DockerNetworkService
	DockerVolumeService
	DockerRegistryService
	RegistrySyncService
}
