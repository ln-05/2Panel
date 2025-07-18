package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/1Panel-dev/1Panel/agent/app/repo"

	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/1Panel-dev/1Panel/agent/buserr"
	"github.com/1Panel-dev/1Panel/agent/constant"
	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/1Panel-dev/1Panel/agent/utils/cmd"
	"github.com/1Panel-dev/1Panel/agent/utils/common"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

type ImageRepoService struct{}

type IImageRepoService interface {
	Page(search dto.SearchWithPage) (int64, interface{}, error)
	List() ([]dto.ImageRepoOption, error)
	Login(req dto.OperateByID) error
	Create(req dto.ImageRepoCreate) error
	Update(req dto.ImageRepoUpdate) error
	Delete(req dto.OperateByID) error
}

func NewIImageRepoService() IImageRepoService {
	return &ImageRepoService{}
}

func (u *ImageRepoService) Page(req dto.SearchWithPage) (int64, interface{}, error) {
	total, ops, err := imageRepoRepo.Page(req.Page, req.PageSize, repo.WithByLikeName(req.Info), repo.WithOrderBy("created_at desc"))
	var dtoOps []dto.ImageRepoInfo
	for _, op := range ops {
		var item dto.ImageRepoInfo
		if err := copier.Copy(&item, &op); err != nil {
			return 0, nil, buserr.WithDetail("ErrStructTransform", err.Error(), nil)
		}
		dtoOps = append(dtoOps, item)
	}
	return total, dtoOps, err
}

func (u *ImageRepoService) Login(req dto.OperateByID) error {
	repo, err := imageRepoRepo.Get(repo.WithByID(req.ID))
	if err != nil {
		return err
	}
	if repo.Auth {
		if err := u.CheckConn(repo.DownloadUrl, repo.Username, repo.Password); err != nil {
			_ = imageRepoRepo.Update(repo.ID, map[string]interface{}{"status": constant.StatusFailed, "message": err.Error()})
			return err
		}
	}
	_ = imageRepoRepo.Update(repo.ID, map[string]interface{}{"status": constant.StatusSuccess})
	return nil
}

func (u *ImageRepoService) List() ([]dto.ImageRepoOption, error) {
	ops, err := imageRepoRepo.List(repo.WithOrderBy("created_at desc"))
	var dtoOps []dto.ImageRepoOption
	for _, op := range ops {
		if op.Status == constant.StatusSuccess {
			var item dto.ImageRepoOption
			if err := copier.Copy(&item, &op); err != nil {
				return nil, buserr.WithDetail("ErrStructTransform", err.Error(), nil)
			}
			dtoOps = append(dtoOps, item)
		}
	}
	return dtoOps, err
}

func (u *ImageRepoService) Create(req dto.ImageRepoCreate) error {
	if cmd.CheckIllegal(req.Username, req.Password, req.DownloadUrl) {
		return buserr.New("ErrCmdIllegal")
	}
	imageRepo, _ := imageRepoRepo.Get(repo.WithByName(req.Name))
	if imageRepo.ID != 0 {
		return buserr.New("ErrRecordExist")
	}

	if req.Protocol == "http" {
		if err := u.handleRegistries(req.DownloadUrl, "", "create"); err != nil {
			return fmt.Errorf("create registry %s failed, err: %v", req.DownloadUrl, err)
		}
		if err := stopBeforeUpdateRepo(); err != nil {
			return err
		}
	}

	if req.Auth {
		if err := u.CheckConn(req.DownloadUrl, req.Username, req.Password); err != nil {
			return err
		}
	}

	if err := copier.Copy(&imageRepo, &req); err != nil {
		return buserr.WithDetail("ErrStructTransform", err.Error(), nil)
	}

	imageRepo.Status = constant.StatusSuccess
	return imageRepoRepo.Create(&imageRepo)
}

func (u *ImageRepoService) Delete(req dto.OperateByID) error {
	if req.ID == 1 {
		return errors.New("The default value cannot be edit !")
	}
	itemRepo, _ := imageRepoRepo.Get(repo.WithByID(req.ID))
	if itemRepo.ID == 0 {
		return buserr.New("ErrRecordNotFound")
	}
	if itemRepo.Auth {
		_, _ = cmd.NewCommandMgr().RunWithStdout("docker", "logout", "-i", itemRepo.DownloadUrl)
	}
	if itemRepo.Protocol == "https" {
		return imageRepoRepo.Delete(repo.WithByID(req.ID))
	}
	if err := u.handleRegistries("", itemRepo.DownloadUrl, "delete"); err != nil {
		return fmt.Errorf("delete registry %s failed, err: %v", itemRepo.DownloadUrl, err)
	}
	if err := validateDockerConfig(); err != nil {
		return err
	}
	if err := imageRepoRepo.Delete(repo.WithByID(req.ID)); err != nil {
		return err
	}
	go func() {
		_ = restartDocker()
	}()
	return nil
}

func (u *ImageRepoService) Update(req dto.ImageRepoUpdate) error {
	if req.ID == 1 {
		return errors.New("The default value cannot be deleted !")
	}
	if cmd.CheckIllegal(req.Username, req.Password, req.DownloadUrl) {
		return buserr.New("ErrCmdIllegal")
	}
	repo, err := imageRepoRepo.Get(repo.WithByID(req.ID))
	if err != nil {
		return err
	}
	needRestart := false
	if repo.Protocol == "http" && req.Protocol == "https" {
		if err := u.handleRegistries("", repo.DownloadUrl, "delete"); err != nil {
			return fmt.Errorf("delete registry %s failed, err: %v", repo.DownloadUrl, err)
		}
		needRestart = true
	}
	if repo.Protocol == "http" && req.Protocol == "http" {
		if err := u.handleRegistries(req.DownloadUrl, repo.DownloadUrl, "update"); err != nil {
			return fmt.Errorf("update registry %s => %s failed, err: %v", repo.DownloadUrl, req.DownloadUrl, err)
		}
		needRestart = repo.DownloadUrl == req.DownloadUrl
	}
	if repo.Protocol == "https" && req.Protocol == "http" {
		if err := u.handleRegistries(req.DownloadUrl, "", "create"); err != nil {
			return fmt.Errorf("create registry %s failed, err: %v", req.DownloadUrl, err)
		}
		needRestart = true
	}
	if needRestart {
		if err := stopBeforeUpdateRepo(); err != nil {
			return err
		}
	}
	if repo.Auth {
		_, _ = cmd.NewCommandMgr().RunWithStdout("docker", "logout", "-i", repo.DownloadUrl)
	}
	if req.Auth {
		if err := u.CheckConn(req.DownloadUrl, req.Username, req.Password); err != nil {
			return err
		}
	}

	upMap := make(map[string]interface{})
	upMap["download_url"] = req.DownloadUrl
	upMap["protocol"] = req.Protocol
	upMap["username"] = req.Username
	upMap["password"] = req.Password
	upMap["auth"] = req.Auth
	upMap["status"] = constant.StatusSuccess
	upMap["message"] = ""
	return imageRepoRepo.Update(req.ID, upMap)
}

func (u *ImageRepoService) CheckConn(host, user, password string) error {
	cmdMgr := cmd.NewCommandMgr()
	stdout, err := cmdMgr.RunWithStdout("docker", "login", "-u", user, "-p", password, host)
	if err != nil {
		return fmt.Errorf("stdout: %s, stderr: %v", stdout, err)
	}
	if strings.Contains(string(stdout), "Login Succeeded") {
		return nil
	}
	return errors.New(string(stdout))
}

func (u *ImageRepoService) handleRegistries(newHost, delHost, handle string) error {
	err := createIfNotExistDaemonJsonFile()
	if err != nil {
		return err
	}
	daemonMap := make(map[string]interface{})
	file, err := os.ReadFile(constant.DaemonJsonPath)
	if err != nil {
		return err
	}
	if len(file) != 0 {
		if err := json.Unmarshal(file, &daemonMap); err != nil {
			return err
		}
	}

	iRegistries := daemonMap["insecure-registries"]
	registries, _ := iRegistries.([]interface{})
	switch handle {
	case "create":
		registries = common.RemoveRepeatElement(append(registries, newHost))
	case "update":
		for i, regi := range registries {
			if regi == delHost {
				registries = append(registries[:i], registries[i+1:]...)
			}
		}
		registries = common.RemoveRepeatElement(append(registries, newHost))
	case "delete":
		for i, regi := range registries {
			if regi == delHost {
				registries = append(registries[:i], registries[i+1:]...)
			}
		}
	}
	if len(registries) == 0 {
		delete(daemonMap, "insecure-registries")
	} else {
		daemonMap["insecure-registries"] = registries
	}
	newJson, err := json.MarshalIndent(daemonMap, "", "\t")
	if err != nil {
		return err
	}
	if err := os.WriteFile(constant.DaemonJsonPath, newJson, 0640); err != nil {
		return err
	}
	return nil
}

func stopBeforeUpdateRepo() error {
	if err := validateDockerConfig(); err != nil {
		return err
	}
	if err := restartDocker(); err != nil {
		return err
	}
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	if err := func() error {
		for range ticker.C {
			select {
			case <-ctx.Done():
				cancel()
				return errors.New("the docker service cannot be restarted")
			default:
				stdout, err := cmd.RunDefaultWithStdoutBashC("systemctl is-active docker")
				if string(stdout) == "active\n" && err == nil {
					global.LOG.Info("docker restart with new conf successful!")
					return nil
				}
			}
		}
		return nil
	}(); err != nil {
		return err
	}
	return nil
}
