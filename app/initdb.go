package app

import (
	"time"

	"github.com/talkincode/logsight/common"
	"github.com/talkincode/logsight/models"
)

func (a *Application) checkSuper() {
	var count int64
	a.gormDB.Model(&models.SysOpr{}).Where("username='admin' and level = ?", "super").Count(&count)
	if count == 0 {
		a.gormDB.Create(&models.SysOpr{
			ID:        common.UUIDint64(),
			Realname:  "administrator",
			Mobile:    "0000",
			Email:     "N/A",
			Username:  "admin",
			Password:  common.Sha256HashWithSalt("logsight", common.SecretSalt),
			Level:     "super",
			Status:    "enabled",
			Remark:    "super",
			LastLogin: time.Now(),
		})
	}
}

func (a *Application) checkSettings() {
	var checkConfig = func(sortid int, stype, cname, value, remark string) {
		var count int64
		a.gormDB.Model(&models.SysConfig{}).Where("type = ? and name = ?", stype, cname).Count(&count)
		if count == 0 {
			a.gormDB.Create(&models.SysConfig{ID: 0, Sort: sortid, Type: stype, Name: cname, Value: value, Remark: remark})
		}
	}

	for sortid, name := range ConfigConstants {
		switch name {
		case ConfigSystemTitle:
			checkConfig(sortid, "system", ConfigSystemTitle, "LogSight Management System", "System title")
		case ConfigSystemTheme:
			checkConfig(sortid, "system", ConfigSystemTheme, "light", "System theme")
		case ConfigSystemLoginRemark:
			checkConfig(sortid, "system", ConfigSystemLoginRemark, "Recommended browser: Chrome/Edge", "Login page description")
		case ConfigSystemLoginSubtitle:
			checkConfig(sortid, "system", ConfigSystemLoginSubtitle, "LogSight Community Edition", "Login form title")

		}
	}
}
