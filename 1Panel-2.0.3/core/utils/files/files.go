package files

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/1Panel-dev/1Panel/core/constant"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/1Panel-dev/1Panel/core/utils/cmd"
	"github.com/1Panel-dev/1Panel/core/utils/req_helper"
)

func CopyFile(src, dst string, withName bool) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	if path.Base(src) != path.Base(dst) && !withName {
		dst = path.Join(dst, path.Base(src))
	}
	if _, err := os.Stat(path.Dir(dst)); err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(path.Dir(dst), os.ModePerm)
		}
	}
	target, err := os.OpenFile(dst+"_temp", os.O_RDWR|os.O_CREATE|os.O_TRUNC, constant.FilePerm)
	if err != nil {
		return err
	}
	defer target.Close()

	if _, err = io.Copy(target, source); err != nil {
		return err
	}
	if err = os.Rename(dst+"_temp", dst); err != nil {
		return err
	}
	return nil
}

func CopyItem(isDir, withName bool, src, dst string) error {
	if path.Base(src) != path.Base(dst) && !withName {
		dst = path.Join(dst, path.Base(src))
	}
	srcInfo, err := os.Stat(path.Dir(src))
	if err != nil {
		return err
	}
	if _, err := os.Stat(dst); err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(dst, srcInfo.Mode())
		}
	}
	cmdStr := fmt.Sprintf(`cp -rf %s %s`, src, dst+"/")
	if !isDir {
		cmdStr = fmt.Sprintf(`cp -f %s %s`, src, dst+"/")
	}
	stdout, err := cmd.RunDefaultWithStdoutBashC(cmdStr)
	if err != nil {
		return fmt.Errorf("handle %s failed, stdout: %s, err: %v", cmdStr, stdout, err)
	}
	return nil
}

func HandleTar(sourceDir, targetDir, name, exclusionRules string, secret string) error {
	if _, err := os.Stat(targetDir); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(targetDir, os.ModePerm); err != nil {
			return err
		}
	}

	exMap := make(map[string]struct{})
	excludes := strings.Split(exclusionRules, ",")
	excludeRules := ""
	for _, exclude := range excludes {
		if len(exclude) == 0 {
			continue
		}
		if _, ok := exMap[exclude]; ok {
			continue
		}
		excludeRules += fmt.Sprintf(" --exclude '%s'", exclude)
		exMap[exclude] = struct{}{}
	}
	path := ""
	if strings.Contains(sourceDir, "/") {
		itemDir := strings.ReplaceAll(sourceDir[strings.LastIndex(sourceDir, "/"):], "/", "")
		aheadDir := sourceDir[:strings.LastIndex(sourceDir, "/")]
		if len(aheadDir) == 0 {
			aheadDir = "/"
		}
		path += fmt.Sprintf("-C %s %s", aheadDir, itemDir)
	} else {
		path = sourceDir
	}

	commands := ""

	if len(secret) != 0 {
		extraCmd := "| openssl enc -aes-256-cbc -salt -k '" + secret + "' -out"
		commands = fmt.Sprintf("tar -zcf %s %s %s %s", " -"+excludeRules, path, extraCmd, targetDir+"/"+name)
		global.LOG.Debug(strings.ReplaceAll(commands, fmt.Sprintf(" %s ", secret), "******"))
	} else {
		commands = fmt.Sprintf("tar -zcf %s %s %s", targetDir+"/"+name, excludeRules, path)
		global.LOG.Debug(commands)
	}
	cmdMgr := cmd.NewCommandMgr(cmd.WithTimeout(24*time.Hour), cmd.WithIgnoreExist1())
	stdout, err := cmdMgr.RunWithStdoutBashC(commands)
	if err != nil {
		if len(stdout) != 0 {
			global.LOG.Errorf("do handle tar failed, stdout: %s, err: %v", stdout, err)
			return fmt.Errorf("do handle tar failed, stdout: %s, err: %v", stdout, err)
		}
	}
	return nil
}

func HandleUnTar(sourceFile, targetDir string, secret string) error {
	if _, err := os.Stat(targetDir); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(targetDir, os.ModePerm); err != nil {
			return err
		}
	}
	commands := ""
	if len(secret) != 0 {
		extraCmd := "openssl enc -d -aes-256-cbc -k '" + secret + "' -in " + sourceFile + " | "
		commands = fmt.Sprintf("%s tar -zxvf - -C %s", extraCmd, targetDir+" > /dev/null 2>&1")
		global.LOG.Debug(strings.ReplaceAll(commands, fmt.Sprintf(" %s ", secret), "******"))
	} else {
		commands = fmt.Sprintf("tar zxvfC %s %s", sourceFile, targetDir)
		global.LOG.Debug(commands)
	}

	cmdMgr := cmd.NewCommandMgr(cmd.WithTimeout(24 * time.Hour))
	stdout, err := cmdMgr.RunWithStdoutBashC(commands)
	if err != nil {
		global.LOG.Errorf("do handle untar failed, stdout: %s, err: %v", stdout, err)
		return errors.New(stdout)
	}
	return nil
}

func DownloadFile(url, dst string) error {
	resp, err := req_helper.HandleGet(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create download file [%s] error, err %s", dst, err.Error())
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("save download file [%s] error, err %s", dst, err.Error())
	}
	return nil
}

func DownloadFileWithProxy(url, dst string) error {
	_, resp, err := req_helper.HandleRequestWithProxy(url, http.MethodGet, constant.TimeOut5m)
	if err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create download file [%s] error, err %s", dst, err.Error())
	}
	defer out.Close()

	reader := bytes.NewReader(resp)
	if _, err = io.Copy(out, reader); err != nil {
		return fmt.Errorf("save download file [%s] error, err %s", dst, err.Error())
	}
	return nil
}

func Stat(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func GetFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()

	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
