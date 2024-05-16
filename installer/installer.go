package installer

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/talkincode/logsight/common"
	"github.com/talkincode/logsight/config"
	"gopkg.in/yaml.v3"
)

var installScript = `#!/bin/bash -x
mkdir -p /var/logsight
chmod -R 755 /var/logsight
install -m 755 {{binfile}} /usr/local/bin/logsight
test -d /usr/lib/systemd/system || mkdir -p /usr/lib/systemd/system
cat>/usr/lib/systemd/system/logsight.service<<EOF
[Unit]
Description=logsight
After=network.target
StartLimitIntervalSec=0

[Service]
Restart=always
RestartSec=1
Environment=GODEBUG=x509ignoreCN=0
LimitNOFILE=65535
LimitNPROC=65535
User=root
ExecStart=/usr/local/bin/logsight

[Install]
WantedBy=multi-user.target
EOF

chmod 600 /usr/lib/systemd/system/logsight.service
systemctl enable logsight && systemctl daemon-reload
`

func InitConfig(config *config.AppConfig) error {
	// config.NBI.JwtSecret = common.UUID()
	cfgstr, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile("/etc/logsight.yml", cfgstr, 0644)
}

func Install() error {
	if !common.FileExists("/etc/logsight.yml") {
		err := InitConfig(config.DefaultAppConfig)
		if err != nil {
			return err
		}
	}
	// Get the absolute path of the currently executing file
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	dir := filepath.Dir(path)
	binfile := filepath.Join(dir, "logsight")
	installScript = strings.ReplaceAll(installScript, "{{binfile}}", binfile)
	_ = os.WriteFile("/tmp/logsight_install.sh", []byte(installScript), 0755)

	// 创建用户&组
	if err := exec.Command("/bin/bash", "/tmp/logsight_install.sh").Run(); err != nil {
		return err
	}

	return os.Remove("/tmp/logsight_install.sh")
}

func Uninstall() {
	_ = os.Remove("/usr/lib/systemd/system/logsight.service")
	_ = os.Remove("/usr/local/bin/logsight")
}
