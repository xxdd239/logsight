package config

import (
	"os"
	"path"
	"strconv"

	"github.com/talkincode/logsight/common"
	"gopkg.in/yaml.v3"
)

// DBConfig 数据库(PostgreSQL)配置
type DBConfig struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Passwd   string `yaml:"passwd"`
	MaxConn  int    `yaml:"max_conn"`
	IdleConn int    `yaml:"idle_conn"`
	Debug    bool   `yaml:"debug"`
}

// SysConfig 系统配置
type SysConfig struct {
	Appid    string `yaml:"appid"`
	Location string `yaml:"location"`
	Workdir  string `yaml:"workdir"`
	Debug    bool   `yaml:"debug"`
}

// WebConfig WEB 配置
type WebConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	TlsPort int    `yaml:"tls_port"`
	Secret  string `yaml:"secret"`
}

type LogConfig struct {
	Mode           string `yaml:"mode"`
	ConsoleEnable  bool   `yaml:"console_enable"`
	LokiEnable     bool   `yaml:"loki_enable"`
	FileEnable     bool   `yaml:"file_enable"`
	Filename       string `yaml:"filename"`
	QueueSize      int    `yaml:"queue_size"`
	LokiApi        string `yaml:"loki_api"`
	LokiUser       string `yaml:"loki_user"`
	LokiPwd        string `yaml:"loki_pwd"`
	LokiJob        string `yaml:"loki_job"`
	MetricsStorage string `yaml:"metrics_storage"`
	MetricsHistory int    `yaml:"metrics_history"`
}

type SyslogdConfig struct {
	Host  string `yaml:"host" json:"host"`
	Port  int    `yaml:"port" json:"port"`
	Debug bool   `yaml:"debug" json:"debug"`
}

type AppConfig struct {
	System   SysConfig     `yaml:"system" json:"system"`
	Web      WebConfig     `yaml:"web" json:"web"`
	Database DBConfig      `yaml:"database" json:"database"`
	Syslogd  SyslogdConfig `yaml:"syslogd" json:"syslogd"`
	Logger   LogConfig     `yaml:"logger" json:"logger"`
}

func (c *AppConfig) GetLogDir() string {
	return path.Join(c.System.Workdir, "logs")
}

func (c *AppConfig) GetPublicDir() string {
	return path.Join(c.System.Workdir, "public")
}

func (c *AppConfig) GetPrivateDir() string {
	return path.Join(c.System.Workdir, "private")
}

func (c *AppConfig) GetDataDir() string {
	return path.Join(c.System.Workdir, "data")
}
func (c *AppConfig) GetBackupDir() string {
	return path.Join(c.System.Workdir, "backup")
}

func (c *AppConfig) initDirs() {
	_ = os.MkdirAll(path.Join(c.System.Workdir, "logs"), 0755)
	_ = os.MkdirAll(path.Join(c.System.Workdir, "public"), 0755)
	_ = os.MkdirAll(path.Join(c.System.Workdir, "data"), 0755)
	_ = os.MkdirAll(path.Join(c.System.Workdir, "data/metrics"), 0755)
	_ = os.MkdirAll(path.Join(c.System.Workdir, "private"), 0644)
	_ = os.MkdirAll(path.Join(c.System.Workdir, "backup"), 0644)
}

func setEnvValue(name string, val *string) {
	var evalue = os.Getenv(name)
	if evalue != "" {
		*val = evalue
	}
}

func setEnvBoolValue(name string, val *bool) {
	var evalue = os.Getenv(name)
	if evalue != "" {
		*val = evalue == "true" || evalue == "1" || evalue == "on"
	}
}

func setEnvInt64Value(name string, val *int64) {
	var evalue = os.Getenv(name)
	if evalue == "" {
		return
	}

	p, err := strconv.ParseInt(evalue, 10, 64)
	if err == nil {
		*val = p
	}
}
func setEnvIntValue(name string, val *int) {
	var evalue = os.Getenv(name)
	if evalue == "" {
		return
	}

	p, err := strconv.ParseInt(evalue, 10, 64)
	if err == nil {
		*val = int(p)
	}
}

var DefaultAppConfig = &AppConfig{
	System: SysConfig{
		Appid:    "Logsight",
		Location: "Asia/Shanghai",
		Workdir:  "/var/logsight",
		Debug:    true,
	},
	Web: WebConfig{
		Host:    "0.0.0.0",
		Port:    1816,
		TlsPort: 1817,
		Secret:  "9b6de5cc-0731-1203-xxtt-0f568ac9da37",
	},
	Database: DBConfig{
		Type:     "postgres",
		Host:     "127.0.0.1",
		Port:     5432,
		Name:     "logsight",
		User:     "postgres",
		Passwd:   "myroot",
		MaxConn:  100,
		IdleConn: 10,
		Debug:    false,
	},
	Syslogd: SyslogdConfig{
		Host:  "0.0.0.0",
		Port:  1814,
		Debug: true,
	},
	Logger: LogConfig{
		Mode:           "development",
		ConsoleEnable:  true,
		LokiEnable:     false,
		FileEnable:     true,
		Filename:       "/var/logsight/logsight.log",
		QueueSize:      4096,
		LokiApi:        "http://127.0.0.1:3100",
		LokiUser:       "logsight",
		LokiPwd:        "logsight",
		LokiJob:        "logsight",
		MetricsStorage: "/var/logsight/data/metrics",
		MetricsHistory: 24 * 7,
	},
}

func LoadConfig(cfile string) *AppConfig {
	// 开发环境首先查找当前目录是否存在自定义配置文件
	if cfile == "" {
		cfile = "logsight.yml"
	}
	if !common.FileExists(cfile) {
		cfile = "/etc/logsight.yml"
	}
	cfg := new(AppConfig)
	if common.FileExists(cfile) {
		data := common.Must2(os.ReadFile(cfile))
		common.Must(yaml.Unmarshal(data.([]byte), cfg))
	} else {
		cfg = DefaultAppConfig
	}

	cfg.initDirs()

	setEnvValue("LOGSIGHT_SYSTEM_WORKER_DIR", &cfg.System.Workdir)
	setEnvBoolValue("LOGSIGHT_SYSTEM_DEBUG", &cfg.System.Debug)

	setEnvValue("LOGSIGHT_SYSLOG_HOST", &cfg.Syslogd.Host)
	setEnvIntValue("LOGSIGHT_SYSLOG_PORT", &cfg.Syslogd.Port)
	setEnvBoolValue("LOGSIGHT_SYSLOG_DEBUG", &cfg.Syslogd.Debug)

	// WEB
	setEnvValue("LOGSIGHT_WEB_HOST", &cfg.Web.Host)
	setEnvValue("LOGSIGHT_WEB_SECRET", &cfg.Web.Secret)
	setEnvIntValue("LOGSIGHT_WEB_PORT", &cfg.Web.Port)
	setEnvIntValue("LOGSIGHT_WEB_TLS_PORT", &cfg.Web.TlsPort)

	// DB
	setEnvValue("LOGSIGHT_DB_HOST", &cfg.Database.Host)
	setEnvValue("LOGSIGHT_DB_NAME", &cfg.Database.Name)
	setEnvValue("LOGSIGHT_DB_USER", &cfg.Database.User)
	setEnvValue("LOGSIGHT_DB_PWD", &cfg.Database.Passwd)
	setEnvIntValue("LOGSIGHT_DB_PORT", &cfg.Database.Port)
	setEnvBoolValue("LOGSIGHT_DB_DEBUG", &cfg.Database.Debug)

	setEnvValue("LOGSIGHT_LOKI_JOB", &cfg.Logger.LokiJob)
	setEnvValue("LOGSIGHT_LOKI_SERVER", &cfg.Logger.LokiApi)
	setEnvValue("LOGSIGHT_LOKI_USERNAME", &cfg.Logger.LokiUser)
	setEnvValue("LOGSIGHT_LOKI_PASSWORD", &cfg.Logger.LokiPwd)
	setEnvValue("LOGSIGHT_LOGGER_MODE", &cfg.Logger.Mode)
	setEnvBoolValue("LOGSIGHT_LOKI_ENABLE", &cfg.Logger.LokiEnable)
	setEnvBoolValue("LOGSIGHT_LOGGER_FILE_ENABLE", &cfg.Logger.FileEnable)

	return cfg
}
