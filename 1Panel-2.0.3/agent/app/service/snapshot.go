package service

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/1Panel-dev/1Panel/agent/app/repo"
	"github.com/1Panel-dev/1Panel/agent/app/task"
	"github.com/1Panel-dev/1Panel/agent/buserr"

	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/1Panel-dev/1Panel/agent/constant"
	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/1Panel-dev/1Panel/agent/utils/cmd"
	"github.com/1Panel-dev/1Panel/agent/utils/docker"
	fileUtils "github.com/1Panel-dev/1Panel/agent/utils/files"
	"github.com/docker/docker/api/types/image"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/shirou/gopsutil/v3/host"
)

type SnapshotService struct {
	OriginalPath string
}

type ISnapshotService interface {
	SearchWithPage(req dto.PageSnapshot) (int64, interface{}, error)
	LoadSnapshotData() (dto.SnapshotData, error)
	SnapshotCreate(parentTask *task.Task, req dto.SnapshotCreate, jobID, retry, timeout uint) error
	SnapshotReCreate(id uint) error
	SnapshotRecover(req dto.SnapshotRecover) error
	SnapshotRollback(req dto.SnapshotRecover) error
	SnapshotImport(req dto.SnapshotImport) error
	Delete(req dto.SnapshotBatchDelete) error

	UpdateDescription(req dto.UpdateDescription) error
}

func NewISnapshotService() ISnapshotService {
	return &SnapshotService{}
}

func (u *SnapshotService) SearchWithPage(req dto.PageSnapshot) (int64, interface{}, error) {
	total, records, err := snapshotRepo.Page(req.Page, req.PageSize, repo.WithByLikeName(req.Info), repo.WithOrderRuleBy(req.OrderBy, req.Order))
	if err != nil {
		return 0, nil, err
	}
	var datas []dto.SnapshotInfo
	for i := 0; i < len(records); i++ {
		var item dto.SnapshotInfo
		if err := copier.Copy(&item, &records[i]); err != nil {
			return 0, nil, err
		}
		item.SourceAccounts, item.DownloadAccount, _ = loadBackupNamesByID(records[i].SourceAccountIDs, records[i].DownloadAccountID)
		datas = append(datas, item)
	}
	return total, datas, err
}

func (u *SnapshotService) SnapshotImport(req dto.SnapshotImport) error {
	if len(req.Names) == 0 {
		return fmt.Errorf("incorrect snapshot request body: %v", req.Names)
	}

	for _, snapName := range req.Names {
		if !checkSnapshotIsOk(snapName) {
			return fmt.Errorf("incorrect snapshot name format of %s", snapName)
		}
		snap, _ := snapshotRepo.Get(repo.WithByName(strings.ReplaceAll(snapName, ".tar.gz", "")))
		if snap.ID != 0 {
			return buserr.New("ErrRecordExist")
		}
	}
	for _, snap := range req.Names {
		shortName := strings.ReplaceAll(snap, "snapshot-", "")
		shortName = strings.ReplaceAll(shortName, "1panel-", "")
		shortName = strings.ReplaceAll(shortName, "core-", "")
		shortName = strings.ReplaceAll(shortName, "agent-", "")
		nameItems := strings.Split(shortName, "-linux")
		if strings.HasSuffix(snap, ".tar.gz") {
			snap = strings.ReplaceAll(snap, ".tar.gz", "")
		}
		itemSnap := model.Snapshot{
			Name:              snap,
			SourceAccountIDs:  fmt.Sprintf("%v", req.BackupAccountID),
			DownloadAccountID: req.BackupAccountID,
			Version:           nameItems[0],
			Description:       req.Description,
			Status:            constant.StatusSuccess,
		}
		if err := snapshotRepo.Create(&itemSnap); err != nil {
			return err
		}
	}
	return nil
}

func (u *SnapshotService) LoadSnapshotData() (dto.SnapshotData, error) {
	var (
		data dto.SnapshotData
		err  error
	)
	fileOp := fileUtils.NewFileOp()
	data.AppData, err = loadApps(fileOp)
	if err != nil {
		return data, err
	}
	data.WithDockerConf = true
	data.PanelData, err = loadPanelFile(fileOp)
	if err != nil {
		return data, err
	}
	itemBackups, err := loadFile(global.Dir.LocalBackupDir, 8, fileOp)
	if err != nil {
		return data, err
	}
	i := 0
	for _, item := range itemBackups {
		if item.Label != "app" && item.Label != "system_snapshot" {
			itemBackups[i] = item
			i++
		}
	}
	data.BackupData = itemBackups[:i]

	return data, nil
}

func (u *SnapshotService) UpdateDescription(req dto.UpdateDescription) error {
	return snapshotRepo.Update(req.ID, map[string]interface{}{"description": req.Description})
}

type SnapshotJson struct {
	BaseDir       string `json:"baseDir"`
	OperestyDir   string `json:"operestyDir"`
	BackupDataDir string `json:"backupDataDir"`
	Size          uint64 `json:"size"`
}

func (u *SnapshotService) Delete(req dto.SnapshotBatchDelete) error {
	snaps, _ := snapshotRepo.GetList(repo.WithByIDs(req.Ids))
	for _, snap := range snaps {
		if req.DeleteWithFile {
			accounts, err := NewBackupClientMap(strings.Split(snap.SourceAccountIDs, ","))
			if err != nil {
				return err
			}
			for _, item := range accounts {
				global.LOG.Debugf("remove snapshot file %s.tar.gz from %s", snap.Name, item.name)
				_, _ = item.client.Delete(path.Join(item.backupPath, "system_snapshot", snap.Name+".tar.gz"))
			}
			_ = backupRepo.DeleteRecord(context.Background(), repo.WithByType("snapshot"), backupRepo.WithByFileName(snap.Name+".tar.gz"))
		}

		if err := snapshotRepo.Delete(repo.WithByID(snap.ID)); err != nil {
			return err
		}
	}
	return nil
}

func hasOs(name string) bool {
	return strings.Contains(name, "amd64") ||
		strings.Contains(name, "arm64") ||
		strings.Contains(name, "armv7") ||
		strings.Contains(name, "ppc64le") ||
		strings.Contains(name, "s390x")
}

func loadOs() string {
	hostInfo, _ := host.Info()
	switch hostInfo.KernelArch {
	case "x86_64":
		return "amd64"
	case "armv7l":
		return "armv7"
	default:
		return hostInfo.KernelArch
	}
}

func loadApps(fileOp fileUtils.FileOp) ([]dto.DataTree, error) {
	var data []dto.DataTree
	apps, err := appInstallRepo.ListBy(context.Background())
	if err != nil {
		return data, err
	}
	openrestyID := 0
	for _, app := range apps {
		if app.App.Key == constant.AppOpenresty {
			openrestyID = int(app.ID)
		}
	}
	websites, err := websiteRepo.List()
	if err != nil {
		return data, err
	}
	appRelationMap := make(map[uint]uint)
	for _, website := range websites {
		if website.Type == constant.Deployment && website.AppInstallID != 0 {
			appRelationMap[uint(openrestyID)] = website.AppInstallID
		}
	}
	appRelations, err := appInstallResourceRepo.GetBy()
	if err != nil {
		return data, err
	}
	for _, relation := range appRelations {
		appRelationMap[uint(relation.AppInstallId)] = relation.LinkId
	}
	appMap := make(map[uint]string)
	for _, app := range apps {
		appMap[app.ID] = fmt.Sprintf("%s-%s", app.App.Key, app.Name)
	}

	appTreeMap := make(map[string]dto.DataTree)
	for _, app := range apps {
		itemApp := dto.DataTree{
			ID:    uuid.NewString(),
			Label: fmt.Sprintf("%s - %s", app.App.Name, app.Name),
			Key:   app.App.Key,
			Name:  app.Name,
		}
		appPath := path.Join(global.Dir.DataDir, "apps", app.App.Key, app.Name)
		itemAppData := dto.DataTree{ID: uuid.NewString(), Label: "appData", Key: app.App.Key, Name: app.Name, IsCheck: true, Path: appPath}
		if app.App.Key == constant.AppOpenresty && len(websites) != 0 {
			itemAppData.IsDisable = true
		}
		if val, ok := appRelationMap[app.ID]; ok {
			itemAppData.RelationItemID = appMap[val]
		}
		sizeItem, err := fileOp.GetDirSize(appPath)
		if err == nil {
			itemAppData.Size = uint64(sizeItem)
		}
		itemApp.Size += itemAppData.Size
		data = append(data, itemApp)
		appTreeMap[fmt.Sprintf("%s-%s", itemApp.Key, itemApp.Name)] = itemAppData
	}

	for key, val := range appTreeMap {
		if valRelation, ok := appTreeMap[val.RelationItemID]; ok {
			valRelation.IsDisable = true
			appTreeMap[val.RelationItemID] = valRelation

			val.RelationItemID = valRelation.ID
			appTreeMap[key] = val
		}
	}
	for i := 0; i < len(data); i++ {
		if val, ok := appTreeMap[fmt.Sprintf("%s-%s", data[i].Key, data[i].Name)]; ok {
			data[i].Children = append(data[i].Children, val)
		}
	}
	data = loadAppBackup(data, fileOp)
	data = loadAppImage(data)
	return data, nil
}
func loadAppBackup(list []dto.DataTree, fileOp fileUtils.FileOp) []dto.DataTree {
	for i := 0; i < len(list); i++ {
		appBackupPath := path.Join(global.Dir.LocalBackupDir, "app", list[i].Key, list[i].Name)
		itemAppBackupTree, err := loadFile(appBackupPath, 8, fileOp)
		itemAppBackup := dto.DataTree{ID: uuid.NewString(), Label: "appBackup", IsCheck: true, Children: itemAppBackupTree, Path: appBackupPath}
		if err == nil {
			backupSizeItem, err := fileOp.GetDirSize(appBackupPath)
			if err == nil {
				itemAppBackup.Size = uint64(backupSizeItem)
				list[i].Size += itemAppBackup.Size
			}
			list[i].Children = append(list[i].Children, itemAppBackup)
		}
	}
	return list
}
func loadAppImage(list []dto.DataTree) []dto.DataTree {
	client, err := docker.NewDockerClient()
	if err != nil {
		global.LOG.Errorf("new docker client failed, err: %v", err)
		return list
	}
	defer client.Close()
	imageList, err := client.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		global.LOG.Errorf("load image list failed, err: %v", err)
		return list
	}

	for i := 0; i < len(list); i++ {
		itemAppImage := dto.DataTree{ID: uuid.NewString(), Label: "appImage"}
		stdout, err := cmd.RunDefaultWithStdoutBashCf("cat %s | grep image: ", path.Join(global.Dir.AppDir, list[i].Key, list[i].Name, "docker-compose.yml"))
		if err != nil {
			list[i].Children = append(list[i].Children, itemAppImage)
			continue
		}
		itemAppImage.Name = strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(stdout), "\n", ""), "image: ", "")
		for _, imageItem := range imageList {
			for _, tag := range imageItem.RepoTags {
				if tag == itemAppImage.Name {
					itemAppImage.Size = uint64(imageItem.Size)
					break
				}
			}
		}
		list[i].Children = append(list[i].Children, itemAppImage)
	}
	return list
}

func loadPanelFile(fileOp fileUtils.FileOp) ([]dto.DataTree, error) {
	var data []dto.DataTree
	snapFiles, err := os.ReadDir(global.Dir.DataDir)
	if err != nil {
		return data, err
	}
	for _, fileItem := range snapFiles {
		itemData := dto.DataTree{
			ID:      uuid.NewString(),
			Label:   fileItem.Name(),
			IsCheck: true,
			Path:    path.Join(global.Dir.DataDir, fileItem.Name()),
		}
		switch itemData.Label {
		case "agent", "runtime", "docker", "task", "geo", "secret", "uploads":
			itemData.IsDisable = true
		case "clamav", "download":
			panelPath := path.Join(global.Dir.DataDir, itemData.Label)
			itemData.Children, _ = loadFile(panelPath, 3, fileOp)
		case "apps", "backup", "log", "db", "tmp":
			continue
		}
		if fileItem.IsDir() {
			sizeItem, err := fileOp.GetDirSize(path.Join(global.Dir.DataDir, itemData.Label))
			if err != nil {
				continue
			}
			itemData.Size = uint64(sizeItem)
		} else {
			fileInfo, err := fileItem.Info()
			if err != nil {
				continue
			}
			itemData.Size = uint64(fileInfo.Size())
		}
		if itemData.IsCheck && itemData.Size == 0 {
			itemData.IsCheck = false
			itemData.IsDisable = true
		}

		data = append(data, itemData)
	}

	openrestySite, _ := settingRepo.GetValueByKey("WEBSITE_DIR")
	if len(openrestySite) != 0 && !strings.Contains(openrestySite, global.Dir.DataDir) {
		sizeItem, _ := fileOp.GetDirSize(openrestySite)
		data = append(data, dto.DataTree{
			ID:        uuid.NewString(),
			Label:     "www",
			IsCheck:   true,
			IsDisable: true,
			Path:      openrestySite,
			Size:      uint64(sizeItem),
		})
	}

	return data, nil
}

func loadFile(pathItem string, index int, fileOp fileUtils.FileOp) ([]dto.DataTree, error) {
	var data []dto.DataTree
	snapFiles, err := os.ReadDir(pathItem)
	if err != nil {
		return data, err
	}
	for _, fileItem := range snapFiles {
		itemData := dto.DataTree{
			ID:      uuid.NewString(),
			Label:   fileItem.Name(),
			Name:    fileItem.Name(),
			Path:    path.Join(pathItem, fileItem.Name()),
			IsCheck: true,
		}
		if fileItem.IsDir() {
			sizeItem, err := fileOp.GetDirSize(path.Join(pathItem, itemData.Label))
			if err != nil {
				continue
			}
			itemData.Size = uint64(sizeItem)
			if index > 1 {
				itemData.Children, _ = loadFile(path.Join(pathItem, itemData.Label), index-1, fileOp)
			}
		} else {
			fileInfo, err := fileItem.Info()
			if err != nil {
				continue
			}
			itemData.Size = uint64(fileInfo.Size())
		}
		data = append(data, itemData)
	}
	return data, nil
}

func checkSnapshotIsOk(name string) bool {
	names := []string{"1panel-core-v2.", "1panel-agent-v2.", "1panel-v2.", "snapshot-1panel-core-v2.", "snapshot-1panel-agent-v2."}
	for _, item := range names {
		if strings.HasPrefix(name, item) {
			return true
		}
	}
	return false
}
