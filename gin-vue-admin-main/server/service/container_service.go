package service

import (
	"os/exec"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model"
)

type ContainerInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
	Ports  string `json:"ports"`
}

func GetDockerContainers(page, pageSize int) ([]ContainerInfo, int, error) {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.ID}}|{{.Image}}|{{.Names}}|{{.Status}}|{{.Ports}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, 0, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	total := len(lines)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		return []ContainerInfo{}, total, nil
	}
	if end > total {
		end = total
	}
	var result []ContainerInfo
	for _, line := range lines[start:end] {
		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			continue
		}
		result = append(result, ContainerInfo{
			ID:     parts[0],
			Image:  parts[1],
			Name:   parts[2],
			Status: parts[3],
			Ports:  parts[4],
		})
	}
	return result, total, nil
}

func CreateContainer(container *model.Container) error {
	return global.GVA_DB.Create(container).Error
}

func DeleteContainer(id uint) error {
	return global.GVA_DB.Delete(&model.Container{}, id).Error
}
