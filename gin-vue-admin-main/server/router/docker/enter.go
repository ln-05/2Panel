package docker

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	DockerContainerRouter
	DockerImageRouter
	DockerNetworkRouter
	DockerVolumeRouter
	DockerRegistryRouter
}

var (
	dockerContainerApi = api.ApiGroupApp.DockerApiGroup.DockerContainerApi
	dockerImageApi     = api.ApiGroupApp.DockerApiGroup.DockerImageApi
)
