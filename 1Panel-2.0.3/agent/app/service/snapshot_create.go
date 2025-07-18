package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/1Panel-dev/1Panel/agent/app/repo"

	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/1Panel-dev/1Panel/agent/app/task"
	"github.com/1Panel-dev/1Panel/agent/constant"
	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/1Panel-dev/1Panel/agent/i18n"
	"github.com/1Panel-dev/1Panel/agent/utils/cmd"
	"github.com/1Panel-dev/1Panel/agent/utils/common"
	"github.com/1Panel-dev/1Panel/agent/utils/copier"
	"github.com/1Panel-dev/1Panel/agent/utils/files"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (u *SnapshotService) SnapshotCreate(parentTask *task.Task, req dto.SnapshotCreate, jobID, retry, timeout uint) error {
	versionItem, _ := settingRepo.Get(settingRepo.WithByKey("SystemVersion"))

	scope := "core"
	if !global.IsMaster {
		scope = "agent"
	}
	if jobID == 0 {
		req.Name = fmt.Sprintf("1panel-%s-%s-linux-%s-%s", scope, versionItem.Value, loadOs(), time.Now().Format(constant.DateTimeSlimLayout))
	}
	appItem, _ := json.Marshal(req.AppData)
	panelItem, _ := json.Marshal(req.PanelData)
	backupItem, _ := json.Marshal(req.BackupData)
	snap := model.Snapshot{
		Name:              req.Name,
		TaskID:            req.TaskID,
		Secret:            req.Secret,
		Description:       req.Description,
		SourceAccountIDs:  req.SourceAccountIDs,
		DownloadAccountID: req.DownloadAccountID,

		AppData:          string(appItem),
		PanelData:        string(panelItem),
		BackupData:       string(backupItem),
		WithDockerConf:   req.WithDockerConf,
		WithMonitorData:  req.WithMonitorData,
		WithLoginLog:     req.WithLoginLog,
		WithOperationLog: req.WithOperationLog,
		WithTaskLog:      req.WithTaskLog,
		WithSystemLog:    req.WithSystemLog,
		IgnoreFiles:      strings.Join(req.IgnoreFiles, ","),

		Version: versionItem.Value,
		Status:  constant.StatusWaiting,
	}
	if err := snapshotRepo.Create(&snap); err != nil {
		global.LOG.Errorf("create snapshot record to db failed, err: %v", err)
		return err
	}

	req.ID = snap.ID
	var err error
	taskItem := parentTask
	if parentTask == nil {
		taskItem, err = task.NewTaskWithOps(req.Name, task.TaskCreate, task.TaskScopeSnapshot, req.TaskID, req.ID)
		if err != nil {
			global.LOG.Errorf("new task for create snapshot failed, err: %v", err)
			return err
		}
	}
	if jobID == 0 {
		go func() {
			_ = handleSnapshot(req, taskItem, jobID, 3, 0)
		}()
		return nil
	}

	return handleSnapshot(req, taskItem, jobID, retry, timeout)
}

func (u *SnapshotService) SnapshotReCreate(id uint) error {
	snap, err := snapshotRepo.Get(repo.WithByID(id))
	if err != nil {
		return err
	}
	taskModel, err := taskRepo.GetFirst(taskRepo.WithResourceID(snap.ID), repo.WithByType(task.TaskScopeSnapshot))
	if err != nil {
		return err
	}

	var req dto.SnapshotCreate
	_ = copier.Copy(&req, snap)
	if err := json.Unmarshal([]byte(snap.PanelData), &req.PanelData); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(snap.AppData), &req.AppData); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(snap.BackupData), &req.BackupData); err != nil {
		return err
	}
	req.TaskID = taskModel.ID
	taskItem, err := task.ReNewTaskWithOps(req.Name, task.TaskCreate, task.TaskScopeSnapshot, req.TaskID, req.ID)
	if err != nil {
		global.LOG.Errorf("new task for create snapshot failed, err: %v", err)
		return err
	}
	_ = snapshotRepo.Update(req.ID, map[string]interface{}{"status": constant.StatusWaiting, "message": ""})
	go func() {
		_ = handleSnapshot(req, taskItem, 0, 3, 0)
	}()

	return nil
}

func handleSnapshot(req dto.SnapshotCreate, taskItem *task.Task, jobID, retry, timeout uint) error {
	rootDir := path.Join(global.Dir.TmpDir, "system", req.Name)
	openrestyDir, _ := settingRepo.GetValueByKey("WEBSITE_DIR")
	itemHelper := snapHelper{SnapID: req.ID, Task: *taskItem, FileOp: files.NewFileOp(), Ctx: context.Background(), OpenrestyDir: openrestyDir}
	baseDir := path.Join(rootDir, "base")
	_ = os.MkdirAll(baseDir, os.ModePerm)

	if timeout == 0 {
		timeout = 1800
	}
	taskItem.AddSubTaskWithAliasAndOps(
		"SnapDBInfo",
		func(t *task.Task) error {
			if err := loadDbConn(&itemHelper, rootDir, req); err != nil {
				return err
			}
			_ = itemHelper.snapAgentDB.Where("id = ?", req.ID).Delete(&model.Snapshot{}).Error
			if jobID != 0 {
				_ = itemHelper.snapAgentDB.Where("id = ?", jobID).Delete(&model.JobRecords{}).Error
			}
			return nil
		}, nil, int(retry), time.Duration(timeout)*time.Second,
	)

	if len(req.InterruptStep) == 0 || req.InterruptStep == "SnapBaseInfo" {
		taskItem.AddSubTaskWithAliasAndOps(
			"SnapBaseInfo",
			func(t *task.Task) error { return snapBaseData(itemHelper, baseDir, req.WithDockerConf) },
			nil, int(retry), time.Duration(timeout)*time.Second,
		)
		req.InterruptStep = ""
	}
	if len(req.InterruptStep) == 0 || req.InterruptStep == "SnapInstallApp" {
		taskItem.AddSubTaskWithAliasAndOps(
			"SnapInstallApp",
			func(t *task.Task) error { return snapAppImage(itemHelper, req, rootDir) },
			nil, int(retry), time.Duration(timeout)*time.Second,
		)
		req.InterruptStep = ""
	}
	if len(req.InterruptStep) == 0 || req.InterruptStep == "SnapLocalBackup" {
		taskItem.AddSubTaskWithAliasAndOps(
			"SnapLocalBackup",
			func(t *task.Task) error { return snapBackupData(itemHelper, req, rootDir) },
			nil, int(retry), time.Duration(timeout)*time.Second,
		)
		req.InterruptStep = ""
	}
	if len(req.InterruptStep) == 0 || req.InterruptStep == "SnapPanelData" {
		taskItem.AddSubTaskWithAliasAndOps(
			"SnapPanelData",
			func(t *task.Task) error { return snapPanelData(itemHelper, req, rootDir) },
			nil, int(retry), time.Duration(timeout)*time.Second,
		)
		req.InterruptStep = ""
	}

	taskItem.AddSubTaskWithAliasAndOps(
		"SnapCloseDBConn",
		func(t *task.Task) error {
			taskItem.Log("---------------------- 6 / 8 ----------------------")
			common.CloseDB(itemHelper.snapAgentDB)
			common.CloseDB(itemHelper.snapCoreDB)
			return nil
		}, nil, int(retry), time.Duration(timeout)*time.Second,
	)
	if len(req.InterruptStep) == 0 || req.InterruptStep == "SnapCompress" {
		taskItem.AddSubTaskWithAliasAndOps(
			"SnapCompress",
			func(t *task.Task) error { return snapCompress(itemHelper, rootDir, req.Secret) },
			nil, int(retry), time.Duration(timeout)*time.Second,
		)
		req.InterruptStep = ""
	}
	if len(req.InterruptStep) == 0 || req.InterruptStep == "SnapUpload" {
		taskItem.AddSubTaskWithAliasAndOps(
			"SnapUpload",
			func(t *task.Task) error {
				return snapUpload(itemHelper, req.SourceAccountIDs, fmt.Sprintf("%s.tar.gz", rootDir))
			}, nil, int(retry), time.Duration(timeout)*time.Second,
		)
		req.InterruptStep = ""
	}
	if err := taskItem.Execute(); err != nil {
		_ = snapshotRepo.Update(req.ID, map[string]interface{}{"status": constant.StatusFailed, "message": err.Error(), "interrupt_step": taskItem.Task.CurrentStep})
		return err
	}
	_ = snapshotRepo.Update(req.ID, map[string]interface{}{"status": constant.StatusSuccess, "interrupt_step": ""})
	_ = os.RemoveAll(rootDir)
	return nil
}

type snapHelper struct {
	SnapID      uint
	SnapName    string
	snapAgentDB *gorm.DB
	snapCoreDB  *gorm.DB
	Ctx         context.Context
	FileOp      files.FileOp
	Wg          *sync.WaitGroup
	Task        task.Task

	OpenrestyDir string
}

func loadDbConn(snap *snapHelper, targetDir string, req dto.SnapshotCreate) error {
	snap.Task.Log("---------------------- 1 / 8 ----------------------")
	snap.Task.LogStart(i18n.GetMsgByKey("SnapDBInfo"))
	pathDB := path.Join(global.Dir.DataDir, "db")

	err := snap.FileOp.CopyDir(pathDB, targetDir)
	snap.Task.LogWithStatus(i18n.GetWithName("SnapCopy", pathDB), err)
	if err != nil {
		return err
	}
	_ = os.Remove(path.Join(targetDir, "db/session.db"))

	agentDb, err := common.LoadDBConnByPathWithErr(path.Join(targetDir, "db/agent.db"), "agent.db")
	snap.Task.LogWithStatus(i18n.GetWithName("SnapNewDB", "agent"), err)
	if err != nil {
		return err
	}
	snap.snapAgentDB = agentDb
	coreDb, err := common.LoadDBConnByPathWithErr(path.Join(targetDir, "db/core.db"), "core.db")
	snap.Task.LogWithStatus(i18n.GetWithName("SnapNewDB", "core"), err)
	if err != nil {
		return err
	}
	snap.snapCoreDB = coreDb

	if !req.WithMonitorData {
		err = os.Remove(path.Join(targetDir, "db/monitor.db"))
		snap.Task.LogWithStatus(i18n.GetMsgByKey("SnapDeleteMonitor"), err)
		if err != nil {
			return err
		}
	}
	if !req.WithTaskLog {
		err = os.Remove(path.Join(targetDir, "db/task.db"))
		if err != nil {
			return err
		}
	} else {
		taskDB, err := common.LoadDBConnByPathWithErr(path.Join(targetDir, "db/task.db"), "core.db")
		snap.Task.LogWithStatus(i18n.GetWithName("SnapNewDB", "task"), err)
		if err != nil {
			return err
		}
		_ = taskDB.Where("id = ?", req.TaskID).Delete(&model.Task{}).Error
	}
	if !req.WithOperationLog && global.IsMaster {
		err = snap.snapCoreDB.Exec("DELETE FROM operation_logs").Error
		snap.Task.LogWithStatus(i18n.GetMsgByKey("SnapDeleteOperationLog"), err)
		if err != nil {
			return err
		}
	}
	if !req.WithLoginLog && global.IsMaster {
		err = snap.snapCoreDB.Exec("DELETE FROM login_logs").Error
		snap.Task.LogWithStatus(i18n.GetMsgByKey("SnapDeleteLoginLog"), err)
		if err != nil {
			return err
		}
	}
	return nil
}

func snapBaseData(snap snapHelper, targetDir string, withDockerConf bool) error {
	snap.Task.Log("---------------------- 2 / 8 ----------------------")
	snap.Task.LogStart(i18n.GetMsgByKey("SnapBaseInfo"))

	if global.IsMaster {
		err := snap.FileOp.CopyFile("/usr/local/bin/1panel-core", targetDir)
		snap.Task.LogWithStatus(i18n.GetWithName("SnapCopy", "/usr/local/bin/1panel-core"), err)
		if err != nil {
			return err
		}
		err = snap.FileOp.CopyFile("/usr/local/bin/1pctl", targetDir)
		snap.Task.LogWithStatus(i18n.GetWithName("SnapCopy", "/usr/local/bin/1pctl"), err)
		if err != nil {
			return err
		}
	}

	err := snap.FileOp.CopyFile("/usr/local/bin/1panel-agent", targetDir)
	snap.Task.LogWithStatus(i18n.GetWithName("SnapCopy", "/usr/local/bin/1panel-agent"), err)
	if err != nil {
		return err
	}

	if global.IsMaster {
		err = snap.FileOp.CopyFile("/etc/systemd/system/1panel-core.service", targetDir)
		snap.Task.LogWithStatus(i18n.GetWithName("SnapCopy", "/etc/systemd/system/1panel-core.service"), err)
		if err != nil {
			return err
		}
	}

	err = snap.FileOp.CopyFile("/etc/systemd/system/1panel-agent.service", targetDir)
	snap.Task.LogWithStatus(i18n.GetWithName("SnapCopy", "/etc/systemd/system/1panel-agent.service"), err)
	if err != nil {
		return err
	}

	if withDockerConf {
		if snap.FileOp.Stat(constant.DaemonJsonPath) {
			err = snap.FileOp.CopyFile(constant.DaemonJsonPath, targetDir)
			snap.Task.LogWithStatus(i18n.GetWithName("SnapCopy", constant.DaemonJsonPath), err)
			if err != nil {
				return err
			}
		}
	}

	remarkInfo, _ := json.MarshalIndent(SnapshotJson{
		BaseDir:       global.Dir.BaseDir,
		OperestyDir:   snap.OpenrestyDir,
		BackupDataDir: global.Dir.LocalBackupDir,
	}, "", "\t")
	err = os.WriteFile(path.Join(targetDir, "snapshot.json"), remarkInfo, 0640)
	snap.Task.LogWithStatus(i18n.GetWithName("SnapCopy", path.Join(targetDir, "snapshot.json")), err)
	if err != nil {
		return err
	}

	return nil
}

func snapAppImage(snap snapHelper, req dto.SnapshotCreate, targetDir string) error {
	snap.Task.Log("---------------------- 3 / 8 ----------------------")
	snap.Task.LogStart(i18n.GetMsgByKey("SnapInstallApp"))

	var appInstalls []model.AppInstall
	_ = snap.snapAgentDB.Where("1 = 1").Find(&appInstalls).Error
	for _, item := range appInstalls {
		if err := snap.snapAgentDB.
			Model(&model.AppInstall{}).
			Where("id = ?", item.ID).
			Updates(map[string]interface{}{"status": constant.StatusWaitingRestart}).
			Error; err != nil {
			global.LOG.Errorf("update app %s status failed, err: %v", item.Name, err)
		}
	}

	var imageList []string
	existStr, _ := cmd.RunDefaultWithStdoutBashC("docker images | awk '{print $1\":\"$2}' | grep -v REPOSITORY:TAG")
	existImages := strings.Split(existStr, "\n")
	for _, app := range req.AppData {
		for _, item := range app.Children {
			if item.Label == "appImage" && item.IsCheck {
				for _, existImage := range existImages {
					if len(existImage) == 0 {
						continue
					}
					if existImage == item.Name {
						imageList = append(imageList, item.Name)
					}
				}
			}
		}
	}

	if len(imageList) != 0 {
		snap.Task.Log(strings.Join(imageList, " "))
		snap.Task.Logf("docker save %s | gzip -c > %s", strings.Join(imageList, " "), path.Join(targetDir, "images.tar.gz"))
		std, err := cmd.NewCommandMgr(cmd.WithTimeout(10*time.Minute)).RunWithStdoutBashCf("docker save %s | gzip -c > %s", strings.Join(imageList, " "), path.Join(targetDir, "images.tar.gz"))
		if err != nil {
			snap.Task.LogFailedWithErr(i18n.GetMsgByKey("SnapDockerSave"), errors.New(std))
			return errors.New(std)
		}
		snap.Task.LogSuccess(i18n.GetMsgByKey("SnapDockerSave"))
	} else {
		snap.Task.Log(i18n.GetMsgByKey("SnapInstallAppImageEmpty"))
	}
	return nil
}

func snapBackupData(snap snapHelper, req dto.SnapshotCreate, targetDir string) error {
	snap.Task.Log("---------------------- 4 / 8 ----------------------")
	snap.Task.LogStart(i18n.GetMsgByKey("SnapLocalBackup"))

	excludes := loadBackupExcludes(snap, req.BackupData)
	excludes = append(excludes, req.IgnoreFiles...)
	excludes = append(excludes, "./system_snapshot")
	for _, item := range req.AppData {
		for _, itemApp := range item.Children {
			if itemApp.Label == "appBackup" {
				excludes = append(excludes, loadAppBackupExcludes([]dto.DataTree{itemApp})...)
			}
		}
	}
	err := snap.FileOp.TarGzCompressPro(false, global.Dir.LocalBackupDir, path.Join(targetDir, "1panel_backup.tar.gz"), "", strings.Join(excludes, ","))
	snap.Task.LogWithStatus(i18n.GetMsgByKey("SnapCompressBackup"), err)

	return err
}
func loadBackupExcludes(snap snapHelper, req []dto.DataTree) []string {
	var excludes []string
	for _, item := range req {
		if len(item.Children) == 0 {
			if item.IsCheck {
				continue
			}
			if err := snap.snapAgentDB.Where("file_dir = ? AND file_name = ?", strings.TrimPrefix(path.Dir(item.Path), global.Dir.LocalBackupDir+"/"), path.Base(item.Path)).Delete(&model.BackupRecord{}).Error; err != nil {
				snap.Task.LogWithStatus("delete backup file from database", err)
			}
			itemDir, _ := filepath.Rel(item.Path, global.Dir.LocalBackupDir)
			excludes = append(excludes, itemDir)
		} else {
			excludes = append(excludes, loadBackupExcludes(snap, item.Children)...)
		}
	}
	return excludes
}
func loadAppBackupExcludes(req []dto.DataTree) []string {
	var excludes []string
	for _, item := range req {
		if len(item.Children) == 0 {
			if !item.IsCheck {
				itemDir, _ := filepath.Rel(item.Path, global.Dir.LocalBackupDir)
				excludes = append(excludes, itemDir)
			}
		} else {
			excludes = append(excludes, loadAppBackupExcludes(item.Children)...)
		}
	}
	return excludes
}

func snapPanelData(snap snapHelper, req dto.SnapshotCreate, targetDir string) error {
	snap.Task.Log("---------------------- 5 / 8 ----------------------")
	snap.Task.LogStart(i18n.GetMsgByKey("SnapPanelData"))

	excludes := loadPanelExcludes(req.PanelData)
	for _, item := range req.AppData {
		for _, itemApp := range item.Children {
			if itemApp.Label == "appData" {
				excludes = append(excludes, loadPanelExcludes([]dto.DataTree{itemApp})...)
			}
		}
	}
	excludes = append(excludes, "./cache")
	excludes = append(excludes, "./db")
	excludes = append(excludes, "./tmp")
	if !req.WithSystemLog {
		excludes = append(excludes, "./log/1Panel*")
	}
	if !req.WithTaskLog {
		excludes = append(excludes, "./log/task")
	}

	rootDir := global.Dir.DataDir
	if strings.Contains(global.Dir.LocalBackupDir, rootDir) {
		itemDir, _ := filepath.Rel(rootDir, global.Dir.LocalBackupDir)
		excludes = append(excludes, itemDir)
	}
	if len(snap.OpenrestyDir) != 0 && strings.Contains(snap.OpenrestyDir, rootDir) {
		itemDir, _ := filepath.Rel(rootDir, snap.OpenrestyDir)
		excludes = append(excludes, itemDir)
	}
	excludes = append(excludes, req.IgnoreFiles...)
	err := snap.FileOp.TarGzCompressPro(false, rootDir, path.Join(targetDir, "1panel_data.tar.gz"), "", strings.Join(excludes, ","))
	snap.Task.LogWithStatus(i18n.GetMsgByKey("SnapCompressPanel"), err)
	if err != nil {
		return err
	}
	if len(snap.OpenrestyDir) != 0 {
		err := snap.FileOp.TarGzCompressPro(false, snap.OpenrestyDir, path.Join(targetDir, "website.tar.gz"), "", strings.Join(req.IgnoreFiles, ","))
		snap.Task.LogWithStatus(i18n.GetMsgByKey("SnapWebsite"), err)
		if err != nil {
			return err
		}
	}

	return err
}
func loadPanelExcludes(req []dto.DataTree) []string {
	var excludes []string
	for _, item := range req {
		if len(item.Children) == 0 {
			if !item.IsCheck {
				itemDir, _ := filepath.Rel(item.Path, path.Join(global.Dir.BaseDir, "1panel"))
				excludes = append(excludes, itemDir)
			}
		} else {
			excludes = append(excludes, loadPanelExcludes(item.Children)...)
		}
	}
	return excludes
}

func snapCompress(snap snapHelper, rootDir string, secret string) error {
	snap.Task.Log("---------------------- 7 / 8 ----------------------")
	snap.Task.LogStart(i18n.GetMsgByKey("SnapCompress"))

	tmpDir := path.Join(global.Dir.TmpDir, "system")
	fileName := fmt.Sprintf("%s.tar.gz", path.Base(rootDir))
	err := snap.FileOp.TarGzCompressPro(true, rootDir, path.Join(tmpDir, fileName), secret, "")
	snap.Task.LogWithStatus(i18n.GetMsgByKey("SnapCompressFile"), err)
	if err != nil {
		return err
	}

	stat, err := os.Stat(path.Join(tmpDir, fileName))
	snap.Task.LogWithStatus(i18n.GetMsgByKey("SnapCheckCompress"), err)
	if err != nil {
		return err
	}

	size := common.LoadSizeUnit2F(float64(stat.Size()))
	snap.Task.Log(i18n.GetWithName("SnapCompressSize", size))
	_ = os.RemoveAll(rootDir)
	return nil
}

func snapUpload(snap snapHelper, accounts string, file string) error {
	snap.Task.Log("---------------------- 8 / 8 ----------------------")
	snap.Task.LogStart(i18n.GetMsgByKey("SnapUpload"))

	source := path.Join(global.Dir.TmpDir, "system", path.Base(file))
	accountMap, err := NewBackupClientMap(strings.Split(accounts, ","))
	snap.Task.LogWithStatus(i18n.GetMsgByKey("SnapLoadBackup"), err)
	if err != nil {
		return err
	}

	targetAccounts := strings.Split(accounts, ",")
	for _, item := range targetAccounts {
		snap.Task.LogStart(i18n.GetWithName("SnapUploadTo", fmt.Sprintf("[%s] %s", accountMap[item].name, path.Join("system_snapshot", path.Base(file)))))
		_, err := accountMap[item].client.Upload(source, path.Join(accountMap[item].backupPath, "system_snapshot", path.Base(file)))
		snap.Task.LogWithStatus(i18n.GetWithName("SnapUploadRes", accountMap[item].name), err)
		if err != nil {
			return err
		}
	}
	_ = os.Remove(source)
	return nil
}
