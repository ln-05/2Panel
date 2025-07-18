package constant

import "sync/atomic"

type DBContext string

const (
	TimeOut5s  = 5
	TimeOut20s = 20
	TimeOut5m  = 300

	DateLayout         = "2006-01-02" // or use time.DateOnly while go version >= 1.20
	DefaultDate        = "1970-01-01"
	DateTimeLayout     = "2006-01-02 15:04:05" // or use time.DateTime while go version >= 1.20
	DateTimeSlimLayout = "20060102150405"

	OrderDesc = "descending"
	OrderAsc  = "ascending"

	// backup
	S3          = "S3"
	OSS         = "OSS"
	Sftp        = "SFTP"
	OneDrive    = "OneDrive"
	MinIo       = "MINIO"
	Cos         = "COS"
	Kodo        = "KODO"
	WebDAV      = "WebDAV"
	Local       = "LOCAL"
	UPYUN       = "UPYUN"
	ALIYUN      = "ALIYUN"
	GoogleDrive = "GoogleDrive"

	OneDriveRedirectURI = "http://localhost/login/authorized"
	GoogleRedirectURI   = "https://localhost:8080"
)

const (
	DirPerm  = 0755
	FilePerm = 0644
)

const (
	SyncSystemProxy                  = "SyncSystemProxy"
	SyncScripts                      = "SyncScripts"
	SyncBackupAccounts               = "SyncBackupAccounts"
	SyncAlertSetting                 = "SyncAlertSetting"
	SyncCustomApp                    = "SyncCustomApp"
	SyncLanguage                     = "SyncLanguage"
	SyncSystemProxyWithRestartDocker = "SyncSystemProxyWithRestartDocker"
)

var WebUrlMap = map[string]struct{}{
	"/apps":           {},
	"/apps/all":       {},
	"/apps/installed": {},
	"/apps/upgrade":   {},
	"/apps/setting":   {},

	"/ai":       {},
	"/ai/model": {},
	"/ai/gpu":   {},
	"/ai/mcp":   {},

	"/containers":                   {},
	"/containers/container/operate": {},
	"/containers/container":         {},
	"/containers/image":             {},
	"/containers/network":           {},
	"/containers/volume":            {},
	"/containers/repo":              {},
	"/containers/compose":           {},
	"/containers/template":          {},
	"/containers/setting":           {},
	"/containers/dashboard":         {},

	"/cronjobs":                 {},
	"/cronjobs/cronjob":         {},
	"/cronjobs/library":         {},
	"/cronjobs/cronjob/operate": {},

	"/databases":                   {},
	"/databases/mysql":             {},
	"/databases/mysql/remote":      {},
	"/databases/postgresql":        {},
	"/databases/postgresql/remote": {},
	"/databases/redis":             {},
	"/databases/redis/remote":      {},

	"/hosts":                  {},
	"/hosts/files":            {},
	"/hosts/monitor/monitor":  {},
	"/hosts/monitor/setting":  {},
	"/hosts/firewall/port":    {},
	"/hosts/firewall/forward": {},
	"/hosts/firewall/ip":      {},
	"/hosts/process/process":  {},
	"/hosts/process/network":  {},
	"/hosts/ssh/ssh":          {},
	"/hosts/ssh/log":          {},
	"/hosts/ssh/session":      {},

	"/terminal": {},

	"/logs":           {},
	"/logs/operation": {},
	"/logs/login":     {},
	"/logs/website":   {},
	"/logs/system":    {},
	"/logs/ssh":       {},
	"/logs/task":      {},

	"/settings":               {},
	"/settings/panel":         {},
	"/settings/backupaccount": {},
	"/settings/license":       {},
	"/settings/about":         {},
	"/settings/safe":          {},
	"/settings/snapshot":      {},
	"/settings/expired":       {},

	"/toolbox":              {},
	"/toolbox/device":       {},
	"/toolbox/supervisor":   {},
	"/toolbox/clam":         {},
	"/toolbox/clam/setting": {},
	"/toolbox/ftp":          {},
	"/toolbox/fail2ban":     {},
	"/toolbox/clean":        {},

	"/websites":                 {},
	"/websites/ssl":             {},
	"/websites/runtimes/php":    {},
	"/websites/runtimes/node":   {},
	"/websites/runtimes/java":   {},
	"/websites/runtimes/go":     {},
	"/websites/runtimes/python": {},
	"/websites/runtimes/dotnet": {},

	"/login": {},

	"/xpack":                {},
	"/xpack/waf/dashboard":  {},
	"/xpack/waf/global":     {},
	"/xpack/waf/websites":   {},
	"/xpack/waf/log":        {},
	"/xpack/waf/block":      {},
	"/xpack/waf/blackwhite": {},
	"/xpack/waf/stat":       {},

	"/xpack/monitor/dashboard": {},
	"/xpack/monitor/setting":   {},
	"/xpack/monitor/rank":      {},
	"/xpack/monitor/log":       {},
	"/xpack/monitor/trend":     {},
	"/xpack/monitor/websites":  {},

	"/xpack/tamper":          {},
	"/xpack/gpu":             {},
	"/xpack/alert/dashboard": {},
	"/xpack/alert/log":       {},
	"/xpack/alert/setting":   {},
	"/xpack/setting":         {},
	"/xpack/node":            {},
	"/xpack/exchange/file":   {},
	"/xpack/app":             {},
}

var DynamicRoutes = []string{
	`^/containers/composeDetail/[^/]+$`,
	`^/databases/mysql/setting/[^/]+/[^/]+$`,
	`^/databases/postgresql/setting/[^/]+/[^/]+$`,
	`^/websites/[^/]+/config/[^/]+$`,
}

var CertStore atomic.Value

var DaemonJsonPath = "/etc/docker/daemon.json"
