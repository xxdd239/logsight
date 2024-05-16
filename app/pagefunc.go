package app

import (
	"strings"
	"time"

	"github.com/talkincode/logsight/assets"
)

func (a *Application) GetTemplateFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"pagever": func() int64 {
			if a.appConfig.System.Debug {
				return time.Now().Unix()
			} else {
				return int64(time.Now().Hour())
			}
		},
		"buildver": func() string {
			bv := strings.TrimSpace(assets.BuildVersion())
			if bv != "" {
				return bv
			}
			return "develop-" + time.Now().Format(time.RFC3339)
		},
		"moontheme": func() string {
			theme := a.GetSystemTheme()
			if theme == "dark" {
				return "1"
			}
			return "0"
		},
		"theme": func() string {
			return a.GetSystemTheme()
		},
		"sys_config": func(name string) string {
			return a.GetSettingsStringValue("system", name)
		},
	}
}
