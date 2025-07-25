package docker

import (
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/flipped-aurora/gin-vue-admin/server/model/docker/response"
)

// ConvertToContainerInfo 将Docker API的Container类型转换为ContainerInfo响应模型
func ConvertToContainerInfo(dockerContainer types.Container) response.ContainerInfo {
	// 处理容器名称（移除前缀斜杠）
	name := ""
	if len(dockerContainer.Names) > 0 {
		name = strings.TrimPrefix(dockerContainer.Names[0], "/")
	}

	// 转换端口映射
	ports := make([]response.PortMapping, 0, len(dockerContainer.Ports))
	for _, port := range dockerContainer.Ports {
		ports = append(ports, response.PortMapping{
			PrivatePort: int(port.PrivatePort),
			PublicPort:  int(port.PublicPort),
			Type:        port.Type,
			IP:          port.IP,
		})
	}

	return response.ContainerInfo{
		ID:         dockerContainer.ID,
		Name:       name,
		Image:      dockerContainer.Image,
		ImageID:    dockerContainer.ImageID,
		Command:    dockerContainer.Command,
		Created:    dockerContainer.Created,
		Status:     dockerContainer.Status,
		State:      dockerContainer.State,
		Ports:      ports,
		Labels:     dockerContainer.Labels,
		SizeRw:     dockerContainer.SizeRw,
		SizeRootFs: dockerContainer.SizeRootFs,
	}
}

// ConvertToContainerDetail 将Docker API的ContainerJSON转换为ContainerDetail响应模型
func ConvertToContainerDetail(dockerContainer types.ContainerJSON) response.ContainerDetail {
	// 先转换基本信息
	containerInfo := response.ContainerInfo{
		ID:      dockerContainer.ID,
		Name:    strings.TrimPrefix(dockerContainer.Name, "/"),
		Image:   dockerContainer.Config.Image,
		ImageID: dockerContainer.Image,
		Command: strings.Join(dockerContainer.Config.Cmd, " "),
		Created: time.Now().Unix(), // 临时使用当前时间，实际应该解析dockerContainer.Created
		Status:  dockerContainer.State.Status,
		State:   dockerContainer.State.Status,
		Labels:  dockerContainer.Config.Labels,
	}

	// 转换端口映射
	ports := make([]response.PortMapping, 0)
	if dockerContainer.NetworkSettings != nil && dockerContainer.NetworkSettings.Ports != nil {
		for portKey, bindings := range dockerContainer.NetworkSettings.Ports {
			parts := strings.Split(string(portKey), "/")
			if len(parts) == 2 {
				privatePort := 0
				if port, err := strconv.Atoi(parts[0]); err == nil {
					privatePort = port
				}
				
				for _, binding := range bindings {
					publicPort := 0
					if binding.HostPort != "" {
						if port, err := strconv.Atoi(binding.HostPort); err == nil {
							publicPort = port
						}
					}
					
					ports = append(ports, response.PortMapping{
						PrivatePort: privatePort,
						PublicPort:  publicPort,
						Type:        parts[1],
						IP:          binding.HostIP,
					})
				}
			}
		}
	}
	containerInfo.Ports = ports

	// 转换配置信息
	config := response.ContainerConfig{
		Hostname:     dockerContainer.Config.Hostname,
		Domainname:   dockerContainer.Config.Domainname,
		User:         dockerContainer.Config.User,
		Env:          dockerContainer.Config.Env,
		Cmd:          dockerContainer.Config.Cmd,
		Image:        dockerContainer.Config.Image,
		WorkingDir:   dockerContainer.Config.WorkingDir,
		Entrypoint:   dockerContainer.Config.Entrypoint,
		Labels:       dockerContainer.Config.Labels,
		ExposedPorts: convertExposedPorts(dockerContainer.Config.ExposedPorts),
	}

	// 转换主机配置
	hostConfig := convertHostConfig(dockerContainer.HostConfig)

	// 转换网络设置
	networkSettings := convertNetworkSettings(dockerContainer.NetworkSettings)

	// 转换挂载点
	mounts := convertMounts(dockerContainer.Mounts)

	return response.ContainerDetail{
		ContainerInfo:   containerInfo,
		Config:          config,
		HostConfig:      hostConfig,
		NetworkSettings: networkSettings,
		Mounts:          mounts,
	}
}

// convertHostConfig 转换主机配置
func convertHostConfig(hostConfig *container.HostConfig) response.HostConfig {
	if hostConfig == nil {
		return response.HostConfig{}
	}

	// 转换端口绑定
	portBindings := make(map[string][]response.PortBinding)
	for port, bindings := range hostConfig.PortBindings {
		responseBindings := make([]response.PortBinding, 0, len(bindings))
		for _, binding := range bindings {
			responseBindings = append(responseBindings, response.PortBinding{
				HostIP:   binding.HostIP,
				HostPort: binding.HostPort,
			})
		}
		portBindings[string(port)] = responseBindings
	}

	// 转换重启策略
	restartPolicy := response.RestartPolicy{
		Name:              hostConfig.RestartPolicy.Name,
		MaximumRetryCount: hostConfig.RestartPolicy.MaximumRetryCount,
	}

	// 转换日志配置
	logConfig := response.LogConfig{
		Type:   hostConfig.LogConfig.Type,
		Config: hostConfig.LogConfig.Config,
	}

	return response.HostConfig{
		Binds:           hostConfig.Binds,
		ContainerIDFile: hostConfig.ContainerIDFile,
		LogConfig:       logConfig,
		NetworkMode:     string(hostConfig.NetworkMode),
		PortBindings:    portBindings,
		RestartPolicy:   restartPolicy,
		AutoRemove:      hostConfig.AutoRemove,
		VolumeDriver:    hostConfig.VolumeDriver,
		VolumesFrom:     hostConfig.VolumesFrom,
		CapAdd:          hostConfig.CapAdd,
		CapDrop:         hostConfig.CapDrop,
		DNS:             hostConfig.DNS,
		DNSOptions:      hostConfig.DNSOptions,
		DNSSearch:       hostConfig.DNSSearch,
		ExtraHosts:      hostConfig.ExtraHosts,
		GroupAdd:        hostConfig.GroupAdd,
		IpcMode:         string(hostConfig.IpcMode),
		Links:           hostConfig.Links,
		OomScoreAdj:     hostConfig.OomScoreAdj,
		PidMode:         string(hostConfig.PidMode),
		Privileged:      hostConfig.Privileged,
		PublishAllPorts: hostConfig.PublishAllPorts,
		ReadonlyRootfs:  hostConfig.ReadonlyRootfs,
		SecurityOpt:     hostConfig.SecurityOpt,
		ShmSize:         hostConfig.ShmSize,
		UTSMode:         string(hostConfig.UTSMode),
		UsernsMode:      string(hostConfig.UsernsMode),
		Sysctls:         hostConfig.Sysctls,
		Runtime:         hostConfig.Runtime,
	}
}

// convertNetworkSettings 转换网络设置
func convertNetworkSettings(networkSettings *types.NetworkSettings) response.NetworkSettings {
	if networkSettings == nil {
		return response.NetworkSettings{}
	}

	// 转换端口映射
	ports := make(map[string][]response.PortBinding)
	for port, bindings := range networkSettings.Ports {
		responseBindings := make([]response.PortBinding, 0, len(bindings))
		for _, binding := range bindings {
			responseBindings = append(responseBindings, response.PortBinding{
				HostIP:   binding.HostIP,
				HostPort: binding.HostPort,
			})
		}
		ports[string(port)] = responseBindings
	}

	// 转换网络端点
	networks := make(map[string]response.EndpointSettings)
	for name, endpoint := range networkSettings.Networks {
		var ipamConfig *response.EndpointIPAMConfig
		if endpoint.IPAMConfig != nil {
			ipamConfig = &response.EndpointIPAMConfig{
				IPv4Address: endpoint.IPAMConfig.IPv4Address,
				IPv6Address: endpoint.IPAMConfig.IPv6Address,
			}
		}

		networks[name] = response.EndpointSettings{
			IPAMConfig:          ipamConfig,
			Links:               endpoint.Links,
			Aliases:             endpoint.Aliases,
			NetworkID:           endpoint.NetworkID,
			EndpointID:          endpoint.EndpointID,
			Gateway:             endpoint.Gateway,
			IPAddress:           endpoint.IPAddress,
			IPPrefixLen:         endpoint.IPPrefixLen,
			IPv6Gateway:         endpoint.IPv6Gateway,
			GlobalIPv6Address:   endpoint.GlobalIPv6Address,
			GlobalIPv6PrefixLen: endpoint.GlobalIPv6PrefixLen,
			MacAddress:          endpoint.MacAddress,
			DriverOpts:          endpoint.DriverOpts,
		}
	}

	return response.NetworkSettings{
		Bridge:                 networkSettings.Bridge,
		SandboxID:              networkSettings.SandboxID,
		HairpinMode:            networkSettings.HairpinMode,
		LinkLocalIPv6Address:   networkSettings.LinkLocalIPv6Address,
		LinkLocalIPv6PrefixLen: networkSettings.LinkLocalIPv6PrefixLen,
		Ports:                  ports,
		SandboxKey:             networkSettings.SandboxKey,
		EndpointID:             networkSettings.EndpointID,
		Gateway:                networkSettings.Gateway,
		GlobalIPv6Address:      networkSettings.GlobalIPv6Address,
		GlobalIPv6PrefixLen:    networkSettings.GlobalIPv6PrefixLen,
		IPAddress:              networkSettings.IPAddress,
		IPPrefixLen:            networkSettings.IPPrefixLen,
		IPv6Gateway:            networkSettings.IPv6Gateway,
		MacAddress:             networkSettings.MacAddress,
		Networks:               networks,
	}
}

// convertMounts 转换挂载点
func convertMounts(mounts []types.MountPoint) []response.Mount {
	responseMounts := make([]response.Mount, 0, len(mounts))
	
	for _, mount := range mounts {
		responseMount := response.Mount{
			Target:      mount.Destination,
			Source:      mount.Source,
			Type:        string(mount.Type),
			ReadOnly:    !mount.RW,
			// Consistency字段在某些Docker版本中可能不存在，暂时注释掉
			// Consistency: mount.Consistency,
		}

		// 根据挂载类型设置相应的选项
		switch mount.Type {
		case "bind":
			if mount.Propagation != "" {
				responseMount.BindOptions = &response.BindOptions{
					Propagation: string(mount.Propagation),
				}
			}
		case "volume":
			if mount.Driver != "" {
				responseMount.VolumeOptions = &response.VolumeOptions{
					DriverConfig: &response.VolumeDriverConfig{
						Name: mount.Driver,
					},
				}
			}
		}

		responseMounts = append(responseMounts, responseMount)
	}

	return responseMounts
}
// convertExposedPorts 转换暴露端口
func convertExposedPorts(exposedPorts nat.PortSet) map[string]struct{} {
	result := make(map[string]struct{})
	for port := range exposedPorts {
		result[string(port)] = struct{}{}
	}
	return result
}